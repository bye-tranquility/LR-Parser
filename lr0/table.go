package lr0

import (
	"fmt"
)

const (
	BuildSuccess = true
	BuildFailure = false
)

func (p *Parser) BuildTable() bool {
	states, transitions := p.BuildAutomaton()
	numStates := len(states)

	p.Table = make([]map[string]TableEntry, numStates)
	for i := range p.Table {
		p.Table[i] = make(map[string]TableEntry)
	}

	for stateIdx, state := range states {
		for _, item := range state {
			if item.Dot < len(item.Rule.Right) {
				nextSymbol := item.Rule.Right[item.Dot]
				if nextState, exists := transitions[stateIdx][nextSymbol]; exists {
					actionType := Shift
					if !p.Grammar.IsTerminal(nextSymbol) {
						actionType = Goto
					}
					if !p.setAction(stateIdx, nextSymbol, TableEntry{actionType, nextState}) {
						return BuildFailure
					}
				}
			} else {
				if item.Rule.Left == "S'" {
					if !p.setAction(stateIdx, "$", TableEntry{Accept, 0}) {
						return BuildFailure
					}
				} else {
					ruleIdx := p.GetRuleIndex(item.Rule)
					for _, terminal := range p.Grammar.Terminals {
						if !p.setAction(stateIdx, terminal, TableEntry{Reduce, ruleIdx}) {
							return BuildFailure
						}
					}
				}
			}
		}
	}
	return BuildSuccess
}

func (p *Parser) setAction(stateIdx int, symbol string, new TableEntry) bool {
	if prev, exists := p.Table[stateIdx][symbol]; exists {
		if prev.Action != new.Action || prev.Value != new.Value {
			_, _ = fmt.Fprintf(p.Output, "CONFLICT in state %d on `%s`:\n", stateIdx, symbol)
			_, _ = fmt.Fprintf(p.Output, "No way to decide whether to \n- %s\nor\n- %s\n",
				DescribeAction(prev), DescribeAction(new))
			_, _ = fmt.Fprintln(p.Output, "The grammar is not LR(0)")
			return BuildFailure
		}
	}
	p.Table[stateIdx][symbol] = new
	return BuildSuccess
}
