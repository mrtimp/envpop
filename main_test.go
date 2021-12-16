package main

import (
	_ "embed"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseEnvLine(t *testing.T) {
	tests := []struct {
		line     string
		env      string
		expected string
	}{
		{line: "STRING=", expected: "STRING="},
		{line: "STRING=\"VALUE\"", expected: "STRING=\"VALUE\""},
		{line: "STRING=", env: "STRING=VALUE", expected: "STRING=\"VALUE\""},

		{line: "BOOLEAN=true", expected: "BOOLEAN=true"},
		{line: "BOOLEAN=", env: "BOOLEAN=true", expected: "BOOLEAN=true"},
		{line: "BOOLEAN=", env: "BOOLEAN=false", expected: "BOOLEAN=false"},

		{line: "NUMBER=0", expected: "NUMBER=0"},
		{line: "NUMBER=", env: "NUMBER=999999", expected: "NUMBER=999999"},

		{line: "NUMBER=0.1", expected: "NUMBER=0.1"},
		{line: "NUMBER=", env: "NUMBER=0.1", expected: "NUMBER=0.1"},

		{line: "NULL=null", expected: "NULL=null"},
		{line: "NULL=null", env: "NULL=null", expected: "NULL=null"},
	}

	for _, tc := range tests {
		if tc.env != "" {
			env := strings.Split(tc.env, "=")
			err := os.Setenv(env[0], env[1])

			if err != nil {
				//
			}
		}

		actual := parseEnvLine(tc.line)

		assert.Equal(t, actual, tc.expected)
	}
}

//go:embed testing/fixtures/.env.dist
var dotenvdist string

//go:embed testing/fixtures/.env.dist.expected
var dotenvdistexpected string

func Test_parseEnvFile(t *testing.T) {
	setEnvVar("DYNAMIC_STRING", "VALUE")
	setEnvVar("DYNAMIC_COMMENT_STRING", "VALUE")
	setEnvVar("DYNAMIC_NUMBER", "123")
	setEnvVar("DYNAMIC_BOOLEAN", "true")
	setEnvVar("DYNAMIC_FLOAT", "1.23")
	setEnvVar("DYNAMIC_NULL", "null")

	envFile := strings.NewReader(dotenvdist)

	actual, _ := Parse(envFile)

	// account for potential whitespace EOL
	actual = strings.TrimSpace(actual)
	dotenvdistexpected = strings.TrimSpace(dotenvdistexpected)

	assert.Equal(t, actual, dotenvdistexpected)
}

func setEnvVar(key string, val string) {
	err := os.Setenv(key, val)

	if err != nil {
		// ignore
	}
}
