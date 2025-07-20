package lr0

import (
	"io"
	"parser/grammar"
	"strings"
)

const (
	Accept = iota
	Reject
	Shift
	Reduce
	Goto
)

type TableEntry struct {
	Action int
	Value  int
}

type Item struct {
	Rule grammar.Rule
	Dot  int
}

type Parser struct {
	Grammar *grammar.Grammar
	Table   []map[string]TableEntry
	Output  io.Writer
}

func (p *Parser) Algo(word string) int {
	statesStack := []int{0}
	symbolsStack := []string{}
	tokens := strings.Split(word, "")
	tokens = append(tokens, "$")

	currTokenIdx := 0
	for currTokenIdx < len(tokens) {
		currToken := tokens[currTokenIdx]
		currState := statesStack[len(statesStack)-1]

		entry, exists := p.Table[currState][currToken]
		if !exists {
			return Reject
		}

		switch entry.Action {
		case Shift:
			symbolsStack = append(symbolsStack, currToken)
			statesStack = append(statesStack, entry.Value)
			currTokenIdx++

		case Reduce:
			rule := p.Grammar.Rules[entry.Value]
			popCount := len(rule.Right)

			if len(statesStack) < popCount || len(symbolsStack) < popCount {
				return Reject
			}

			statesStack = statesStack[:len(statesStack)-popCount]
			symbolsStack = symbolsStack[:len(symbolsStack)-popCount]

			newSymbol := rule.Left
			newState := statesStack[len(statesStack)-1]

			gotoEntry, exists := p.Table[newState][newSymbol]
			if !exists || gotoEntry.Action != Goto {
				return Reject
			}

			symbolsStack = append(symbolsStack, newSymbol)
			statesStack = append(statesStack, gotoEntry.Value)

		case Accept:
			if currToken == "$" && currTokenIdx == len(tokens)-1 {
				return Accept
			}
			return Reject

		default:
			return Reject
		}
	}
	return Reject
}
