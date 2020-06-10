package cache

import (
	"reflect"
	"testing"
)

func TestBytesToString(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"#1", args{b: []byte("hello")}, "hello"},
		{"#2", args{b: []byte("AAAAAA")}, "AAAAAA"},
		{"#3", args{b: []byte("_*")}, "_*"},
		{"#4", args{b: []byte("1")}, "1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesToString(tt.args.b); got != tt.want {
				t.Errorf("BytesToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringToBytes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"#1", args{s: "hello"}, []byte("hello")},
		{"#2", args{s: "AAA"}, []byte("AAA")},
		{"#3", args{s: "_*"}, []byte("_*")},
		{"#4", args{s: "1"}, []byte("1")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringToBytes(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
