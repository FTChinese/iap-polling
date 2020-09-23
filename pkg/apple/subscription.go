package apple

// Subscription contains a user's subscription data.
// It it built from app store's verification response.
// The original transaction id is used to uniquely identify a user.
type Subscription struct {
	Environment           Environment `json:"environment" db:"environment"`
	OriginalTransactionID string      `json:"originalTransactionId" db:"original_transaction_id"`
}

func (s Subscription) String() string {
	return s.OriginalTransactionID + "," + s.Environment.String()
}

// ReceiptFileName builds the file name when persisting latest receipt to disk.
func (s Subscription) ReceiptFileName() string {
	return s.OriginalTransactionID + "_" + s.Environment.String() + ".txt"
}

// Create key used in redis: `iap:receipt:1000000922681985-Sandbox`
func (s Subscription) ReceiptKeyName() string {
	return "iap:receipt:" + s.OriginalTransactionID + "-" + s.Environment.String()
}
