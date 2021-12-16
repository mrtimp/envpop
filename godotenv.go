// thanks to https://github.com/joho/godotenv
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func readFile(filename string) (file *os.File, err error) {
	file, err = os.Open(filename)

	if err != nil {
		return
	}

	// @todo handle errors

	return
}

// Parse reads an env file from io.Reader, returning a map of keys and values.
func Parse(r io.Reader) (envFile string, err error) {
	var lines []string

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return
	}

	for i, line := range lines {
		lines[i] = parseEnvLine(line)
	}

	envFile = strings.Join(lines, "\n")

	return
}

func parseEnvLine(line string) string {
	var key string
	var val string
	key, val, _ = parseEnvLineToKeyValue(line)

	if len(key) > 0 && !strings.HasPrefix(key, "#") {
		return fmt.Sprintf("%s=%s", key, val)
	}

	return key
}

func parseEnvLineToKeyValue(line string) (key string, value string, err error) {
	line = strings.TrimSpace(line)

	if len(line) == 0 {
		// blank lines are fine
		return
	}

	if strings.HasPrefix(line, "#") {
		key = line
		return
	}

	var comments string

	// ditch the comments (but keep quoted hashes)
	if strings.Contains(line, "#") {
		// @todo handle the comments!
		segmentsBetweenHashes := strings.Split(line, "#")
		quotesAreOpen := false
		var segmentsToKeep []string
		for _, segment := range segmentsBetweenHashes {
			if strings.Count(segment, "\"") == 1 || strings.Count(segment, "'") == 1 {
				if quotesAreOpen {
					quotesAreOpen = false
					segmentsToKeep = append(segmentsToKeep, segment)
				} else {
					quotesAreOpen = true
				}
			}

			if len(segmentsToKeep) == 0 || quotesAreOpen {
				segmentsToKeep = append(segmentsToKeep, segment)
			}
		}

		// keep the comment so we can append it later
		comments = strings.TrimSpace(segmentsBetweenHashes[1])

		line = strings.Join(segmentsToKeep, "#")
		line = strings.TrimSpace(line)
	}

	firstEquals := strings.Index(line, "=")
	firstColon := strings.Index(line, ":")
	splitString := strings.SplitN(line, "=", 2)
	if firstColon != -1 && (firstColon < firstEquals || firstEquals == -1) {
		// this is a yaml-style line
		splitString = strings.SplitN(line, ":", 2)
	}

	if len(splitString) != 2 {
		// @todo handle?
		err = errors.New("can't separate key from value")
		return
	}

	// Parse the key
	key = splitString[0]
	if strings.HasPrefix(key, "export") {
		key = strings.TrimPrefix(key, "export")
	}
	key = strings.TrimSpace(key)

	var exportRegex = regexp.MustCompile(`^\s*(?:export\s+)?(.*?)\s*$`)

	key = exportRegex.ReplaceAllString(splitString[0], "$1")

	value = splitString[1]

	parsedValue := parseValue(value)

	if len(parsedValue) == 0 {
		// look up the key to find out if it is already set in the environment
		// if it is then override the value from the .env file
		if _, present := os.LookupEnv(key); present {
			parsedValue = os.Getenv(key)

			value = handleValueType(parsedValue)
		}
	}

	if len(comments) > 0 {
		value = fmt.Sprintf("%s # %s", value, comments)
	}

	return
}

func parseValue(value string) string {
	var (
		singleQuotesRegex  = regexp.MustCompile(`\A'(.*)'\z`)
		doubleQuotesRegex  = regexp.MustCompile(`\A"(.*)"\z`)
		escapeRegex        = regexp.MustCompile(`\\.`)
		unescapeCharsRegex = regexp.MustCompile(`\\([^$])`)
	)

	// trim
	value = strings.Trim(value, " ")

	// check if we've got quoted values or possible escapes
	if len(value) > 1 {
		singleQuotes := singleQuotesRegex.FindStringSubmatch(value)

		doubleQuotes := doubleQuotesRegex.FindStringSubmatch(value)

		if singleQuotes != nil || doubleQuotes != nil {
			// pull the quotes off the edges
			value = value[1 : len(value)-1]
		}

		if doubleQuotes != nil {
			// expand newlines
			value = escapeRegex.ReplaceAllStringFunc(value, func(match string) string {
				c := strings.TrimPrefix(match, `\`)
				switch c {
				case "n":
					return "\n"
				case "r":
					return "\r"
				default:
					return match
				}
			})
			// unescape characters
			value = unescapeCharsRegex.ReplaceAllString(value, "$1")
		}

		if singleQuotes == nil {
			value = expandVariables(value)
		}
	}

	return value
}

func handleValueType(value string) string {
	// boolean type
	if value != "0" && value != "1" {
		if val, err := strconv.ParseBool(value); err == nil {
			return fmt.Sprintf("%v", val)
		}
	}

	// integer type
	if intVal, err := strconv.Atoi(value); err == nil {
		// treat numbers prefixed with 0 as a string
		if !strings.HasPrefix(value, "0") || intVal == 0 {
			return fmt.Sprintf("%v", intVal)
		}
	}

	// float type
	if strings.Contains(value, ".") {
		if val, err := strconv.ParseFloat(value, 10); err == nil {
			return fmt.Sprintf("%v", val)
		}
	}

	// null type
	if strings.ToLower(value) == "null" {
		return "null"
	}

	if value != "" {
		// default everything else to string
		return fmt.Sprintf("\"%s\"", value)
	}

	return ""
}

func expandVariables(v string) string {
	var expandVarRegex = regexp.MustCompile(`(\\)?(\$)(\()?\{?([A-Z0-9_]+)?\}?`)

	return expandVarRegex.ReplaceAllStringFunc(v, func(s string) string {
		submatch := expandVarRegex.FindStringSubmatch(s)

		if submatch == nil {
			return s
		}

		// if submatch[1] == "\\" || submatch[2] == "(" {
		//     return submatch[0][1:]
		// } else if submatch[4] != "" {
		//     return m[submatch[4]]
		// }

		return s
	})
}
