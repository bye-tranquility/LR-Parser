package lr1

import (
	"fmt"
	"parser/lr0"
)

func (p *Parser) BuildTable() bool {
	states, transitions := p.BuildAutomaton()
	numStates := len(states)

	p.Table = make([]map[string]lr0.TableEntry, numStates)
	for i := range p.Table {
		p.Table[i] = make(map[string]lr0.TableEntry)
	}

	for stateIdx, state := range states {
		for _, item := range state {
			if item.Dot < len(item.Rule.Right) {
				nextSymbol := item.Rule.Right[item.Dot]
				if nextState, exists := transitions[stateIdx][nextSymbol]; exists {
					actionType := lr0.Shift
					if !p.Grammar.IsTerminal(nextSymbol) {
						actionType = lr0.Goto
					}
					if !p.setAction(stateIdx, nextSymbol, lr0.TableEntry{actionType, nextState}) {
						return lr0.BuildFailure
					}
				}
			} else {
				if item.Rule.Left == "S'" {
					if !p.setAction(stateIdx, "$", lr0.TableEntry{lr0.Accept, 0}) {
						return lr0.BuildFailure
					}
				} else {
					ruleIdx := p.GetRuleIndex(item.Rule)
					if !p.setAction(stateIdx, item.lookahead, lr0.TableEntry{lr0.Reduce, ruleIdx}) {
						return lr0.BuildFailure
					}
				}
			}
		}
	}
	return lr0.BuildSuccess
}

func (p *Parser) setAction(stateIdx int, symbol string, new lr0.TableEntry) bool {
	if prev, exists := p.Table[stateIdx][symbol]; exists {
		if prev.Action != new.Action || prev.Value != new.Value {
			_, _ = fmt.Fprintf(p.Output, "CONFLICT in state %d on `%s`:\n", stateIdx, symbol)
			_, _ = fmt.Fprintf(p.Output, "No way to decide whether to \n- %s\nor\n- %s\n",
				lr0.DescribeAction(prev), lr0.DescribeAction(new))
			_, _ = fmt.Fprintln(p.Output, "The grammar is not LR(1)")
			return lr0.BuildFailure
		}
	}
	p.Table[stateIdx][symbol] = new
	return lr0.BuildSuccess
}
