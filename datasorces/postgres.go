package datasorces

import (
	"fmt"
	"github.com/article-publish-server1/config"
	"github.com/jinzhu/gorm"
	"time"
)

type pqdb struct {
	*gorm.DB
}

func (p *pqdb) initDB() {
	db, err := gorm.Open("postgres", config.Postgres.GetURI())
	if err != nil {
		msg := fmt.Sprintf("init postgres db err, msg:%s", err.Error())
		panic(msg)
	}

	if err := db.DB().Ping(); err != nil {
		msg := fmt.Sprintf("ping postgres db error, host=%s, port=%s, errmsg:%s",
			config.Postgres.Host, config.Postgres.Port, err.Error())

		panic(msg)
	}

	db.LogMode(true)

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(30)
	db.DB().SetConnMaxLifetime(time.Hour)

	p.DB = db
}

var PqDB = &pqdb{}
