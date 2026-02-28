package render

import (
	"goJsonDiff/pkg/render/style"
	"strings"
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
