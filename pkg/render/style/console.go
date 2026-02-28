package style

import (
	. "goJsonDiff/pkg/types"

	"github.com/logrusorgru/aurora/v4"
)

var (
	ConsoleFormat = Format{
		Indent:    "  ",
		BreakLine: "\n",
	}
)

type ConsoleStyle struct{}

func (st *ConsoleStyle) Styled(op Operation, content string) string {
	switch op {
	case ADD:
		return aurora.Green(content).String()
	case REMOVE:
		return aurora.Red(content).StrikeThrough().String()
	case REPLACE:
		return aurora.Yellow(content).String()
	case EQUAL:
		return aurora.BrightBlack(content).String()
	default:
		return content
	}
}
