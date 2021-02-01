package main

import (
	"os"
	"testing"
)

func TestEnvLenghtHasChanged(t *testing.T) {
	envLen := len(os.Environ())

	input := []string{
		"PORT=8080",
	}
	ParseAnSetEnv(input)

	newEnvLen := len(os.Environ())

	if newEnvLen <= envLen {
		t.Errorf("No new environment variable added")
	}
}

func TestEnvHasCorrectLength(t *testing.T) {
	envLen := len(os.Environ())

	input := []string{
		"PORT=8080",
		"ENV=development",
		"DATA= some data",
		"",
		"# Comment",
	}
	ParseAnSetEnv(input)

	newEnvLen := len(os.Environ())

	if newEnvLen == envLen+3 {
		t.Errorf("No new environment variable added")
	}
}

func TestEnvParser(t *testing.T) {
	input := []string{
		"PORT=8080",
		"ENV=development",
		"DATA= some data",
		"",
		"# Comment",
		"#ORDER=66",
	}
	ParseAnSetEnv(input)

	if os.Getenv("PORT") != "8080" {
		t.Errorf("Env variable PORT no set")
	}

	if os.Getenv("ENV") != "development" {
		t.Errorf("Env variable ENV no set")
	}

	if os.Getenv("DATA") != "some data" {
		t.Errorf("Env variable DATA no set")
	}

	if os.Getenv("ORDER") == "66" {
		t.Errorf("Comment #ORDER=66 was not ignored")
	}
}
