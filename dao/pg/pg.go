package pg

import (
	"IM-Backend/configs"
	"IM-Backend/global"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cf configs.AppConf) *gorm.DB {
	dsn := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai`, cf.DB.Host, cf.DB.User, cf.DB.PassWord, cf.DB.DBName, cf.DB.Port)
	// 连接到 PostgreSQL 数据库
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("connect to database failed:%v", err))
	}
	//if err := db.AutoMigrate(&table.Svc{}); err != nil {
	//	panic(err)
	//}

	global.Log.Info("connect to database successfully")

	return db
}
