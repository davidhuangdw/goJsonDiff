package types

type Operation int
type JsonValue = any

const (
	EQUAL Operation = iota
	ADD
	REMOVE
	REPLACE
	//MOVE
	//IGNORED
)

type DeltaLeaf struct {
	Op        Operation
	Value     JsonValue
	ValueFrom JsonValue
	//KeyFrom   any // for 'MOVE' case, nil | int | string,
}

type Delta interface {
	// *DeltaLeaf | []Delta | map[string]Delta
}
