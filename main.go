package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "debug-lexer" {
		input := "x = 5"
		lexer := NewLexer(input)

		for {
			token := lexer.NextToken()
			fmt.Printf("Token: %s, Value: '%s'\n", token.Type.String(), token.Value)
			if token.Type == EOF {
				break
			}
		}
		return
	}

	if len(os.Args) > 1 {
		// Run file
		filename := os.Args[1]
		if !strings.HasSuffix(filename, ".alo") {
			fmt.Println("Error: Alonso files must have .alo extension")
			os.Exit(1)
		}

		content, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}

		interpreter := NewInterpreter()
		err = interpreter.Execute(string(content))
		if err != nil {
			fmt.Printf("Runtime error: %v\n", err)
			os.Exit(1)
		}
	} else {
		// REPL mode
		fmt.Println("Welcome to Alonso - The F1 Programming Language!")
		fmt.Println("Type 'pit' to exit")

		interpreter := NewInterpreter()
		scanner := bufio.NewScanner(os.Stdin)

		for {
			fmt.Print("alonso> ")
			if !scanner.Scan() {
				break
			}

			line := strings.TrimSpace(scanner.Text())
			if line == "pit" {
				fmt.Println("Thanks for racing with Alonso!")
				break
			}

			if line == "" {
				continue
			}

			err := interpreter.Execute(line)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}
}
