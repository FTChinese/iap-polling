package apple

const StmtSubs = `
SELECT environment,
	original_transaction_id
FROM premium.apple_subscription
WHERE DATEDIFF(expires_date_utc, UTC_DATE()) < 3
    AND auto_renewal = 1
    AND environment = 'production'
ORDER BY expires_date_utc`
