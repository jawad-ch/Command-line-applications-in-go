package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"
)

const (
	header = `<!DOCTYPE html>
<html>
	<head> 
		<meta http-equiv="content-type" content="text/html; charset=utf-8">
		<title>Markdown Preview Tool</title> 
	</head> 
<body>
`
	footer = ` </body>
</html> `
)

func main() {
	// Parse flags
	filename := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	flag.Parse()

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename, os.Stdout, *skipPreview); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename string, out io.Writer, skipPreview bool) error {
	input, err := os.ReadFile(filename)

	if err != nil {
		return err
	}

	htmlData := parseContent(input)

	// Create temporary file and check for errors
	temp, err := os.CreateTemp("", "mdp*.html")

	if err != nil {
		return err
	}

	if err := temp.Close(); err != nil {
		return err
	}

	outName := temp.Name()
	//outName := fmt.Sprintf("%s.html", filepath.Base(filename))

	fmt.Println(out, outName)

	if err := saveHtml(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			_ = fmt.Errorf("errrr: %q", err)
		}
	}(outName)

	return preview(outName)
}

func parseContent(input []byte) []byte {
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	var buffer bytes.Buffer

	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)

	return buffer.Bytes()
}

func saveHtml(outFileName string, data []byte) error {
	return os.WriteFile(outFileName, data, 0644)
}

func preview(fName string) error {
	cName := ""
	var cParams []string
	// Define executable based on OS
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}
	// Append filename to parameters slice
	cParams = append(cParams, fName)
	// Locate executable in PATH
	cPath, err := exec.LookPath(cName)

	if err != nil {
		return err
	}
	// Open the file using default program
	err = exec.Command(cPath, cParams...).Run()

	// Give the browser some time to open the file before deleting it
	time.Sleep(2 * time.Second)
	return err
}
