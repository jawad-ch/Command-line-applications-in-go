package main

import (
	"os"
	"testing"
)

type testCase struct {
	name     string
	files    []string
	exts     []string
	minSize  int64
	expected bool
}

func TestFilterOut(t *testing.T) {
	testCases := []testCase{
		{"FilterNoExtension", []string{"testdata/dir.log"}, []string{}, 0, false},
		{"FilterExtensionMatch", []string{"testdata/dir.log"}, []string{".log"}, 0, false},
		{"FilterExtensionNoMatch", []string{"testdata/dir.log"}, []string{".sh"}, 0, true},
		{"FilterExtensionSizeMatch", []string{"testdata/dir.log"}, []string{".log"}, 10, false},
		{"FilterExtensionSizeNoMatch", []string{"testdata/dir.log"}, []string{".log"}, 20, true},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, file := range tc.files {
				info, err := os.Stat(file)
				if err != nil {
					t.Fatal(err)
				}

				f := filterOut(file, tc.exts, tc.minSize, info)

				if f != tc.expected {
					t.Errorf("Expected '%t', got '%t' instead at index = %d\n", tc.expected, f, i)
				}
			}
		})
	}
}
