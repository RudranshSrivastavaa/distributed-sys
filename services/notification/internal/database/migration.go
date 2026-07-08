package database

import (
	"github.com/rudransh/distributed-commerce/notification/internal/model"
)

func Migrate() error {

	return DB.AutoMigrate(

		&model.Notification{},
		&model.DeadLetterNotification{},
		&model.InboxMessage{},
	)

}