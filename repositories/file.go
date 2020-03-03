package repositories

import (
	"github.com/article-publish-server1/datamodels"
	"github.com/article-publish-server1/datasorces"
	"github.com/jinzhu/gorm"
)

type FileRepository interface {
	Create(file *datamodels.File) error
	Get(query interface{}, args ...interface{}) (*datamodels.File, error)
	CreateWithTx(tx *gorm.DB, file *datamodels.File) error
}

type fileRepository struct {
	commonRepository
}

func NewFileRepository() FileRepository {
	return &fileRepository{
		commonRepository: commonRepository{
			db: datasorces.PqDB.DB,
		},
	}
}

func (f fileRepository) Create(file *datamodels.File) error {
	return f.commonRepository.Create(file)
}

func (f fileRepository) Get(query interface{}, args ...interface{}) (*datamodels.File, error) {
	record, err := f.commonRepository.Get(&datamodels.File{}, query, args...)
	if record == nil {
		return nil, err
	}

	return record.(*datamodels.File), err
}

func (f fileRepository) CreateWithTx(tx *gorm.DB, file *datamodels.File) error {
	return f.commonRepository.CreateWithTx(tx, file)
}
