package pg

import (
	"github.com/go-pg/pg/v10"
)

type CrudRepositoryWrapper struct {
	db *pg.DB
}

func Wrap(db *pg.DB) *CrudRepositoryWrapper {
	return &CrudRepositoryWrapper{db: db}
}

func (wrapper *CrudRepositoryWrapper) Save(model interface{}) error {
	count, err := wrapper.db.Model(model).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		_, err = wrapper.db.Model(model).Update()
		if err != nil {
			return err
		}
		return nil
	}
	_, err = wrapper.db.Model(model).Insert()
	if err != nil {
		return err
	}
	return nil
}

func (wrapper *CrudRepositoryWrapper) Load(model interface{}) error {

	err := wrapper.db.Model(model).WherePK().Select()
	if err == pg.ErrNoRows {
		return nil
	}
	return err
}

func (wrapper *CrudRepositoryWrapper) Remove(model interface{}) error {
	err := wrapper.db.Model(model).WherePK().Select()
	if err == pg.ErrNoRows {
		return nil
	}
	return err
}
