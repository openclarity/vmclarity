package utils

import (
	"testing"

	"github.com/openclarity/vmclarity/shared/pkg/utils"
)

func TestTrimMountPath(t *testing.T) {
	type args struct {
		toTrim    string
		mountPath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "root path",
			args: args{
				toTrim:    "/mnt",
				mountPath: "/mnt",
			},
			want: "/",
		},
		{
			name: "not root path",
			args: args{
				toTrim:    "/mnt/foo",
				mountPath: "/mnt",
			},
			want: "/foo",
		},
		{
			name: "no trim",
			args: args{
				toTrim:    "/bar/foo",
				mountPath: "/mnt",
			},
			want: "/bar/foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimMountPath(tt.args.toTrim, tt.args.mountPath); got != tt.want {
				t.Errorf("TrimMountPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveMountPathSubStringIfNeeded(t *testing.T) {
	type args struct {
		toTrim    string
		mountPath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "from prefix",
			args: args{
				toTrim:    "/mnt/foo some massage",
				mountPath: "/mnt",
			},
			want: "/foo some massage",
		},
		{
			name: "from prefix replace with /",
			args: args{
				toTrim:    "/mnt some massage",
				mountPath: "/mnt",
			},
			want: "/ some massage",
		},
		{
			name: "from middle",
			args: args{
				toTrim:    "some message /mnt/foo some massage",
				mountPath: "/mnt",
			},
			want: "some message /foo some massage",
		},
		{
			name: "from middle replace with /",
			args: args{
				toTrim:    "some message /mnt some massage",
				mountPath: "/mnt",
			},
			want: "some message / some massage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveMountPathSubStringIfNeeded(tt.args.toTrim, tt.args.mountPath); got != tt.want {
				t.Errorf("RemoveMountPathSubStringIfNeeded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShouldStripInputPath(t *testing.T) {
	type args struct {
		inputShouldStrip  *bool
		familyShouldStrip bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "inputShouldStrip = nil, familyShouldStrip = false",
			args: args{
				inputShouldStrip:  nil,
				familyShouldStrip: false,
			},
			want: false,
		},
		{
			name: "inputShouldStrip = nil, familyShouldStrip = true",
			args: args{
				inputShouldStrip:  nil,
				familyShouldStrip: true,
			},
			want: true,
		},
		{
			name: "inputShouldStrip = false, familyShouldStrip = false",
			args: args{
				inputShouldStrip:  utils.PointerTo(false),
				familyShouldStrip: false,
			},
			want: false,
		},
		{
			name: "inputShouldStrip = false, familyShouldStrip = true",
			args: args{
				inputShouldStrip:  utils.PointerTo(false),
				familyShouldStrip: true,
			},
			want: false,
		},
		{
			name: "inputShouldStrip = true, familyShouldStrip = false",
			args: args{
				inputShouldStrip:  utils.PointerTo(true),
				familyShouldStrip: false,
			},
			want: true,
		},
		{
			name: "inputShouldStrip = true, familyShouldStrip = true",
			args: args{
				inputShouldStrip:  utils.PointerTo(true),
				familyShouldStrip: true,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShouldStripInputPath(tt.args.inputShouldStrip, tt.args.familyShouldStrip); got != tt.want {
				t.Errorf("ShouldStripInputPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
