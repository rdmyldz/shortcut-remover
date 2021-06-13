package main

import (
	"fmt"
	"strings"
	"unicode/utf16"

	"golang.org/x/sys/windows"
)

func parseDrives() ([]string, error) {
	n, err := windows.GetLogicalDriveStrings(0, nil)
	if err != nil {
		return nil, err
	}
	fmt.Printf("len of n:%d\n", n)
	buf := make([]uint16, n)
	_, err = windows.GetLogicalDriveStrings(n, &buf[0])
	if err != nil {
		return nil, err
	}
	fmt.Printf("buf: %q --len of buf %d\n", buf, len(buf))
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
		fmt.Println("returned", windows.GetDriveType(&encoded[0]))
		if windows.DRIVE_REMOVABLE == windows.GetDriveType(&encoded[0]) {
			fmt.Printf("drive %s is removable\n", d)
			removables = append(removables, d)
		}
		if windows.DRIVE_FIXED == windows.GetDriveType(&encoded[0]) {
			fmt.Printf("drive %s is fixed\n", d)
			fixeds = append(fixeds, d)
		}
	}

	return removables, fixeds
}

func main() {
	drives, err := parseDrives()
	if err != nil {
		panic(err)
	}

	removables, fixeds := getDrives(drives)
	fmt.Printf("removables: %+v --- fixeds: %+v\n", removables, fixeds)

}
