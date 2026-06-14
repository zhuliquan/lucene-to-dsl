package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	lucene_to_dsl "github.com/zhuliquan/lucene-to-dsl"
)

func main() {
	var mappingPath string
	var luceneQuery string

	flag.StringVar(&mappingPath, "m", "", "mapping file path")
	flag.StringVar(&mappingPath, "mapping", "", "mapping file path")
	flag.StringVar(&luceneQuery, "q", "", "lucene query")
	flag.StringVar(&luceneQuery, "query", "", "lucene query")
	flag.Parse()

	if luceneQuery == "" {
		fmt.Fprintln(os.Stderr, "Error: lucene query is required")
		os.Exit(1)
	}

	var opts []lucene_to_dsl.Option

	if mappingPath != "" {
		mappingData, err := os.ReadFile(mappingPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading mapping file: %v\n", err)
			os.Exit(1)
		}
		opts = append(opts, lucene_to_dsl.WithMappingData(mappingData))
	}

	dsl, err := lucene_to_dsl.LuceneToDSL(luceneQuery, opts...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	jsonBytes, err := json.MarshalIndent(dsl, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling to JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonBytes))
}
