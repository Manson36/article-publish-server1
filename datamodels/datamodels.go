package datamodels

import "github.com/jinzhu/gorm"

type PlatformType int8

const (
	ZingglobalPlatform   PlatformType = 1
	ZhidreamPlatform     PlatformType = 2
	HealthEnginePlatform PlatformType = 3
)

func GetModeList() []interface{} {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, name string) string {
		return "t_ow_" + name
	}

	return []interface{}{
		&AdminUser{},
	}
}
