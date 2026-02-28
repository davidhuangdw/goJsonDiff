package render

import . "goJsonDiff/pkg/types"

type DeltaView interface {
	Render(delta Delta) (string, error)
}
