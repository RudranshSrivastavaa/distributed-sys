package repository

import "gorm.io/gorm"

type TransactionManager interface {

	WithTransaction(
		fn func(tx *gorm.DB) error,
	) error

}