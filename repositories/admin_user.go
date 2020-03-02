package repositories

import (
	"errors"
	"github.com/article-publish-server1/datamodels"
	"github.com/article-publish-server1/datasorces"
)

type AdminUserRepository interface {
	Create(user *datamodels.AdminUser) error
	Get(query interface{}, args ...interface{}) (*datamodels.AdminUser, error)
}

type adminUserRepository struct {
	commonRepository
}

func NewAdminUserRepository() AdminUserRepository {
	return &adminUserRepository{
		commonRepository: commonRepository{
			db: datasorces.PqDB.DB,
		},
	}
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
