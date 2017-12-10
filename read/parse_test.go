package read

import "testing"

func TestParse(t *testing.T) {
	prog := "(hello world (a b))"
	parse("test", prog)
}
