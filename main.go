package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	file := flag.String("file", "", "The .env template file to read")

	flag.Parse()

	envFile, err := readFile(*file)

	if err != nil {
		log.Fatalf("Error loading the .env template file '%s'", *file)
	}

	var env, _ = Parse(envFile)

	if err != nil {
		log.Fatalf("Error parsing the .env template file '%s'", *file)
	}

	fmt.Printf("%s", env)
}
