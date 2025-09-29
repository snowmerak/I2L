package main

import (
	"context"
	"flag"
	"fmt"
	"i2l"
	"os"
	"os/signal"
)

func main() {
	inputFile := flag.String("f", "", "Path to the source code file to analyze.")
	targetLang := flag.String("l", "", "Target language for code generation.")
	outputFile := flag.String("o", "", "Name of the output file to save the generated code.")
	flag.Parse()

	if *inputFile == "" || *targetLang == "" || *outputFile == "" {
		fmt.Println("Usage: go run . -f <path/to/source.go> -l <language> -o <path/to/output.txt>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	code, err := os.ReadFile(*inputFile)
	if err != nil {
		panic(fmt.Errorf("failed to read input file %s: %w", *inputFile, err))
	}

	il, err := i2l.DefaultGoogleAIRAG(ctx)
	if err != nil {
		panic(err)
	}

	graph, err := il.GenerateGraphFromCode(ctx, string(code))
	if err != nil {
		panic(err)
	}

	fmt.Println("Extracted graph tuples:")
	for _, t := range graph {
		fmt.Println(t.String())
	}

	codeResult, err := il.GenerateCodeFromGraph(ctx, *targetLang, graph)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\nGenerated code (saving to %s):\n", *outputFile)
	fmt.Println(codeResult.Code)

	if err := os.WriteFile(*outputFile, []byte(codeResult.Code), 0644); err != nil {
		panic(fmt.Errorf("failed to write output to file %s: %w", *outputFile, err))
	}
}
