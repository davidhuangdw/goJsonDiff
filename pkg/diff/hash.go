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
}

func (hb *HashBuilder) buildHashTree(json types.JsonValue, path []any) (*HashTree, error) {
	switch jsonVal := json.(type) {
	case string, float64, bool, nil:
		return &HashTree{Hash: SimpleHash(jsonVal)}, nil
	case []types.JsonValue:
		chd := make([]*HashTree, len(jsonVal))
		codes := make([]uint64, len(jsonVal))
		for i, subVal := range jsonVal {
			subTree, err := hb.buildHashTree(subVal, append(path, i))
			if err != nil {
				return nil, err
			}
			chd[i] = subTree
			codes[i] = subTree.Hash
		}
		return &HashTree{Hash: SimpleHash(codes), Children: chd}, nil
	case map[string]types.JsonValue:
		chd := make(map[string]*HashTree, len(jsonVal))
		var hash uint64
		for k, subVal := range jsonVal {
			subTree, err := hb.buildHashTree(subVal, append(path, k))
			if err != nil {
				return nil, err
			}
			chd[k] = subTree
			hash ^= SimpleHash([2]uint64{SimpleHash(k), subTree.Hash}) // xor sub-hash to keep unordered
		}
		return &HashTree{Hash: hash, Children: chd}, nil
	default:
		return nil, fmt.Errorf("unexpected json type on path %#v: %#v", path, json)
	}
}

func SimpleHash(value any) uint64 {
	return xxh3.HashString(fmt.Sprintf("%T|%+v", value, value))
}
