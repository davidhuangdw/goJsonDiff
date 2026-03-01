package docs

import "github.com/davidhuangdw/goJsonDiff"

func ConsoleDiff(fromJson, toJson string) (string, error) {
	delta, err := goJsonDiff.DiffJsonStr(fromJson, toJson)
	if err != nil {
		return "", err
	}

	return goJsonDiff.RenderConsole(delta)
}

func HtmlDiff(fromJson, toJson string) (string, error) {
	delta, err := goJsonDiff.DiffJsonStr(fromJson, toJson)
	if err != nil {
		return "", err
	}

	return goJsonDiff.RenderHtml(delta)
}

func JsonPatchDiff(fromJson, toJson string) (string, error) {
	delta, err := goJsonDiff.DiffJsonStr(fromJson, toJson)
	if err != nil {
		return "", err
	}

	return goJsonDiff.RenderJsonPatch(delta)
}
