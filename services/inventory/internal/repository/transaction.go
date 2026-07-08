package repository

import "gorm.io/gorm"

type TransactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(
	db *gorm.DB,
) *TransactionManager {

	return &TransactionManager{
		db: db,
	}
}

func (tm *TransactionManager) Execute(
	fn func(tx *gorm.DB) error,
) error {

	return tm.db.Transaction(func(tx *gorm.DB) error {

		return fn(tx)

	})

}