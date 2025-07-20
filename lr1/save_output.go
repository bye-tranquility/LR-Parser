package lr1

import (
	"bytes"
	"context"
	"fmt"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

func (p *Parser) renderAutomaton(ctx context.Context) ([]byte, error) {
	g, err := graphviz.New(ctx)
	if err != nil {
		return nil, err
	}
	graph, err := g.Graph()
	if err != nil {
		return nil, err
	}
	defer func() {
		graph.Close()
		g.Close()
	}()

	states, transitions := p.BuildAutomaton()

	var nodes []*cgraph.Node
	for i, state := range states {
		node, err := graph.CreateNodeByName(fmt.Sprintf("State %d\n%s", i, p.stateToString(state)))
		node.SetShape(graphviz.RectShape)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}

	for from, transition := range transitions {
		for symbol, to := range transition {
			edge, err := graph.CreateEdgeByName("", nodes[from], nodes[to])
			if err != nil {
				return nil, err
			}
			edge.SetLabel(symbol)
		}
	}

	var buf bytes.Buffer
	if err := g.Render(ctx, graph, "dot", &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *Parser) SaveAutomatonToPng(filename string) error {
	ctx := context.Background()

	automatonBytes, err := p.renderAutomaton(ctx)
	if err != nil {
		return err
	}
	if err := p.SaveToPng(ctx, automatonBytes, filename); err != nil {
		return err
	}
	return nil
}
