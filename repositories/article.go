package repositories

import "github.com/article-publish-server1/datamodels"

type ArticleRepository interface {
	Create(article *datamodels.Article) error
	Remove(query interface{}, args ...interface{}) error
	Update(doc map[string]interface{}, query interface{}, args ...interface{}) error
	Get(query interface{}, args ...interface{}) (*datamodels.Article, error)
	List(order, limit, offset, query interface{}, args ...interface{}) ([]datamodels.Article, error)
	Count(query interface{}, args ...interface{}) (int64, error)
	Save(article *datamodels.Article) error
}
