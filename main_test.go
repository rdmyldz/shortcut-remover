package main

import (
	"testing"
)

func Test_scanRemovables(t *testing.T) {
	type args struct {
		walkdir string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{name: "1st", args: args{walkdir: `d:\`}},
		{name: "2nd", args: args{walkdir: `e:\`}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanRemovables(tt.args.walkdir)
		})
	}
}

func Test_getDirs(t *testing.T) {
	_ = getDirs()
}

func Test_fixAttrs(t *testing.T) {
	type args struct {
		fPath string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixAttrs(tt.args.fPath)
		})
	}
	// scanRemovables(`C:\Users\rdm\Desktop\repos\github.com\rdmyldz\shortcut-virus\temp`)
	scanRemovables(`d:\`)
}
