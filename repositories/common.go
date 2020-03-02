package repositories

import (
	"github.com/jinzhu/gorm"
)

type CommonRepository interface {
	Create(entity interface{}) error
	Save(entity interface{}) error
	Remove(entity, query interface{}, args ...interface{}) error
	Update(entity, query interface{}, docs map[string]interface{}, args ...interface{}) error
	Get(entity, query interface{}, args ...interface{}) (interface{}, error)
	List(entities, order, limit, offset, query interface{}, args ...interface{}) error
	ListAll(entities, order, query interface{}, args ...interface{}) error
	Count(entity, query interface{}, args ...interface{}) (int64, error)
}

type commonRepository struct {
	db *gorm.DB
}

func (c commonRepository) Create(entity interface{}) error {
	return c.db.Create(entity).Error
}

func (c commonRepository) Save(entity interface{}) error {
	return c.db.Save(entity).Error
}

func (c commonRepository) Remove(entity, query interface{}, args ...interface{}) error {
	return c.db.Model(entity).Where(query, args...).Updates(map[string]interface{}{
		"removed":    true,
		"removed_at": gorm.Expr("now()"),
	}).Error
}

func (c commonRepository) Update(entity, query interface{}, docs map[string]interface{}, args ...interface{}) error {
	return c.db.Model(entity).Where(query, args...).Update(docs).Error
}

func (c commonRepository) Get(entity, query interface{}, args ...interface{}) (interface{}, error) {
	if err := c.db.Where(query, args...).Take(entity).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
		return nil, nil
	}

	return entity, nil
}

func (c commonRepository) List(entities, order, limit, offset, query interface{}, args ...interface{}) error {
	return c.db.Order(order).Limit(limit).Offset(offset).Where(query, args...).Find(entities).Error
}

func (c commonRepository) ListAll(entities, order, query interface{}, args ...interface{}) error {
	return c.List(entities, order, -1, -1, query, args...)
}

func (c commonRepository) Count(entity, query interface{}, args ...interface{}) (int64, error) {
	var total int64
	if err := c.db.Model(entity).Where(query, args...).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
