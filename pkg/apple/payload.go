package apple

type VerificationPayload struct {
	ReceiptData            string `json:"receipt-data"`             // Required. The Base64 encoded receipt data.
	Password               string `json:"password"`                 // Required. Your app's shared secret (a hexadecimal string).
	ExcludeOldTransactions bool   `json:"exclude-old-transactions"` // Set this value to true for the response to include only the latest renewal transaction for any subscriptions. Applicable only to auto-renewable subscriptions.
}
