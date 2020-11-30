package apple

import (
	"context"
	"encoding/json"
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"github.com/FTChinese.com/iap-polling/pkg/db"
	"github.com/FTChinese/go-rest/render"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"log"
)

// LinkInput defines the request body to link IAP to ftc account.
type LinkInput struct {
	FtcID        string `json:"ftcId" db:"ftc_user_id"`
	OriginalTxID string `json:"originalTxId" db:"original_transaction_id"`
}

type LinkMigration struct {
	db     *sqlx.DB
	api    SubsAPI
	logger *zap.Logger
}

const stmtLinkedID = `
SELECT original_transaction_id,
	ftc_id AS ftc_user_id
FROM premium.apple_id_mapping
WHERE ftc_id IS NOT NULL
	AND original_transaction_id != ''`

const stmtLinkCol = `
error_field = :error_field,
error_code = :error_code,
error_message = :error_message,
updated_utc = UTC_TIMESTAMP()
`

const stmtLinkErr = `
INSERT INTO premium.apple_link_log
SET original_transaction_id = :original_transaction_id,
	ftc_id = :ftc_user_id,
` + stmtLinkCol + `,
	created_utc = UTC_TIMESTAMP()
ON DUPLICATE KEY UPDATE
` + stmtLinkCol

type LinkErrLog struct {
	LinkInput
	Field   string `db:"error_field"`
	Code    string `db:"error_code"`
	Message string `db:"error_message"`
}

func NewLinkMigration(logger *zap.Logger) LinkMigration {
	return LinkMigration{
		db:     db.MustNewDB(config.MustDBConn(false)),
		api:    NewSubsAPI(true),
		logger: logger,
	}
}

func (m LinkMigration) retrieveLink() <-chan LinkInput {
	defer m.logger.Sync()
	sugar := m.logger.Sugar()

	ch := make(chan LinkInput)

	go func() {
		defer close(ch)

		rows, err := m.db.Queryx(stmtLinkedID)
		if err != nil {
			sugar.Error(err)
			return
		}

		var link LinkInput
		for rows.Next() {
			err := rows.StructScan(&link)
			if err != nil {
				sugar.Error(err)
				continue
			}

			ch <- link
		}
	}()

	return ch
}

func (m LinkMigration) saveLinkErrLog(l LinkErrLog) error {
	_, err := m.db.NamedExec(stmtLinkErr, l)
	if err != nil {
		return err
	}

	return nil
}

func (m LinkMigration) link(input LinkInput) error {
	defer m.logger.Sync()
	sugar := m.logger.Sugar()

	resp, err := m.api.Link(input)
	if err != nil {
		sugar.Error(err)
		return err
	}

	sugar.Infof("Link response status code %d", resp.StatusCode)

	if resp.StatusCode < 400 {
		return nil
	}

	sugar.Infof("Link response error %d, %s", resp.StatusCode, resp.Body)

	var respErr render.ResponseError
	if err := json.Unmarshal(resp.Body, &respErr); err != nil {
		sugar.Error(err)
		return err
	}

	if respErr.Invalid != nil {
		return m.saveLinkErrLog(LinkErrLog{
			LinkInput: input,
			Field:     respErr.Invalid.Field,
			Code:      string(respErr.Invalid.Code),
			Message:   respErr.Message,
		})
	}

	return m.saveLinkErrLog(LinkErrLog{
		LinkInput: input,
		Message:   respErr.Message,
	})
}

func (m LinkMigration) Start() error {
	defer m.logger.Sync()
	sugar := m.logger.Sugar()

	ctx := context.Background()

	ch := m.retrieveLink()

	for input := range ch {
		err := sem.Acquire(ctx, 1)
		if err != nil {
			break
		}

		go func(i LinkInput) {
			defer sem.Release(1)

			err := m.link(i)
			if err != nil {
				sugar.Error(err)
			}
		}(input)
	}

	if err := sem.Acquire(ctx, int64(maxWorkers)); err != nil {
		log.Printf("Failed to acquire semaphore: %v", err)
		return nil
	}

	return nil
}
