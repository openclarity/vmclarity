package models

import (
	"testing"
	"encoding/json"

	"github.com/google/go-cmp/cmp"
)

type Document struct {
	NullableString nullable[string] `json:"nullableString,omitempty"`
}

func TestNullable_MarshalJSON(t *testing.T) {
	someValue := "foo"

	tests := []struct {
		name string
		doc Document
		want string
	}{
		{
			name: "unset field should not be sent in the json",
			doc: Document{},
			want: "{}",
		},
		{
			name: "nil nullable should send null",
			doc: Document{
				NullableString: Null[string](),
			},
			want: "{\"nullableString\":null}",
		},
		{
			name: "if value set should send value",
			doc: Document{
				NullableString: Nullable[string](someValue),
			},
			want: "{\"nullableString\":\"foo\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.doc)
			if err != nil {
				t.Fatalf("unexpected error marshalling doc: %v", err)
			}
			got := string(b)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Nullable MarshalJSON mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNullable_UnmarshalJSON(t *testing.T) {
	someValue := "foo"

	tests := []struct {
		name string
		input string
		want Document
		IsPresent bool
		IsNull bool
		GetValue string
	}{
		{
			name: "unset field should not be set in the document",
			input: "{}",
			want: Document{},
			IsPresent: false,
			IsNull: true,
		},
		{
			name: "null field should be set in the document with value nil",
			input: "{\"nullableString\":null}",
			want: Document{
				NullableString: Null[string](),
			},
			IsPresent: true,
			IsNull: true,
		},
		{
			name: "value field should be set in the document with value",
			input: "{\"nullableString\":\"foo\"}",
			want: Document{
				NullableString: Nullable[string](someValue),
			},
			IsPresent: true,
			IsNull: false,
			GetValue: someValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Document
			err := json.Unmarshal([]byte(tt.input), &got)
			if err != nil {
				t.Fatalf("unexpected error marshalling doc: %v", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Nullable MarshalJSON mismatch (-want +got):\n%s", diff)
			}

			isSet := got.NullableString.IsPresent()
			if diff := cmp.Diff(tt.IsPresent, isSet); diff != "" {
				t.Errorf("Nullable IsPresent mismatch (-want +got):\n%s", diff)
			}

			isNull := got.NullableString.IsNull()
			if diff := cmp.Diff(tt.IsNull, isNull); diff != "" {
				t.Errorf("Nullable IsNull mismatch (-want +got):\n%s", diff)
			}

			getValue := got.NullableString.GetValue()
			if diff := cmp.Diff(tt.GetValue, getValue); diff != "" {
				t.Errorf("Nullable GetValue mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
