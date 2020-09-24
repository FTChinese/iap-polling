package migrate

import (
	"context"
	"encoding/json"
	"github.com/FTChinese.com/iap-polling/pkg/apple"
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"github.com/FTChinese.com/iap-polling/pkg/db"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/semaphore"
	"io/ioutil"
	"log"
	"runtime"
)

var (
	maxWorkers = runtime.GOMAXPROCS(0)
	sem        = semaphore.NewWeighted(int64(maxWorkers))
)

type Worker struct {
	db         *sqlx.DB // Always use local database for logging.
	subsClient apple.SubsClient
}

func NewWorker(prodApi bool) Worker {
	return Worker{
		db:         db.MustNewDB(config.MustDBConn(false)),
		subsClient: apple.NewSubsClient(prodApi),
	}
}

func (w Worker) Verify(f string) (apple.Subscription, error) {
	receipt, err := ioutil.ReadFile(f)
	if err != nil {
		return apple.Subscription{}, err
	}

	body, err := w.subsClient.VerifyReceipt(string(receipt))
	if err != nil {
		return apple.Subscription{}, err
	}

	log.Printf("%s", body)

	var s apple.Subscription
	if err := json.Unmarshal(body, &s); err != nil {
		return apple.Subscription{}, err
	}

	return s, nil
}

func (w Worker) SaveMapping(m IDMapping) error {
	_, err := w.db.NamedExec(StmtSaveMapping, m)
	if err != nil {
		return err
	}

	return nil
}

func (w Worker) Start(k NamingKind) error {
	fileCh := make(chan string)

	ctx := context.Background()

	go func() {
		err := WalkDir(fileCh, k)
		if err != nil {
			log.Println(err)
		}
	}()

	for f := range fileCh {
		err := sem.Acquire(ctx, 1)
		if err != nil {
			break
		}

		m := NewIDMapping(f, k)

		go func(m IDMapping) {
			defer sem.Release(1)

			subs, err := w.Verify(m.AbsFilePath)

			if err != nil {
				log.Println(err)
				return
			}

			m.TxID = subs.OriginalTransactionID

			err = w.SaveMapping(m)
			if err != nil {
				log.Println(err)
			}
		}(m)
	}

	if err := sem.Acquire(ctx, int64(maxWorkers)); err != nil {
		log.Printf("Failed to acquire semaphore: %v", err)
		return nil
	}

	return nil
}
