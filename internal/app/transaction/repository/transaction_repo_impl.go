package repository

import (
	"github.com/go-pg/pg/v10"
	"github.com/tunaiku/mobilebanking/internal/app/domain"
)

type TransactionRepository struct {
	db *pg.DB
}

func NewTransactionRepository(db *pg.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

// this func will act as upsert operation for the transaction
func (txrepo *TransactionRepository) SaveTransaction(tx *domain.Transaction) error {

	_, err := txrepo.db.Model(tx).OnConflict("(id) DO UPDATE").Insert(tx)
	if err != nil {
		return err
	}
	return nil

}

func (txrepo *TransactionRepository) ReadTransaction(ID string) (*domain.Transaction, error) {
	tx := &domain.Transaction{ID: ID}
	err := txrepo.db.Model(tx).WherePK().Select()
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return tx, err
}
