package local

import (
	"errors"
	"github.com/openclarity/vmclarity/scanner/server/store"
	"gorm.io/gorm"
)

type repo[T any] struct {
	DB    *gorm.DB
	Model T
}

type getParams struct {
	filters [][]interface{}
}

func (r *repo[T]) GetAll(params getParams, dest *[]T) error {
	tx := r.DB.Model(r.Model)
	for _, filter := range params.filters {
		if len(filter) >= 1 {
			tx = tx.Where(filter[0], filter[1:]...)
		}
	}
	return tx.Find(dest).Error
}

func (r *repo[T]) Get(cond *T, dest *T) error {
	err := r.DB.Where(cond).First(dest).Error
	return extractErr(err)
}

func (r *repo[T]) Update(cond *T, updatedColumns *T) error {
	err := r.DB.Model(r.Model).Where(cond).Updates(updatedColumns).Error
	return extractErr(err)
}

func (r *repo[T]) Delete(cond *T) error {
	err := r.DB.Model(r.Model).Delete(cond)
	return extractErr(err.Error)
}

func (r *repo[T]) Create(data *T) error {
	return r.DB.Create(data).Error
}

func (r *repo[T]) CreateMany(data *[]*T) error {
	return r.DB.Create(data).Error
}

func (r *repo[T]) getTx() *gorm.DB {
	return r.DB.Model(r.Model)
}

func newRepo[T any](db *gorm.DB, model T) *repo[T] {
	return &repo[T]{
		DB:    db,
		Model: model,
	}
}

func extractErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return store.ErrNotFound
	}
	return err
}
