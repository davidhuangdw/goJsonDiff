package diff

import (
	"fmt"
	"goJsonDiff/pkg/types"

	"github.com/zeebo/xxh3"
)

type HashTree struct {
	Hash     uint64
	Children any // nil | []*HashTree | map[string]*HashTree
}

type HashBuilder struct {
	// extra metadata
	path []any
}

func (hb *HashBuilder) buildHashTree(json types.JsonValue) (*HashTree, error) {
	switch jsonVal := json.(type) {
	case string, float64, bool, nil:
		return &HashTree{Hash: SimpleHash(jsonVal)}, nil
	case []types.JsonValue:
		chd := make([]*HashTree, len(jsonVal))
		codes := make([]uint64, len(jsonVal))
		for i, subVal := range jsonVal {
			hb.path = append(hb.path, i)
			subTree, err := hb.buildHashTree(subVal)
			if err != nil {
				return nil, err
			}
			hb.path = hb.path[:len(hb.path)-1] // pop
			chd[i] = subTree
			codes[i] = subTree.Hash
		}
		return &HashTree{Hash: SimpleHash(codes), Children: chd}, nil
	case map[string]types.JsonValue:
		chd := make(map[string]*HashTree, len(jsonVal))
		var hash uint64
		for k, subVal := range jsonVal {
			hb.path = append(hb.path, k)
			subTree, err := hb.buildHashTree(subVal)
			if err != nil {
				return nil, err
			}
			hb.path = hb.path[:len(hb.path)-1] // pop
			chd[k] = subTree
			hash ^= SimpleHash([2]uint64{SimpleHash(k), subTree.Hash}) // xor sub-hash to keep unordered
		}
		return &HashTree{Hash: hash, Children: chd}, nil
	default:
		return nil, fmt.Errorf("unexpected json type on path %#v: %#v", hb.path, json)
	}
}

func SimpleHash(value any) uint64 {
	return xxh3.HashString(fmt.Sprintf("%T|%+v", value, value))
}
