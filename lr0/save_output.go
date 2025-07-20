package lr0

import (
	"bytes"
	"context"
	"fmt"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"log"
	"os"
	"path/filepath"
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

func (p *Parser) renderTable(ctx context.Context) ([]byte, error) {
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

	var symbols []string
	symbols = append(symbols, p.Grammar.NonTerminals...)
	for _, terminal := range p.Grammar.Terminals {
		if terminal != "$" {
			symbols = append(symbols, terminal)
		}
	}
	symbols = append(symbols, "$")

	label := "{"
	for stateIdx := 0; stateIdx < len(p.Table); stateIdx++ {
		label += fmt.Sprintf("| %d ", stateIdx)
	}
	label += "}"

	for _, symbol := range symbols {
		row := fmt.Sprintf("{ %s ", symbol)
		for stateIdx := 0; stateIdx < len(p.Table); stateIdx++ {
			entry, exists := p.Table[stateIdx][symbol]
			if !exists {
				row += "|     "
				continue
			}
			switch entry.Action {
			case Shift:
				row += fmt.Sprintf("| s(%d) ", entry.Value)
			case Reduce:
				row += fmt.Sprintf("| r(%d) ", entry.Value+1) // rules are numbered starting from 1
			case Goto:
				row += fmt.Sprintf("| %d ", entry.Value)
			case Accept:
				row += "| acc "
			default:
				row += "|     "
			}
		}
		row += "}"
		label += "|" + row
	}

	node, err := graph.CreateNodeByName("table")
	if err != nil {
		return nil, err
	}
	node.SetShape("record")
	node.SetLabel(label)

	var buf bytes.Buffer
	if err := g.Render(ctx, graph, "dot", &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *Parser) SaveToPng(ctx context.Context, graphBytes []byte, filename string) error {
	graph, err := graphviz.ParseBytes(graphBytes)
	if err != nil {
		return err
	}
	g, err := graphviz.New(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()

	if err := os.Mkdir("output", 0755); os.IsNotExist(err) && err != nil {
		return err
	}

	if err := g.RenderFilename(ctx, graph, "png", filepath.Join("output", filename)); err != nil {
		return err
	}
	return nil
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

func (p *Parser) SaveTableToPng(filename string) error {
	ctx := context.Background()

	tableBytes, err := p.renderTable(ctx)
	if err != nil {
		return err
	}
	if err := p.SaveToPng(ctx, tableBytes, filename); err != nil {
		return err
	}
	return nil
}
