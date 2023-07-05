package models

import (
	"log"
	"encoding/json"
)

type nullable[T any] []*T

func (n *nullable[T]) UnmarshalJSON(data []byte) error {
	log.Print("called with ", string(data))
	var val *T
	err := json.Unmarshal(data, &val)
	if err != nil {
		return err
	}
	*n = append((*n)[0:0], val)
	return nil
}

func (n nullable[T]) MarshalJSON() ([]byte, error) {
	log.Print("State ", n)
	return json.Marshal(n[0])
}

func (n nullable[T]) IsPresent() bool {
	return n != nil
}

func (n nullable[T]) IsNull() bool {
	if len(n) == 0 {
		return true
	}
	return n[0] == nil
}

func (n nullable[T]) GetValue() T {
	if len(n) == 0 || n[0] == nil {
		var v T
		return v
	}
	return *n[0]
}

func Null[T any]() nullable[T] {
	return nullable[T]{nil}
}

func Nullable[T any](val T) nullable[T] {
	return nullable[T]{&val}
}

