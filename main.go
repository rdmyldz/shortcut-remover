package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf16"

	"golang.org/x/sys/windows"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

const warning = `
***********************************************************************************************
*     if you have javascript files(end with '.js' ) in your removables(e.g. USB flash drives) *
* all will be deleted!!!                                                                      *
*                                                                                             *
*                                                                                             *
*				written by @rdmyldz                                           *
***********************************************************************************************
`

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

func getDrives(drives []string) ([]string, []string, error) {
	var removables []string
	var fixeds []string
	for _, d := range drives {
		encoded, err := windows.UTF16PtrFromString(d)
		if err != nil {
			return nil, nil, err
		}
		switch driveType := windows.GetDriveType(encoded); driveType {
		case windows.DRIVE_REMOVABLE:
			removables = append(removables, d)
		case windows.DRIVE_FIXED:
			fixeds = append(fixeds, d)
		}
	}

	return removables, fixeds, nil
}

func scanRemovables(walkdir string) ([]string, error) {
	var viruses []string
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

		if walkdir == path {
			return nil
		}
		if !info.IsDir() && isVirus(path) {
			viruses = append(viruses, path)
			return nil // not necessery to fix attrs because we'll delete it
		}
		err = fixAttrs(path)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error while walking: %w", err)
	}
	return viruses, nil

}

func fixAttrs(fPath string) error {
	ePath, err := windows.UTF16PtrFromString(fPath)
	if err != nil {
		return err
	}

	fa, err := windows.GetFileAttributes(ePath)
	if err != nil {
		return err
	}

	attrs := windows.FILE_ATTRIBUTE_READONLY | windows.FILE_ATTRIBUTE_HIDDEN | windows.FILE_ATTRIBUTE_SYSTEM

	fa = fa &^ uint32(attrs)
	err = windows.SetFileAttributes(ePath, fa)
	if err != nil {
		return err
	}

	return nil
}

func getDirs() ([]string, error) {
	var dirs []string
	h, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("err while getting homeDir: %w", err)
	}
	roaming := "AppData/Roaming"
	start := "AppData/Roaming/Microsoft/Windows/Start Menu/Programs/Startup"
	r := filepath.Join(h, roaming)
	s := filepath.Join(h, start)
	dirs = append(dirs, r, s)

	return dirs, err
}

func scanDirs(dirs []string) ([]string, error) {
	var viruses []string
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("error while reading the dir %s, err:%w", dir, err)
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
			viruses = append(viruses, p)
		}
	}
	return viruses, nil
}

func isVirus(name string) bool {
	return strings.HasSuffix(name, ".js")
}

func removeViruses(viruses []string) []error {
	var errs []error

	for i, path := range viruses {
		err := os.Remove(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("err-%d: while removing %s, err: %v", i, path, err))
		}
	}
	return errs
}

func getTmpFile() (*os.File, error) {
	fTime := time.Now().Format("20060102-150405")

	tmpDir, err := os.MkdirTemp("", "shortcut-"+fTime+"_*")
	if err != nil {
		return nil, fmt.Errorf("error while creating temp directory: %w", err)
	}

	f, err := os.CreateTemp(tmpDir, fTime+"_*.txt")
	if err != nil {
		return nil, fmt.Errorf("error while creating file: %w", err)
	}
	return f, nil
}

func warnUser() error {
	fmt.Println(warning)

	r := bufio.NewReader(os.Stdin)
	// fmt.Printf("want to continue? [y[es]/n[o]]:")
	fmt.Print("want to continue? [y[es]/n[o]]:")
	input, err := r.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error while reading input: %w", err)
	}
	input = strings.ToLower(strings.TrimSpace(input))
	if input == "y" || input == "yes" {
		fmt.Println("\nscanning is starting...")
		return nil
	}
	fmt.Printf("You typed '%s' if you wanna scan your device, type 'y' or 'yes'", input)
	fmt.Println("\nexiting...")
	return fmt.Errorf("the user said 'no'")

}

func printInfo(logger *log.Logger, format string, a ...interface{}) {
	fmt.Printf(format, a...)
	logger.Printf(format, a...)
}

func main() {
	tmpFile, err := getTmpFile()
	if err != nil {
		log.Fatalf("inside getTmpFile(): %v\n", err)
	}
	defer tmpFile.Close()

	InfoLogger = log.New(tmpFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(tmpFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	err = warnUser()
	if err != nil {
		ErrorLogger.Fatalln(err)
	}

	drives, err := parseDrives()
	if err != nil {
		ErrorLogger.Fatalf("inside parseDrives(): %v\n", err)
	}

	removables, fixeds, err := getDrives(drives)
	if err != nil {
		ErrorLogger.Fatalf("inside getDrives(): %v\n", err)
	}
	printInfo(InfoLogger, "removables: %+v --- fixeds: %+v\n", removables, fixeds)
	var found []string
	for _, removable := range removables {
		v, err := scanRemovables(removable)
		if err != nil {
			ErrorLogger.Printf("inside scanRemovables(): removable: %s,err: %v\n", removable, err)
		}
		found = append(found, v...)
	}

	dirs, err := getDirs()
	if err != nil {
		ErrorLogger.Fatalf("inside getDirs(): %v\n", err)
	}

	v, err := scanDirs(dirs)
	if err != nil {
		ErrorLogger.Fatalf("inside scanDirs(): %v\n", err)
	}

	found = append(found, v...)
	if len(found) > 0 {
		printInfo(InfoLogger, "%d viruses found: %v\n", len(found), found)
		errs := removeViruses(found)
		if len(errs) > 0 {
			ErrorLogger.Printf("returned from removeViruses(): %v\n", errs)
			for _, err := range errs {
				log.Println(err)
			}

		} else {
			printInfo(InfoLogger, "viruses got removed...\n")
		}
	} else {
		printInfo(InfoLogger, "not found anything\n")
	}

	fmt.Printf("press 'enter' to exit...")
	fmt.Scanln()
}
