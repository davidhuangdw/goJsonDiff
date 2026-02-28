package style

import . "goJsonDiff/pkg/types"

type Format struct {
	Indent    string
	BreakLine string
}

type Style interface {
	Styled(op Operation, content string) string
}
