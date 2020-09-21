package apple

import "github.com/FTChinese/go-rest/chrono"

// BaseSchema contains the shared fields of all schema.
type BaseSchema struct {
	Environment           Environment `json:"environment" db:"environment"`
	OriginalTransactionID string      `json:"originalTransactionId" db:"original_transaction_id"`
}

// ReceiptFileName builds the file name when persisting latest receipt to disk.
func (s BaseSchema) ReceiptFileName() string {
	return s.OriginalTransactionID + "_" + s.Environment.String() + ".txt"
}

// Create key used in redis: `iap:receipt:1000000922681985-Sandbox`
func (s BaseSchema) ReceiptKeyName() string {
	return "iap:receipt:" + s.OriginalTransactionID + "-" + s.Environment.String()
}

// Subscription contains a user's subscription data.
// It it built from app store's verification response.
// The original transaction id is used to uniquely identify a user.
type Subscription struct {
	BaseSchema
	LastTransactionID string      `json:"lastTransactionId" db:"last_transaction_id"`
	ProductID         string      `json:"productId" db:"product_id"`
	PurchaseDateUTC   chrono.Time `json:"purchaseDateUtc" db:"purchase_date_utc"`
	ExpiresDateUTC    chrono.Time `json:"expiresDateUtc" db:"expires_date_utc"`
	AutoRenewal       bool        `json:"autoRenewal" db:"auto_renewal"`
	CreatedUTC        chrono.Time `json:"createdUtc" db:"created_utc"`
	UpdatedUTC        chrono.Time `json:"updatedUtc" db:"updated_utc"`
}
