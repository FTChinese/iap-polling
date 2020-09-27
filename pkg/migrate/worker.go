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
)

var (
	maxWorkers = 16
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
	log.Printf("Verify: start verifying %s", f)

	receipt, err := ioutil.ReadFile(f)
	if err != nil {
		log.Printf("Veify: error reading receipt %s: %v", f, err)
		return apple.Subscription{}, err
	}

	body, err := w.subsClient.VerifyReceipt(string(receipt))
	if err != nil {
		return apple.Subscription{}, err
	}

	log.Printf("Verify: subscription response %s", body)

	var s apple.Subscription
	if err := json.Unmarshal(body, &s); err != nil {
		log.Printf("Verify: unmarshal subscription error %v", err)
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

func (w Worker) Start(dir string) error {
	fileCh := make(chan string)

	ctx := context.Background()

	go func() {
		err := WalkDir(fileCh, dir)
		if err != nil {
			log.Println(err)
		}
	}()

	for f := range fileCh {
		err := sem.Acquire(ctx, 1)
		if err != nil {
			break
		}

		go func(filename string) {
			defer sem.Release(1)

			_, _ = w.Verify(filename)
		}(f)
	}

	if err := sem.Acquire(ctx, int64(maxWorkers)); err != nil {
		log.Printf("Failed to acquire semaphore: %v", err)
		return nil
	}

	return nil
}
