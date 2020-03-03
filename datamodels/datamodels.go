package datamodels

import (
	"database/sql/driver"
	"github.com/gin-gonic/gin/internal/json"
	"github.com/jinzhu/gorm"
	pg "github.com/jinzhu/gorm/dialects/postgres"
)

type PlatformType int8

const (
	ZingglobalPlatform   PlatformType = 1
	ZhidreamPlatform     PlatformType = 2
	HealthEnginePlatform PlatformType = 3
)

type JsonNumArray []int64

func (t JsonNumArray) Value() (driver.Value, error) {
	buf, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	j := pg.Jsonb{RawMessage: buf}
	return j.Value()
}

func (t JsonNumArray) Scan(value interface{}) error {
	j := &pg.Jsonb{}
	err := j.Scan(value)
	if err != nil {
		return err
	}

	return json.Unmarshal(j.RawMessage, t)
}

func GetModeList() []interface{} {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, name string) string {
		return "t_ow_" + name
	}

	return []interface{}{
		&AdminUser{},
		&Article{},
		&File{},
		&Image{},
	}
}
