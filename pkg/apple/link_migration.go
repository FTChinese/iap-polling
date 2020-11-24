package apple

import (
	"context"
	"encoding/json"
	"fmt"
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
	Force        bool   `json:"force"` // If FtcID is already linked to another apple id, it is not allowed to use this OriginalTxId's subscription unless Force is true.
}

type LinkMigration struct {
	db     *sqlx.DB
	api    SubsAPI
	logger *zap.Logger
}

const stmtLinkedID = `
SELECT i.vip_id AS ftc_user_id,
    i.transaction_id AS original_transaction_id
FROM premium.ftc_vip_ios AS i
    LEFT JOIN (
        SELECT COUNT(*) AS tx_count,
            transaction_id
        FROM premium.ftc_vip_ios
        WHERE vip_id_alias IS NULL 
            AND transaction_id IS NOT NULL
            AND transaction_id != ''
        GROUP BY transaction_id
    ) AS dup
    ON i.transaction_id = dup.transaction_id
WHERE i.vip_id_alias IS NULL 
    AND i.transaction_id IS NOT NULL
    AND i.transaction_id != ''
    AND dup.tx_count = 1`

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

func NewLinkMigration(prod bool, logger *zap.Logger) LinkMigration {
	return LinkMigration{
		db:     db.MustNewDB(config.MustDBConn(prod)),
		api:    NewSubsAPI(prod),
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

	if resp.StatusCode != 422 {
		return fmt.Errorf("link response error status: %d", resp.StatusCode)
	}

	var respErr render.ResponseError
	if err := json.Unmarshal(resp.Body, &respErr); err != nil {
		sugar.Error(err)
		return err
	}

	return m.saveLinkErrLog(LinkErrLog{
		LinkInput: input,
		Field:     respErr.Invalid.Field,
		Code:      string(respErr.Invalid.Code),
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
