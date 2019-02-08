package names

import "fmt"

// Generator generates unique names
type Generator interface {
	Generate() string
}

type sequencialPrefixedName struct {
	prefix string
	format string
	next   int
}

func NewSequentialPrefixedName(prefix string) Generator {
	return &sequencialPrefixedName{
		prefix: prefix,
		format: "%s%4.4d",
		next:   1000,
	}
}

func (g *sequencialPrefixedName) Generate() string {
	g.next++
	name := fmt.Sprintf(g.format, g.prefix, g.next)
	return name
}
