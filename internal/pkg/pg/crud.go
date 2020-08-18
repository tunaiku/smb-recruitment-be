package pg

import (
	"github.com/go-pg/pg/v10"
)

// Refactoring crud repository
type CrudRepositoryWrapper struct {
	db *pg.DB
}

func Wrap(db *pg.DB) *CrudRepositoryWrapper {
	return &CrudRepositoryWrapper{db: db}
}

func (wrapper *CrudRepositoryWrapper) Save(model interface{}) error {
	_, err := wrapper.db.Model(model).OnConflict("DO UPDATE").Insert(model)
	return err
}

func (wrapper *CrudRepositoryWrapper) Load(model interface{}) error {

	err := wrapper.db.Model(model).WherePK().Select()
	if err == pg.ErrNoRows {
		return nil
	}
	return err
}

func (wrapper *CrudRepositoryWrapper) Remove(model interface{}) error {
	_, err := wrapper.db.Model(model).WherePK().Delete()
	if err == pg.ErrNoRows {
		return nil
	}
	return err
}
