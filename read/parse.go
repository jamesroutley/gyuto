package read

import "fmt"

type parser struct {
	items chan item
	ast   []interface{}
}

func parse(name, input string) {
	_, items := lex(name, input)
	p := &parser{
		items: items,
	}
	p.run()
	fmt.Println(p.ast)
}

func (p *parser) run() {
	switch i := p.next(); i.typ {
	case itemLeftParen:
		p.ast = p.parseSexpr()
	default:
		panic("parser run: should not be here")
	}
}

func (p *parser) next() item {
	return <-p.items
}

func (p *parser) parseSexpr() []interface{} {
	var expr []interface{}
	i := p.next()
	for {
		switch i.typ {
		case itemRightParen:
			return expr
		case itemLeftParen:
			expr = append(expr, p.parseSexpr())
		case itemAtom:
			expr = append(expr, i)
		default:
			panic("parser: shouldn't be here")
		}
		i = p.next()
	}

}
