package apple

const StmtSubs = `
SELECT environment,
	original_transaction_id,
	last_transaction_id,
	product_id,
	purchase_date_utc,
	expires_date_utc,
	tier,
	cycle,
	auto_renewal,
	created_utc,
	updated_utc
FROM premium.apple_subscription
WHERE DATEDIFF(expires_date_utc, UTC_DATE()) < 3
    AND auto_renewal = 1
    AND environment = 'production'
ORDER BY expires_date_utc;`
