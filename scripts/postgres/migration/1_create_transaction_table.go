package main

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table transactions...")
		_, err := db.Exec(`
			create table if not exists transactions(
				id varchar primary key,
				user_id varchar not null,
				state numeric not null,
				authorization_method numeric not null,
				transaction_code varchar not null,
				amount numeric not null,
				source_account varchar not null,
				destination_account varchar not null,
				created_at timestamp not null
			);
		`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table transactions...")
		_, err := db.Exec(`DROP TABLE transactions`)
		return err
	})
}
