package lr1

import (
	"parser/lr0"
)

type Item struct {
	lr0.Item
	lookahead string
}

type Parser struct {
	lr0.Parser
	nullable map[string]bool
	first    map[string]map[string]bool
}
