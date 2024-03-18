package local

import (
	"errors"
	"github.com/openclarity/vmclarity/scanner/types"
	"gorm.io/gorm"
)

type repo[T any] struct {
	DB    *gorm.DB
	Model T
}

type getParams struct {
	filters []interface{}
}

func (r *repo[T]) GetAll(params getParams, dest *[]T) error {
	return r.DB.Model(r.Model).Find(dest, params.filters...).Error
}

func (r *repo[T]) Get(cond *T, dest *T) error {
	err := r.DB.Where(cond).First(dest).Error
	return extractErr(err)
}

func (r *repo[T]) Update(cond *T, updatedColumns *T) error {
	err := r.DB.Model(r.Model).Select("*").Where(cond).Updates(updatedColumns).Error
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
		return types.ErrNotFound
	}
	return err
}
