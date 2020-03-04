package repositories

import (
	"github.com/article-publish-server1/datamodels"
	"github.com/article-publish-server1/datasorces"
)

type TagRepository interface {
	Create(tag *datamodels.Tag) error
	Remove(query interface{}, args ...interface{}) error
	Update(docs map[string]interface{}, query interface{}, args ...interface{}) error
	Get(query interface{}, args ...interface{}) (*datamodels.Tag, error)
	ListAll(order, query interface{}, args ...interface{}) ([]datamodels.Tag, error)
}

type tagRepository struct {
	commonRepository
}

func NewTagRepository() TagRepository {
	return &tagRepository{
		commonRepository: commonRepository{
			db: datasorces.PqDB.DB,
		},
	}
}

func (t tagRepository) Create(tag *datamodels.Tag) error {
	return t.commonRepository.Create(tag)
}

func (t tagRepository) Remove(query interface{}, args ...interface{}) error {
	return t.commonRepository.Remove(&datamodels.Tag{}, query, args...)
}

func (t tagRepository) Update(docs map[string]interface{}, query interface{}, args ...interface{}) error {
	return t.commonRepository.Update(&datamodels.Tag{}, query, docs, args)
}

func (t tagRepository) Get(query interface{}, args ...interface{}) (*datamodels.Tag, error) {
	record, err := t.commonRepository.Get(&datamodels.Tag{}, query, args...)
	if record == nil {
		return nil, err
	}

	return record.(*datamodels.Tag), err
}

func (t tagRepository) ListAll(order, query interface{}, args ...interface{}) ([]datamodels.Tag, error) {
	var list []datamodels.Tag

	if err := t.commonRepository.ListAll(&list, order, query, args...); err != nil {
		return nil, err
	}

	return list, nil
}
