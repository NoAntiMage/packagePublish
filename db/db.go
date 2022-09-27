package db

import (
	"PackageServer/config"
	"PackageServer/logger"
	"PackageServer/model"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	//	ormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var Db *gorm.DB

func InitDb() {
	if config.ServerConf.Name == "manager" {
		InitDbForManager()
	} else if config.ServerConf.Name == "worker" {
		logger.Log.Info("worker db in next version...")
	} else {
		logger.Log.Error("Invalid Server Type Name for Db")
	}
}

func InitDbForManager() {
	dbFileLocation := fmt.Sprintf("%v/%v.db", config.ServerConf.DataDir, config.ServerConf.Name)
	db, err := gorm.Open(sqlite.Open(dbFileLocation), &gorm.Config{
		//		Logger: ormlogger.Default.LogMode(ormlogger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		logger.Log.Error(err)
	}

	logger.Log.Debugf("initDb: db: %v", db)

	db.AutoMigrate(
		&model.AreaInfo{},
		&model.ServiceOnline{},
		&model.PublishPlanLog{},
	)

	Db = db
}
