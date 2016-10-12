package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	manifestBuf := bytes.NewBuffer([]byte{})
	startRead := false

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		if string(line) == "---" {
			startRead = true
			continue
		}

		if startRead {
			if strings.Contains(string(line), "D, ") || strings.Contains(string(line), "I, ") {
				break
			}

			_, err := manifestBuf.Write(line)
			if err != nil {
				panic(err)
			}

			err = manifestBuf.WriteByte(byte('\n'))
			if err != nil {
				panic(err)
			}
		}
	}

	rawManifest, err := ioutil.ReadAll(manifestBuf)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(rawManifest))
}
