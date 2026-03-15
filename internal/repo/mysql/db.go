package mysql

import (
	"github.com/boxuanduan/gochat/config"
	"github.com/boxuanduan/gochat/pkg/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Init(cfg *config.MySQLConfig) (*gorm.DB, error) {
	// Use GORM to connect to MySQL using the DSN from the config
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	// Register models here
	err = db.AutoMigrate(
		&model.User{},
		&model.Friend{},
		&model.GroupInfo{},
		&model.GroupMember{},
		&model.Message{},
		&model.UserConversation{},
		&model.Outbox{},
	)

	return db, err
}
