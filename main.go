package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf16"

	"golang.org/x/sys/windows"
)

func parseDrives() ([]string, error) {
	n, err := windows.GetLogicalDriveStrings(0, nil)
	if err != nil {
		return nil, err
	}
	buf := make([]uint16, n)
	_, err = windows.GetLogicalDriveStrings(n, &buf[0])
	if err != nil {
		return nil, err
	}
	s := string(utf16.Decode(buf))

	return strings.Split(strings.TrimRight(s, "\x00"), "\x00"), nil
}

func getDrives(drives []string) ([]string, []string) {
	var encoded []uint16
	var removables []string
	var fixeds []string
	for _, d := range drives {
		encoded = utf16.Encode([]rune(d))
		fmt.Printf("a: %q\n", encoded)
		encoded = append(encoded, 0)

		switch driveType := windows.GetDriveType(&encoded[0]); driveType {
		case windows.DRIVE_REMOVABLE:
			fmt.Printf("drive %s is removable\n", d)
			removables = append(removables, d)
		case windows.DRIVE_FIXED:
			fmt.Printf("drive %s is fixed\n", d)
			fixeds = append(fixeds, d)
		}
	}

	return removables, fixeds
}

func scanRemovables(walkdir string) {
	err := filepath.Walk(walkdir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			var e *fs.PathError
			if errors.As(err, &e) {
				if strings.HasSuffix(e.Path, `:\System Volume Information`) {
					return nil
				}
			}
			return err
		}

		fmt.Printf("filepath: %s\n", path)

		if !info.IsDir() && isVirus(path) {
			fmt.Println("**************************")
			fmt.Printf("FOUND VIRUS: %s\n", path)
			fmt.Println("**************************")
			err = os.Remove(path)
			if err != nil {
				return err
			}
			fmt.Println("**************************")
			fmt.Printf("DELETED %s\n", path)
			fmt.Println("**************************")
			return nil // not necessery to fix attrs because we deleted it
		}
		err = fixAttrs(path)
		if err != nil {
			fmt.Printf("after fix attrs err: %v\n", err)
			return err
		}

		return nil
	})
	if err != nil {
		log.Printf("error while walking: err:%+#v ---- %T\n ", err, err)
	}

}

func fixAttrs(fPath string) error {
	ePath := utf16.Encode([]rune(fPath))
	ePath = append(ePath, 0)
	fa, err := windows.GetFileAttributes(&ePath[0])
	if err != nil {
		return err
	}
	fmt.Printf("ePath: %v\n", fa)
	if fa == windows.FILE_ATTRIBUTE_DIRECTORY {
		// skipping this path for bootable removables(0x57 - syscall.Errno)
		return nil
	}

	attrs := windows.FILE_ATTRIBUTE_READONLY | windows.FILE_ATTRIBUTE_HIDDEN | windows.FILE_ATTRIBUTE_SYSTEM

	fa = fa &^ uint32(attrs)
	err = windows.SetFileAttributes(&ePath[0], fa)
	if err != nil {
		return err
	}

	return nil
}

func getDirs() []string {
	var dirs []string
	h, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("err while getting home")
	}
	roaming := "AppData/Roaming"
	start := "AppData/Roaming/Microsoft/Windows/Start Menu/Programs/Startup"
	r := filepath.Join(h, roaming)
	s := filepath.Join(h, start)
	dirs = append(dirs, r, s)

	return dirs
}

func scanDirs(dirs []string) {
	for _, dir := range dirs {
		fmt.Printf("scanning the dir %s\n", dir)
		entries, err := os.ReadDir(dir)
		if err != nil {
			log.Fatalf("error while reading the dir %s, err:%v\n", dir, err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			n := entry.Name()
			if !isVirus(n) {
				continue
			}

			p := filepath.Join(dir, n)
			fmt.Printf("virus FOUND!! -- %s\n", p)
			err := os.Remove(p)
			if err != nil {
				log.Fatalf("error while removing. err:%v\n", err)
			}

		}
	}
}

func isVirus(name string) bool {
	return strings.HasSuffix(name, ".js")
}

func main() {
	drives, err := parseDrives()
	if err != nil {
		log.Fatalf("error while parsing drives. err:%v\n", err)
	}

	removables, fixeds := getDrives(drives)
	fmt.Printf("removables: %+v --- fixeds: %+v\n", removables, fixeds)
	for _, removable := range removables {
		scanRemovables(removable)
	}

	dirs := getDirs()
	fmt.Printf("dirs: %+v\n", dirs)
	scanDirs(dirs)

}
