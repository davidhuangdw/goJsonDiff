package style

import (
	. "github.com/davidhuangdw/goJsonDiff/pkg/types"
)

var (
	HtmlFormat = Format{
		Indent:    "&emsp;",
		BreakLine: "<br>\n",
	}
)

type HtmlStyle struct{}

func (st *HtmlStyle) Styled(op Operation, content string) string {
	switch op {
	case EQUAL:
		return `<span style="color: gray">` + content + `</span>`
	case ADD:
		return `<span style="background-color: #bbffbb">` + content + `</span>`
	case REMOVE:
		return `<span style="background-color: #ffbbbb"><strike>` + content + `</strike></span>`
	case REPLACE:
		return `<span style="background-color: #ffffbb">` + content + `</span>`
	default:
		return content
	}
}
