package main

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
)


func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("seedig transaction table...")
		_, err := db.Exec(`
		INSERT INTO transactions
		(
			id, 
		 	user_id,
		 	state,
		 	transaction_code, 
		 	amount, 
		 	source_account, 
		 	destination_account, 
			created_at,
			authorization_method
		)
		VALUES(
			'16e70cef-04ba-45eb-a509-5460ea4c5316', 
			'4c4d8623-89ee-4d30-a42b-cecef89c352a', 
			1, 
			'T001',
			 3000,
			'10001',
			'10002', 
			'2020-07-31 07:56:15.092',
			1
		),(
			'e3dfd59e-e735-4ffd-8816-7127f81542d4', 
			'4c4d8623-89ee-4d30-a42b-cecef89c352a', 
			1, 
			'T002',
			 3000,
			'10001',
			'10002', 
			'2020-07-31 07:56:15.092',
			2
		),(
			'34a58731-284a-4e98-aa42-05a032da2469', 
			'4c4d8623-89ee-4d30-a42b-cecef89c352a', 
			2, 
			'T001',
			 3000,
			'10001',
			'10002', 
			'2020-07-31 07:56:15.092',
			2
		),
		(
			'3b48dea5-ef3a-45b3-8205-82e606e4dc92', 
			'4c4d8623-89ee-4d30-a42b-cecef89c352a', 
			3, 
			'T001',
			 3000,
			'10001',
			'10002', 
			'2020-07-31 07:56:15.092',
			2
		);	
		`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("truncate transaction...")
		_, err := db.Exec(`TRUNCATE TABLE transactions`)
		return err
	})
}