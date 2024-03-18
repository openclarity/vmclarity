package local

import (
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
	return r.DB.Where(cond).First(dest).Error
}

func (r *repo[T]) Update(cond *T, updatedColumns *T) error {
	return r.DB.Model(r.Model).Select("*").Where(cond).Updates(updatedColumns).Error
}

func (r *repo[T]) Delete(cond *T) error {
	if err := r.DB.Model(r.Model).Delete(cond); err != nil {
		return err.Error
	}
	return nil
}

func (r *repo[T]) Create(data *T) error {
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
