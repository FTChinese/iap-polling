package apple

// Subscription contains a user's subscription data.
// It it built from app store's verification response.
// The original transaction id is used to uniquely identify a user.
type Subscription struct {
	Environment           Environment `json:"environment" db:"environment"`
	OriginalTransactionID string      `json:"originalTransactionId" db:"original_transaction_id"`
}
