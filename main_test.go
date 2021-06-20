package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"
	"unicode/utf16"

	"golang.org/x/sys/windows"
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

func Test_app(t *testing.T) {
	p := "."
	tmpDir, err := os.MkdirTemp(".", "shortcut-log")
	if err != nil {
		t.Errorf("error creating temp directory: %v\n", err)
	}
	fmt.Printf("tempDir: %+v\n", tmpDir)

	err = os.MkdirAll(filepath.Join(tmpDir, p), 0755)
	defer os.RemoveAll(tmpDir)
	if err != nil {
		t.Errorf("err: %v\n", err)
	}

	f, err := os.CreateTemp(tmpDir, "temp*.js")
	if err != nil {
		t.Errorf("error while creating file: %v\n", err)
	}
	f.Close()

	makeHidden(f.Name())
	time.Sleep(20 * time.Second)
	want := 0
	viruses, err := scanRemovables(tmpDir)
	if err != nil {
		t.Errorf("inside scanRemovables(): err: %v\n", err)
	}
	got := len(viruses)
	if got > 0 {
		_ = removeViruses(viruses)
	}

	viruses, err = scanRemovables(tmpDir)
	if err != nil {
		t.Errorf("inside scanRemovables(): err: %v\n", err)
	}
	got = len(viruses)
	if got != want {
		t.Errorf("got: %d, want: %d\n", got, want)
	}
}

func getAttrs(fPath string) {
	ePath, err := syscall.UTF16PtrFromString(fPath)
	if err != nil {
		log.Fatalln("error while convertint to ptr: ", err)
	}
	fa, err := windows.GetFileAttributes(ePath)
	if err != nil {
		log.Fatalf("error while getting attrs: %v\n", err)
	}
	fmt.Printf("fa: %v\n", fa)
}

func setAttrs(fPath string) {
	ePath := utf16.Encode([]rune(fPath))
	ePath = append(ePath, 0)
	attrs := windows.FILE_ATTRIBUTE_READONLY | windows.FILE_ATTRIBUTE_HIDDEN | windows.FILE_ATTRIBUTE_SYSTEM
	fmt.Printf("attrs: %#X\n", uint32(attrs))
	err := windows.SetFileAttributes(&ePath[0], uint32(attrs))
	if err != nil {
		log.Fatalln("error while setting attrs: ", err)
	}
}

func makeHidden(fPath string) {
	fmt.Printf("fPath: %s\n", fPath)
	getAttrs(fPath)

	setAttrs(fPath)

	fmt.Println("after setting hidden")
	getAttrs(fPath)

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
}
