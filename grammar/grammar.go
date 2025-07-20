package grammar

import (
	"bufio"
	"fmt"
	"strings"
)

type Rule struct {
	Left  string
	Right []string
}

type Grammar struct {
	Rules        []Rule
	Terminals    []string
	NonTerminals []string
	StartSymbol  string
}

func NewGrammar() *Grammar {
	return &Grammar{
		Rules:        []Rule{},
		Terminals:    []string{"$"},
		NonTerminals: []string{},
		StartSymbol:  "",
	}
}

func (gr *Grammar) AddRule(rawRule string) {
	rawRule = strings.TrimSpace(rawRule)
	splitted := strings.Split(rawRule, "->")

	left := strings.TrimSpace(splitted[0])
	right := strings.TrimSpace(splitted[1])
	right = strings.ReplaceAll(right, " ", "")
	gr.Rules = append(gr.Rules, Rule{left, strings.Split(right, "")})
}

func (gr *Grammar) ReadInput(scanner *bufio.Scanner) {
	scanner.Scan()
	line := scanner.Text()
	for _, letter := range line {
		gr.NonTerminals = append(gr.NonTerminals, string(letter))
	}

	scanner.Scan()
	line = scanner.Text()
	for _, letter := range line {
		gr.Terminals = append(gr.Terminals, string(letter))
	}

	scanner.Scan()
	line = scanner.Text()
	var p int
	_, _ = fmt.Sscanf(line, "%d", &p)

	for range p {
		scanner.Scan()
		line = scanner.Text()
		line = strings.ReplaceAll(line, " ", "")
		gr.AddRule(line)
	}

	scanner.Scan()
	line = scanner.Text()
	gr.StartSymbol = line
}

func (gr *Grammar) IsTerminal(letter string) bool {
	for _, term := range gr.Terminals {
		if term == letter {
			return true
		}
	}
	return false
}

func (gr *Grammar) PrintGrammar() {
	for _, rule := range gr.Rules {
		fmt.Println(rule.Left, rule.Right)
	}
	fmt.Println("Start Symbol:", gr.StartSymbol)
	fmt.Println("Nonterminals:", gr.NonTerminals)
	fmt.Println("Terminals:", gr.Terminals)
}
