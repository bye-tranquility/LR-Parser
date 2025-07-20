package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"parser/grammar"
	"parser/lr0"
	"parser/lr1"
)

const Accept = lr0.Accept

type Parser interface {
	BuildTable() bool
	Algo(word string) int
	SaveAutomatonToPng(filename string) error
	SaveTableToPng(filename string) error
}

func run(reader io.Reader, writer io.Writer, parserType string, automatonFileName string, tableFileName string) int {
	gr := grammar.NewGrammar()

	scanner := bufio.NewScanner(reader)
	gr.ReadInput(scanner)

	var parser Parser
	switch parserType {
	case "lr0":
		parser = &lr0.Parser{
			Grammar: gr,
			Output:  writer,
		}
	case "lr1":
		parser = &lr1.Parser{
			Parser: lr0.Parser{
				Grammar: gr,
				Output:  writer,
			},
		}
	}

	if !parser.BuildTable() {
		return 1
	}

	if automatonFileName != "" {
		err := parser.SaveAutomatonToPng(automatonFileName)
		if err != nil {
			_, _ = fmt.Fprintln(writer, err)
		}
	}

	if tableFileName != "" {
		err := parser.SaveTableToPng(tableFileName)
		if err != nil {
			_, _ = fmt.Fprintln(writer, err)
		}
	}

	for scanner.Scan() {
		word := scanner.Text()

		if word == "exit" {
			break
		}
		_, _ = fmt.Fprintln(writer, "Word:", word)

		_, _ = fmt.Fprint(writer, "Belongs: ")

		if parser.Algo(word) == Accept {
			_, _ = fmt.Fprintln(writer, "YES")
		} else {
			_, _ = fmt.Fprintln(writer, "NO")
		}
	}
	return 0
}

func main() {
	parserType := flag.String("parser", "lr1", "Parser type: lr1 or lr0")
	var automatonFileName, tableFileName string
	flag.StringVar(&automatonFileName, "export-automaton", "", "Name a file in the format: file_to_be_saved_to.png")
	flag.StringVar(&tableFileName, "export-table", "", "Name a file in the format: file_to_be_saved_to.png¨")
	flag.Parse()
	os.Exit(run(os.Stdin, os.Stdout, *parserType, automatonFileName, tableFileName))
}
