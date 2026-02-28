package goJsonDiff

import (
	"goJsonDiff/pkg/diff"
	"goJsonDiff/pkg/render"
	"goJsonDiff/pkg/types"
)

func ConsoleDiff(fromJson, toJson string) (string, error) {
	delta, err := DiffJsonStr(fromJson, toJson)
	if err != nil {
		return "", err
	}

	view := render.NewConsoleView()
	consoleStr, err := view.Render(delta)
	return consoleStr, err
}

func HtmlDiff(fromJson, toJson string) (string, error) {
	delta, err := DiffJsonStr(fromJson, toJson)
	if err != nil {
		return "", err
	}

	view := render.NewHtmlView()
	htmlStr, err := view.Render(delta)
	return htmlStr, err
}

func DiffJson(from, to types.JsonValue) (types.Delta, error) {
	return (&diff.JsonDiffer{}).Diff(from, to)
}

func DiffJsonStr(from, to string) (types.Delta, error) {
	return (&diff.JsonDiffer{}).DiffJsonStr(from, to)
}
