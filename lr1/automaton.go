package lr1

import (
	"parser/grammar"
	"parser/lr0"
)

func (p *Parser) BuildAutomaton() (states [][]Item, transitions []map[string]int) {
	p.computeNullable()
	p.computeFirst()

	seen := make(map[string]bool)

	initialItem := Item{lr0.Item{grammar.Rule{"S'", []string{p.Grammar.StartSymbol}}, 0}, "$"}
	initialState := p.closure([]Item{initialItem})

	states = append(states, initialState)
	transitions = append(transitions, make(map[string]int))

	seen[p.stateToString(initialState)] = true

	for i := 0; i < len(states); i++ {
		currState := states[i]
		nextSymbols := getNextSymbols(currState)

		for _, symbol := range nextSymbols {
			nextState := p.gotoState(currState, symbol)
			nextStateString := p.stateToString(nextState)
			if !seen[nextStateString] {
				states = append(states, nextState)
				seen[nextStateString] = true
				transitions = append(transitions, make(map[string]int))
			}
			nextStateIdx := p.getStateIndex(states, nextState)
			transitions[i][symbol] = nextStateIdx
		}
	}
	return
}

func (p *Parser) closure(state []Item) []Item {
	result := append([]Item{}, state...)
	closureSet := make(map[string]bool)
	for _, item := range state {
		closureSet[p.itemToString(item)] = true
	}

	added := true
	for added {
		added = false
		currentLen := len(result)
		for i := 0; i < currentLen; i++ {
			currItem := result[i]
			if currItem.Dot >= len(currItem.Rule.Right) {
				continue
			}

			nextSymbol := currItem.Rule.Right[currItem.Dot]
			if !p.Grammar.IsTerminal(nextSymbol) {
				var rest []string
				if currItem.Dot+1 < len(currItem.Rule.Right) {
					rest = currItem.Rule.Right[currItem.Dot+1:]
				}
				sequence := append(rest, currItem.lookahead)
				firstSet := p.getFirst(sequence)

				for _, rule := range p.Grammar.Rules {
					if rule.Left == nextSymbol {
						for terminal := range firstSet {
							newItem := Item{lr0.Item{rule, 0}, terminal}
							itemString := p.itemToString(newItem)
							if !closureSet[itemString] {
								closureSet[itemString] = true
								result = append(result, newItem)
								added = true
							}
						}
					}
				}
			}
		}
	}
	return result
}

func (p *Parser) gotoState(state []Item, symbol string) []Item {
	var result []Item
	for _, item := range state {
		if item.Dot < len(item.Rule.Right) && item.Rule.Right[item.Dot] == symbol {
			newItem := Item{lr0.Item{item.Rule, item.Dot + 1}, item.lookahead}
			result = append(result, newItem)
		}
	}
	return p.closure(result)
}
