package render

import (
	"encoding/json"
	"fmt"
	"strings"

	. "github.com/davidhuangdw/goJsonDiff/pkg/types"
)

type PatchItem struct {
	Op    string `json:"op"`
	Path  string `json:"path,omitempty"`
	Value any    `json:"value,omitempty"`
}

type JsonPatchView struct { // follow RFC6902 patch format
	patches []PatchItem
}

func (vw *JsonPatchView) Render(delta Delta) (string, error) {
	vw.patches = make([]PatchItem, 0)
	if err := vw.render(delta, nil); err != nil {
		return "", err
	}

	if patches, err := json.MarshalIndent(vw.patches, "", "    "); err != nil {
		return "", err
	} else {
		return string(patches), nil
	}
}

func (vw *JsonPatchView) render(delta Delta, path []any) error {
	switch dlt := delta.(type) {
	case *DeltaLeaf:
		pathStr := buildPathStr(path)
		switch dlt.Op {
		case EQUAL:
			vw.patches = append(vw.patches, PatchItem{
				Op:    "test",
				Path:  pathStr,
				Value: dlt.Value,
			})
		case ADD:
			vw.patches = append(vw.patches, PatchItem{
				Op:    "add",
				Path:  pathStr,
				Value: dlt.Value,
			})
		case REMOVE:
			vw.patches = append(vw.patches, PatchItem{
				Op:    "remove",
				Path:  pathStr,
				Value: dlt.Value,
			})
		case REPLACE:
			vw.patches = append(vw.patches, PatchItem{
				Op:    "replace",
				Path:  pathStr,
				Value: dlt.Value,
			})
		default:
			return fmt.Errorf("unknown DeltaLeaf.Op: %#v", dlt)
		}
	case []ArrayItemDelta:
		for _, item := range dlt {
			// recur, prefer use from-side's index/path
			i := item.IndexFrom
			if i < 0 { // 'add' case
				i = item.Index
			}
			if err := vw.render(item.Delta, append(path, i)); err != nil {
				return err
			}
		}
	case map[string]Delta:
		for k, sub := range dlt {
			if err := vw.render(sub, append(path, k)); err != nil { // recur
				return err
			}
		}
	default:
		return fmt.Errorf("unknown Delta type: %#v", dlt)
	}
	return nil
}

func buildPathStr(path []any) string {
	b := &strings.Builder{}
	for _, v := range path {
		b.WriteString("/")
		b.WriteString(fmt.Sprintf("%v", v))
	}
	if len(path) == 0 {
		b.WriteString("/")
	}
	return b.String()
}
