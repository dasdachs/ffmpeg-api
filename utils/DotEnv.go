// The utils.DotEnv package includes a simple .env parser
package utils

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type fileReader interface {
	ReadEnvFile() []string
}

func ParseEnv() {
	fr := &envFileReader{}
	ParseAndSetEnv(fr)
}

func ParseAndSetEnv(fr fileReader) {
	content := fr.ReadEnvFile()
	for _, line := range content {
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}

		vals := strings.Split(line, "=")

		key := vals[0]
		val := strings.TrimSpace(strings.Split(vals[1], "#")[0])

		err := os.Setenv(key, val)
		if err != nil {
			log.Fatalf("Could not parse %s\n", line)
		}
	}
}

type envFileReader struct{}

func (fr envFileReader) ReadEnvFile() []string {
	_, err := os.Stat(".env")
	if os.IsNotExist(err) {
		return make([]string, 0)
	}

	data, err := ioutil.ReadFile(".env")
	if err != nil {
		log.Panic("Could not parse .env file")
	}

	content := strings.Split(string(data), "\n")

	return content
}
