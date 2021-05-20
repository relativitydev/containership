package utils

import (
	"testing"
)

func TestAfter(t *testing.T) {
	type args struct {
		value string
		a     string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "slash",
			args: args{
				value: "thisis/atest",
				a:     "/",
			},
			want: "atest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := After(tt.args.value, tt.args.a); got != tt.want {
				t.Errorf("After() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBefore(t *testing.T) {
	type args struct {
		value string
		a     string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "slash",
			args: args{
				value: "thisis/atest",
				a:     "/",
			},
			want: "thisis",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Before(tt.args.value, tt.args.a); got != tt.want {
				t.Errorf("Before() = %v, want %v", got, tt.want)
			}
		})
	}
}
