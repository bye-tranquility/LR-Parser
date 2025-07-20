package lr1

import (
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

func (p *Parser) getFirst(sequence []string) map[string]bool {
	set := make(map[string]bool)

	for _, symbol := range sequence {
		if p.Grammar.IsTerminal(symbol) {
			set[symbol] = true
			return set
		} else {
			for terminal := range p.first[symbol] {
				set[terminal] = true
			}
			if !p.nullable[symbol] {
				return set
			}
		}
	}
	return set
}

func (p *Parser) computeNullable() {
	p.nullable = make(map[string]bool)

	for _, rule := range p.Grammar.Rules {
		if len(rule.Right) == 0 {
			p.nullable[rule.Left] = true
		}
	}
	changed := true
	for changed {
		changed = false
		for _, rule := range p.Grammar.Rules {
			if p.nullable[rule.Left] {
				continue
			}

			allNullable := true
			for _, symbol := range rule.Right {
				if !p.nullable[symbol] {
					allNullable = false
					break
				}
			}

			if allNullable && !p.nullable[rule.Left] {
				p.nullable[rule.Left] = true
				changed = true
			}
		}
	}
}

func (p *Parser) computeFirst() {
	p.first = make(map[string]map[string]bool)

	for _, nonterminal := range p.Grammar.NonTerminals {
		p.first[nonterminal] = make(map[string]bool)
	}
	for _, terminal := range p.Grammar.Terminals {
		p.first[terminal] = map[string]bool{terminal: true}
	}

	changed := true
	for changed {
		changed = false
		for _, rule := range p.Grammar.Rules {
			for _, symbol := range rule.Right {
				for terminal := range p.first[symbol] {
					if !p.first[rule.Left][terminal] {
						p.first[rule.Left][terminal] = true
						changed = true
					}
				}

				if !p.nullable[symbol] {
					break
				}
			}
		}
	}
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
	b.WriteString(", ")
	b.WriteString(item.lookahead)
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
