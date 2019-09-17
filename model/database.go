package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rs/zerolog/log"
	"GoRestServer/model/auth"
	"GoRestServer/pkg/config"
)

var SQL *gorm.DB

func Setup() {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return config.Database.TablePrefix + defaultTableName
	}

	var err error
	SQL, err = gorm.Open(config.Database.Type, config.Database.User+":"+config.Database.Password+config.Database.Url)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to connect to database")
		panic(err)
	}

	SQL.DB().SetMaxOpenConns(100)
	SQL.DB().SetMaxIdleConns(10)
	SQL.SetLogger(&log.Logger)
	SQL.LogMode(config.Database.ShowSql)
	SQL.SingularTable(true)
}

func Close() {
	SQL.Close()
}

func AutoMigrate() {
	SQL.AutoMigrate(&model_auth.Function{})
	SQL.AutoMigrate(&model_auth.Role{})
	SQL.AutoMigrate(&model_auth.User{})
	SQL.Model(&model_auth.User{}).AddForeignKey("role_id", config.Database.TablePrefix+"role(id)", "no action", "no action")
}
