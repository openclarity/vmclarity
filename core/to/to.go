// Copyright © 2024 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package to

import (
	"cmp"
	"sort"
)

func Ptr[T any](value T) *T {
	return &value
}

// PtrOrNil returns a pointer to t if it has a non-empty value otherwise nil.
func PtrOrNil[T comparable](t T) *T {
	var empty T
	if t == empty {
		return nil
	}

	return &t
}

// ValueOrZero returns the value that the pointer ptr pointers to. It returns
// the zero value if ptr is nil.
func ValueOrZero[T any](ptr *T) T {
	var t T
	if ptr != nil {
		t = *ptr
	}

	return t
}

// Keys returns a slice of keys from m map.
func Keys[K comparable, V any](m map[K]V) []K {
	s := make([]K, 0, len(m))
	for k := range m {
		s = append(s, k)
	}

	return s
}

// SortedKeys returns a slice of keys from m map in sorted order (ascending).
func SortedKeys[K cmp.Ordered, V any](m map[K]V) []K {
	s := make([]K, 0, len(m))
	for k := range m {
		s = append(s, k)
	}

	sort.Slice(s, func(i, j int) bool {
		return cmp.Compare[K](s[i], s[j]) < 0
	})

	return s
}

// Values returns a slice of values from m map.
func Values[K comparable, V any](m map[K]V) []V {
	s := make([]V, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}

	return s
}

// UniqueSlice returns a slice without duplicate elements.
func UniqueSlice[T comparable](items []T) []T {
	var filtered []T
	unique := make(map[T]bool, len(items))
	for _, item := range items {
		if !unique[item] {
			filtered = append(filtered, item)
			unique[item] = true
		}
	}
	return filtered
}

// UniqueSliceByKey returns a slice without duplicate elements using a custom get key function.
func UniqueSliceByKey[T any](items []T, getKey func(T) string) []T {
	var filtered []T

	unique := make(map[string]bool, len(items))
	for _, item := range items {
		if key := getKey(item); key != "" && !unique[key] {
			filtered = append(filtered, item)
			unique[key] = true
		}
	}

	return filtered
}
