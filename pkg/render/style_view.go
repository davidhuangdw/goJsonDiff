package render

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/davidhuangdw/goJsonDiff/pkg/render/style"
	. "github.com/davidhuangdw/goJsonDiff/pkg/types"
)

const KEY_SEP = ": "

type StyleView struct {
	*strings.Builder
	style.Format
	style.Style
}

func (st *StyleView) Render(delta Delta) (string, error) {
	st.Builder.Reset()
	if err := st.renderDeltaTree(delta, nil, "", false); err != nil {
		return "", err
	}

	s := st.String()
	if st.BreakLine != "\n" {
		s = strings.Join(strings.Split(s, "\n"), st.BreakLine)
	}
	return s, nil
}

func (st *StyleView) renderDeltaTree(delta Delta, path []any, key string, isLastItem bool) error {
	indent := st.Indent
	preIndents := repeat(indent, len(path))
	st.WriteString(preIndents)

	var comma string
	if key != "" {
		key = key + KEY_SEP
		if !isLastItem {
			comma = ","
		}
	}

	switch dlt := delta.(type) {
	case *DeltaLeaf:
		switch dlt.Op {
		case EQUAL, ADD, REMOVE:
			val, err := json.MarshalIndent(dlt.Value, preIndents, indent)
			if err != nil {
				return err
			}
			for _, s := range []string{key, string(val), comma} {
				st.write(dlt.Op, s)
			}
		case REPLACE:
			frVal, err := json.MarshalIndent(dlt.ValueFrom, preIndents, indent)
			if err != nil {
				return err
			}
			toVal, err := json.MarshalIndent(dlt.Value, preIndents, indent)
			if err != nil {
				return err
			}
			st.write(REPLACE, key)
			st.write(REMOVE, string(frVal))
			st.write(REPLACE, " => ")
			st.write(ADD, string(toVal))
			st.write(REPLACE, comma)
		default:
			return fmt.Errorf("unknown DeltaLeaf.Op: %#v", dlt)
		}
	case []ArrayItemDelta:
		st.WriteString(key + "[\n")
		for i, item := range dlt {
			if err := st.renderDeltaTree( // recur
				item.Delta, append(path, item.Index),
				renderArrKey(item), i == len(dlt)-1,
			); err != nil {
				return err
			}
		}
		st.WriteString(preIndents + "]" + comma)
	case map[string]Delta:
		st.WriteString(key + "{\n")
		i := 0
		for k, sub := range dlt {
			if err := st.renderDeltaTree( // recur
				sub, append(path, k),
				renderMapKey(k), i == len(dlt)-1,
			); err != nil {
				return err
			}
			i++
		}
		st.WriteString(preIndents + "}" + comma)
	default:
		return fmt.Errorf("unknown Delta type: %#v", dlt)
	}
	st.WriteString("\n")
	return nil
}

func (st *StyleView) write(op Operation, content string) {
	st.WriteString(st.Styled(op, content))
}

func renderArrKey(item ArrayItemDelta) string {
	i, j := item.IndexFrom, item.Index
	if i == j {
		return fmt.Sprintf("%v", j)
	}
	if i < 0 {
		return fmt.Sprintf("_~%v", j)
	}
	if j < 0 {
		return fmt.Sprintf("%v~_", i)
	}
	return fmt.Sprintf("%v~%v", i, j)
}

func renderMapKey(k string) string {
	return fmt.Sprintf("%#v", k)
}

func repeat(s string, k int) string {
	var b strings.Builder
	for i := 0; i < k; i++ {
		b.WriteString(s)
	}
	return b.String()
}
