package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		cfg      config
		expected string
	}{
		{
			name: "NoFilter",
			root: "testdata",
			cfg: config{
				size: 0,
				list: true,
			},
			expected: "testdata/dir.log\ntestdata/dir2/markup.html\ntestdata/dir2/script.sh\n",
		},
		{
			name: "FilterExtensionMatch",
			root: "testdata",
			cfg: config{
				exts: []string{".log"},
				size: 0,
				list: true,
			},
			expected: "testdata/dir.log\n",
		},
		{
			name: "FilterExtensionSizeMatch",
			root: "testdata",
			cfg: config{
				exts: []string{".log"},
				size: 10,
				list: true,
			},
			expected: "testdata/dir.log\n",
		},
		{
			name: "FilterExtensionSizeNoMatch",
			root: "testdata",
			cfg: config{
				exts: []string{".log"},
				size: 20,
				list: true,
			},
			expected: "",
		},
		{
			name: "FilterExtensionNoMatch",
			root: "testdata",
			cfg: config{
				exts: []string{".gz"},
				size: 0,
				list: true,
			},
			expected: "",
		},
		{
			name: "FilterExtensionMixMatch",
			root: "testdata",
			cfg: config{
				exts: []string{".log", ".html"},
				size: 0,
				list: true,
			},
			expected: "testdata/dir.log\ntestdata/dir2/markup.html\n",
		},
		{
			name: "FilterExtensionMixNoMatch",
			root: "testdata",
			cfg: config{
				exts: []string{".zh", ".go"},
				size: 0,
				list: true,
			},
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer

			if err := run(tc.root, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			res := buffer.String()

			if res != tc.expected {
				t.Errorf("Expected '%q', got '%q' instead\n", tc.expected, res)
			}
		})
	}
}

func TestRunDelExt(t *testing.T) {
	testCases := []struct {
		name      string
		cfg       config
		nDelete   int
		nNoDelete int
		expected  string
	}{
		{
			name:      "DeleteExtensionNoMatch",
			cfg:       config{exts: []string{".c"}, del: true},
			nDelete:   0,
			nNoDelete: 20,
			expected:  "",
		},
		{
			name: "DeleteExtensionMatch",
			cfg: config{
				exts: []string{".log"},
				del:  true,
			},
			nDelete:   5,
			nNoDelete: 15,
			expected:  "",
		},
		{
			name: "DeleteExtensionMixed",
			cfg: config{
				exts: []string{".log", ".html"},
				del:  true,
			},
			nDelete:   10,
			nNoDelete: 10,
			expected:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer, logBuffer bytes.Buffer

			tc.cfg.wLog = &logBuffer

			tempDir, cleanup := createTempDir(t, "")

			defer cleanup()

			if err := run(tempDir, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			res := buffer.String()

			if tc.expected != res {
				t.Errorf("Expected %q, got %q instead\n", tc.expected, res)
			}

			filesLeft, err := os.ReadDir(tempDir)

			if err != nil {
				t.Error(err)
			}

			if len(filesLeft) != tc.nNoDelete {
				t.Errorf("Expected %d files left, got %d instead\n",
					tc.nNoDelete, len(filesLeft))
			}

			expectedLogLines := tc.nDelete + 1
			lines := bytes.Split(logBuffer.Bytes(), []byte("\n"))

			if len(lines) != expectedLogLines {
				t.Errorf("Expected %d log lines, got %d instead\n",
					expectedLogLines, len(lines))
			}

		})
	}
}

func createTempDir(t *testing.T, dir string) (dirname string, cleanup func()) {
	t.Helper()
	tempDir := os.TempDir()
	walkDir := filepath.Join(tempDir, "walktest")

	// walkDir := "walktest"

	if dir != "" {
		walkDir = filepath.Join(tempDir, dir)
		// walkDir = dir
	}
	//os.Remove(walkDir)
	_, err := os.Stat(walkDir)
	if os.IsNotExist(err) {
		os.Mkdir(walkDir, 0775)
	}
	if dir == "" {
		createFiles(t, walkDir)
	}
	return walkDir, func() {
		os.RemoveAll(walkDir)
	}

}

func createFiles(t *testing.T, walkDir string) {
	t.Helper()
	exts := []string{".log", ".gz", ".html", ".ts"}

	for _, ext := range exts {
		for i := 1; i <= 5; i++ {
			fname := fmt.Sprintf("file%d%s", i, ext)
			fpath := filepath.Join(walkDir, fname)

			if err := os.WriteFile(fpath, []byte("dummy"), 0664); err != nil {
				t.Log(err)
			}
		}
	}

}

func TestRunArchive(t *testing.T) {
	// Archiving test cases
	testCases := []struct {
		name       string
		cfg        config
		nArchive   int
		nNoArchive int
	}{
		{
			name: "ArchiveExtensionMatch",
			cfg: config{
				exts: []string{".log"},
			},
			nArchive: 5, nNoArchive: 20,
		},
		{
			name: "ArchiveExtensionNoMatch",
			cfg: config{
				exts: []string{".c"},
			},
			nArchive: 0, nNoArchive: 25,
		},
		{
			name: "ArchiveExtensionMixed",
			cfg: config{
				exts: []string{".log", ".html"},
			},
			nArchive: 10, nNoArchive: 15,
		},
	}

	// Execute RunArchive test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Buffer for RunArchive output
			var buffer bytes.Buffer
			// Create temp dirs for RunArchive test
			tempDir, cleanup := createTempDir(t, "")
			defer cleanup()
			archiveDir, cleanupArchive := createTempDir(t, "archivetest")
			defer cleanupArchive()

			tc.cfg.archive = archiveDir
			if err := run(tempDir, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			// pattern := filepath.Join(tempDir, "/*")
			// expFiles, err := filepath.Glob(pattern)
			// if err != nil {
			// 	t.Fatal(err)
			// }

			// expOut := strings.Join(expFiles, "\n")

			// res := strings.TrimSpace(buffer.String())

			// if expOut != res {
			// 	t.Errorf("Expected %q, got %q instead\n", expOut, res)
			// }

			filesArchived, err := os.ReadDir(archiveDir)
			if err != nil {
				t.Fatal(err)
			}

			if len(filesArchived) != tc.nArchive {
				t.Errorf("Expected %d files archived, got %d instead\n",
					tc.nArchive, len(filesArchived))
			}
		})
	}
}
