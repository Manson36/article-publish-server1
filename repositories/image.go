package repositories

import (
	"github.com/article-publish-server1/datamodels"
	"github.com/article-publish-server1/datasorces"
	"github.com/jinzhu/gorm"
)

type ImageRepository interface {
	Create(image *datamodels.Image) error
	CreateWithTx(tx *gorm.DB, iamge *datamodels.Image) error
	List(order, limit, offset, query interface{}, args ...interface{}) ([]datamodels.Image, error)
	Count(interface{}, ...interface{}) (int64, error)
	Remove(interface{}, ...interface{}) error
}

type imageRepository struct {
	commonRepository
}

func NewImageRepository() ImageRepository {
	return &imageRepository{
		commonRepository: commonRepository{
			db: datasorces.PqDB.DB,
		},
	}
}

func (i imageRepository) Create(image *datamodels.Image) error {
	return i.commonRepository.Create(image)
}

func (i imageRepository) CreateWithTx(tx *gorm.DB, image *datamodels.Image) error {
	return i.commonRepository.CreateWithTx(tx, image)
}

func (i imageRepository) List(order, limit, offset, query interface{}, args ...interface{}) ([]datamodels.Image, error) {
	var list []datamodels.Image
	if err := i.commonRepository.List(&list, order, limit, offset, query, args); err != nil {
		return nil, err
	}

	return list, nil
}

func (i imageRepository) Count(query interface{}, args ...interface{}) (int64, error) {
	return i.commonRepository.Count(&datamodels.Image{}, query, args...)
}

func (i imageRepository) Remove(query interface{}, args ...interface{}) error {
	return i.commonRepository.Remove(&datamodels.Image{}, query, args...)
}
