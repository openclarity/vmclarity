package utils

import "testing"

func TestContains_strings(t *testing.T) {
	type args struct {
		s []string
		v string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "nil slice",
			args: args{
				s: nil,
				v: "str",
			},
			want: false,
		},
		{
			name: "empty slice",
			args: args{
				s: []string{},
				v: "str",
			},
			want: false,
		},
		{
			name: "empty value",
			args: args{
				s: []string{"str", "str1"},
				v: "",
			},
			want: false,
		},
		{
			name: "string slice - contains",
			args: args{
				s: []string{"str", "str1"},
				v: "str",
			},
			want: true,
		},
		{
			name: "string slice - contains twice",
			args: args{
				s: []string{"str", "str1", "str"},
				v: "str",
			},
			want: true,
		},
		{
			name: "string slice - does not contains",
			args: args{
				s: []string{"str", "str1"},
				v: "str2",
			},
			want: false,
		},
		//{
		//	name: "nil value",
		//	args: args{
		//		s: []string{"str", "str1"},
		//		v: nil,
		//	},
		//	want: false,
		//},
		//{
		//	name: "int slice - contains",
		//	args: args{
		//		s: []string{1, 2},
		//		v: 1,
		//	},
		//	want: true,
		//},
		//{
		//	name: "int slice - contains twice",
		//	args: args{
		//		s: []string{1, 2, 1},
		//		v: 1,
		//	},
		//	want: true,
		//},
		//{
		//	name: "int slice - does not contains",
		//	args: args{
		//		s: []string{1, 2},
		//		v: 3,
		//	},
		//	want: false,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.s, tt.args.v); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContains_int(t *testing.T) {
	type args struct {
		s []int
		v int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "nil slice",
			args: args{
				s: nil,
				v: 1,
			},
			want: false,
		},
		{
			name: "empty slice",
			args: args{
				s: []int{},
				v: 1,
			},
			want: false,
		},
		{
			name: "int slice - contains",
			args: args{
				s: []int{1, 2},
				v: 1,
			},
			want: true,
		},
		{
			name: "int slice - contains twice",
			args: args{
				s: []int{1, 2, 1},
				v: 1,
			},
			want: true,
		},
		{
			name: "int slice - does not contains",
			args: args{
				s: []int{1, 2},
				v: 3,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.s, tt.args.v); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
