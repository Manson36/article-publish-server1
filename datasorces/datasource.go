package datasorces

import "github.com/article-publish-server1/datamodels"

func init() {
	PqDB.initDB()
	RDS.initDB()

	PqDB.AutoMigrate(datamodels.GetModeList()...)

}
