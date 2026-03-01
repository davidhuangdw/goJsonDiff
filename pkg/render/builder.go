package render

import (
	"strings"

	"github.com/davidhuangdw/goJsonDiff/pkg/render/style"
)

func NewConsoleView() DeltaView {
	return &StyleView{
		Builder: &strings.Builder{},
		Format:  style.ConsoleFormat,
		Style:   &style.ConsoleStyle{},
	}
}

func NewHtmlView() DeltaView {
	return &StyleView{
		Builder: &strings.Builder{},
		Format:  style.HtmlFormat,
		Style:   &style.HtmlStyle{},
	}
}

func NewJsonPatchView() DeltaView {
	return &JsonPatchView{}
}
