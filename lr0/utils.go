package lr0

import (
	"fmt"
	"parser/grammar"
	"sort"
	"strings"
)

func (p *Parser) getStateIndex(states [][]Item, target []Item) int {
	targetStr := p.stateToString(target)
	for i, state := range states {
		if p.stateToString(state) == targetStr {
			return i
		}
	}
	return -1
}

func (p *Parser) GetRuleIndex(target grammar.Rule) int {
	for i, Rule := range p.Grammar.Rules {
		if Rule.Left == target.Left && len(Rule.Right) == len(target.Right) {
			equal := true
			for j := range Rule.Right {
				if Rule.Right[j] != target.Right[j] {
					equal = false
					break
				}
			}
			if equal {
				return i
			}
		}
	}
	return -1
}

func getNextSymbols(state []Item) []string {
	seen := make(map[string]bool)
	for _, item := range state {
		if item.Dot < len(item.Rule.Right) {
			symbol := item.Rule.Right[item.Dot]
			seen[symbol] = true
		}
	}
	var symbols []string
	for s := range seen {
		symbols = append(symbols, s)
	}
	sort.Strings(symbols)
	return symbols
}

func (p *Parser) itemToString(item Item) string {
	var b strings.Builder
	b.WriteString(item.Rule.Left)
	b.WriteString("->")
	for i, symbol := range item.Rule.Right {
		if i == item.Dot {
			b.WriteString(".")
		}
		b.WriteString(symbol)
	}
	if item.Dot == len(item.Rule.Right) {
		b.WriteString(".")
	}
	return b.String()
}

func (p *Parser) stateToString(state []Item) string {
	var result []string
	for _, item := range state {
		result = append(result, p.itemToString(item))
	}
	sort.Strings(result)
	return strings.Join(result, "\n")
}

func DescribeAction(entry TableEntry) string {
	if entry.Action == Reduce {
		return fmt.Sprintf("reduce by Rule %d", entry.Value)
	} else if entry.Action == Shift {
		return fmt.Sprintf("shift to state %d", entry.Value)
	}
	return fmt.Sprint("<unexpected action>")
}
