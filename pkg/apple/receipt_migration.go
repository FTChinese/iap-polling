package apple

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
	"io/ioutil"
	"log"
)

var (
	maxWorkers = 16
	sem        = semaphore.NewWeighted(int64(maxWorkers))
)

type ReceiptMigration struct {
	api    SubsAPI
	logger *zap.Logger
}

func NewReceiptMigration(prod bool, logger *zap.Logger) ReceiptMigration {
	return ReceiptMigration{
		api:    NewSubsAPI(prod),
		logger: logger,
	}
}

func (m ReceiptMigration) Verify(f string) (Subscription, error) {
	log.Printf("Verify: start verifying %s", f)

	receipt, err := ioutil.ReadFile(f)
	if err != nil {
		log.Printf("Veify: error reading receipt %s: %v", f, err)
		return Subscription{}, err
	}

	body, err := m.api.VerifyReceipt(string(receipt))
	if err != nil {
		return Subscription{}, err
	}

	log.Printf("Verify: subscription response %s", body)

	var s Subscription
	if err := json.Unmarshal(body, &s); err != nil {
		log.Printf("Verify: unmarshal subscription error %v", err)
		return Subscription{}, err
	}

	return s, nil
}

func (m ReceiptMigration) Start(dir string) error {

	ctx := context.Background()

	fileCh := WalkReceipts(dir)

	for f := range fileCh {
		err := sem.Acquire(ctx, 1)
		if err != nil {
			break
		}

		go func(filename string) {
			defer sem.Release(1)

			_, _ = m.Verify(filename)
		}(f)
	}

	if err := sem.Acquire(ctx, int64(maxWorkers)); err != nil {
		log.Printf("Failed to acquire semaphore: %v", err)
		return nil
	}

	return nil
}
