package main

import (
	"context"
	"fmt"
	"i2l"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	il, err := i2l.DefaultGoogleAIRAG(ctx)
	if err != nil {
		panic(err)
	}

	graph, err := il.GenerateGraphFromCode(ctx, `package main

	import "fmt"

	func main() {
		fmt.Println("Hello, World!")
	}`)
	if err != nil {
		panic(err)
	}

	fmt.Println("Extracted graph tuples:")
	for _, t := range graph {
		fmt.Println(t.String())
	}

	codeResult, err := il.GenerateCodeFromGraph(ctx, "Java", graph)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n\nGenerated code:")
	fmt.Println(codeResult.Code)
}
