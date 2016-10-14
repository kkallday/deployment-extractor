package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const manifestStartToken = "DirectorJobRunner: Manifest:"

func main() {
	manifest, err := extractManifest()
	if err != nil {
		panic(err)
	}

	fmt.Print(manifest)
}

func extractManifest() (string, error) {
	debugBuf := new(bytes.Buffer)

	_, err := io.Copy(debugBuf, os.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to copy manifest from stdin: %v", err)
	}

	manifestBuf := new(bytes.Buffer)

	err = seekToManifest(debugBuf)
	if err != nil {
		return "", fmt.Errorf("failed to find beginning of manifest in debug log: %v", err)
	}

	err = collectManifest(manifestBuf, *debugBuf)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve manifest from debug log: %v", err)
	}

	manifest, err := ioutil.ReadAll(manifestBuf)
	if err != nil {
		return "", fmt.Errorf("unexpected error: %v", err)
	}

	return string(manifest), nil
}

func seekToManifest(buf *bytes.Buffer) error {
	for {
		rawLine, err := buf.ReadBytes(byte('\n'))
		switch err {
		case nil:
			if strings.Contains(string(rawLine), manifestStartToken) {
				return nil
			}
		case io.EOF:
			return fmt.Errorf(`could not find %s`, manifestStartToken)
		default:
			return fmt.Errorf("unexpected error: %v", err)
		}
	}
}

func collectManifest(dst *bytes.Buffer, debugBuf bytes.Buffer) error {
	for {
		rawLine, err := debugBuf.ReadBytes(byte('\n'))
		switch err {
		case nil:
			line := string(rawLine)

			if strings.Contains(line, "D, ") || strings.Contains(line, "I, ") {
				return nil
			}

			_, err := dst.Write(rawLine)
			if err != nil {
				return fmt.Errorf("unexpected error: %v", err)
			}
		case io.EOF:
			return fmt.Errorf("reached end of log before finding end of manifest")
		default:
			return fmt.Errorf("unexpected error: %v", err)
		}
	}
}
