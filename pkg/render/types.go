package render

import . "github.com/davidhuangdw/goJsonDiff/pkg/types"

type DeltaView interface {
	Render(delta Delta) (string, error)
}
