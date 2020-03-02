package repositories

import (
	"encoding/json"
	"errors"
	"github.com/article-publish-server1/config"
	"github.com/article-publish-server1/datamodels"
	"github.com/article-publish-server1/datasorces"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type AdminUserRepository interface {
	Create(user *datamodels.AdminUser) error
	Get(query interface{}, args ...interface{}) (*datamodels.AdminUser, error)
	GetInfoFromCacheByID(int64) (*datamodels.AdminUser, error)
	SetInfoToCache(user *datamodels.AdminUser) error
}

type adminUserRepository struct {
	commonRepository
	rds *datasorces.Rds
}

func NewAdminUserRepository() AdminUserRepository {
	return &adminUserRepository{
		commonRepository: commonRepository{
			db: datasorces.PqDB.DB,
		},
		rds: datasorces.RDS,
	}
}

func (a adminUserRepository) GetInfoFromCacheByID(uid int64) (*datamodels.AdminUser, error) {
	key := config.Redis.KeyPrefix + "session" + strconv.FormatInt(uid, 10)
	reply, err := redis.String(a.rds.Do("get", key))
	if err != nil && redis.ErrNil != err {
		return nil, err
	}

	if reply == "" {
		return nil, nil
	}

	var u datamodels.AdminUser
	if err := json.Unmarshal([]byte(reply), &u); err != nil {
		_, _ = a.rds.Do("del", key)
		return nil, err
	}

	if u.ID == 0 {
		_, _ = a.rds.Do("del", key)
		return nil, nil
	}

	return &u, nil
}

func (a adminUserRepository) Create(user *datamodels.AdminUser) error {
	if user == nil {
		return errors.New("user params is nil")
	}

	return a.commonRepository.Create(user)
}

func (a adminUserRepository) Get(query interface{}, args ...interface{}) (*datamodels.AdminUser, error) {
	record, err := a.commonRepository.Get(&datamodels.AdminUser{}, query, args...)
	if record == nil {
		return nil, err
	}

	return record.(*datamodels.AdminUser), nil
}

func (a adminUserRepository) SetInfoToCache(user *datamodels.AdminUser) error {
	if user == nil || user.ID == 0 {
		return nil
	}

	key := config.Redis.KeyPrefix + "session:" + strconv.FormatInt(user.ID, 10)
	buf, _ := json.Marshal(user)
	_, err := a.rds.Do("set", key, string(buf), "ex", config.Web.ExpiresAt*3600)
	return err
}
