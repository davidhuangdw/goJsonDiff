package goJsonDiff

import (
	"github.com/davidhuangdw/goJsonDiff/pkg/diff"
	"github.com/davidhuangdw/goJsonDiff/pkg/render"
	. "github.com/davidhuangdw/goJsonDiff/pkg/types"
)

func DiffJson(from, to JsonValue) (Delta, error) {
	return (&diff.JsonDiffer{}).Diff(from, to)
}

func DiffJsonStr(from, to string) (Delta, error) {
	return (&diff.JsonDiffer{}).DiffJsonStr(from, to)
}

func RenderConsole(delta Delta) (string, error) {
	return render.NewConsoleView().Render(delta)
}
func RenderHtml(delta Delta) (string, error) {
	return render.NewHtmlView().Render(delta)
}
func RenderJsonPatch(delta Delta) (string, error) {
	return render.NewJsonPatchView().Render(delta)
}
