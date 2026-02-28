package types

type Operation int
type JsonValue = any

const (
	EQUAL Operation = iota
	ADD
	REMOVE
	REPLACE
)

type Delta interface {
	// *DeltaLeaf | map[string]Delta | []ArrayItemDelta==[]{origin_index, Delta}
}

type DeltaLeaf struct {
	Op               Operation
	Value, ValueFrom JsonValue
}

type ArrayItemDelta struct { // to store original index
	Delta            Delta
	Index, IndexFrom int // index in original array
}

func IsEqual(delta Delta) bool {
	if leaf, ok := delta.(*DeltaLeaf); ok {
		return leaf.Op == EQUAL
	}
	return false
}
