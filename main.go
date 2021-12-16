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

	fmt.Println(Parse(envFile))
}
