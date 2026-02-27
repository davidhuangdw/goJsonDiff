package diff

import (
	"goJsonDiff/pkg/types"
	"reflect"
)

type DiffOpt struct {
	// for future features config
}

type JsonDiffer struct {
	// metadata
	Opt  DiffOpt
	path []any
}

func (jd *JsonDiffer) Diff(from, to types.JsonValue) (types.Delta, error) {
	hb := &HashBuilder{}
	fromHash, err := hb.buildHashTree(from)
	if err != nil {
		return nil, err
	}
	toHash, err := hb.buildHashTree(to)
	if err != nil {
		return nil, err
	}
	return jd.diffWithHash(from, to, fromHash, toHash)
}

func (jd *JsonDiffer) diffWithHash(from, to types.JsonValue, fromHash, toHash *HashTree) (types.Delta, error) {
	if reflect.TypeOf(from) != reflect.TypeOf(to) {
		return &types.DeltaLeaf{
			Op:        types.REPLACE,
			Value:     to,
			ValueFrom: from,
		}, nil
	}
	if fromHash.Hash == toHash.Hash {
		return &types.DeltaLeaf{
			Op:    types.EQUAL,
			Value: to,
		}, nil
	}

	switch fromVal := from.(type) {
	case string, float64, bool, nil:
		return &types.DeltaLeaf{
			Op:        types.REPLACE,
			Value:     to,
			ValueFrom: from,
		}, nil
	case []any: // LCS for array case
		delta, toVal := make([]types.Delta, 0), to.([]any)
		fromHashes, toHashes := fromHash.Children.([]*HashTree), toHash.Children.([]*HashTree)
		lcs := longestCommonSubHashes(fromHashes, toHashes)
		i, j, k, n, m := 0, 0, 0, len(fromHashes), len(toHashes)
		for i < n || j < m {
			switch {
			case i < n && j < m && fromHashes[i].Hash == toHashes[j].Hash:
				delta = append(delta, &types.DeltaLeaf{Op: types.EQUAL, Value: fromVal[i]})
				k++
				i++
				j++
			case i >= n || (i < n && k < len(lcs) && lcs[k] == fromHashes[i].Hash):
				delta = append(delta, &types.DeltaLeaf{Op: types.ADD, ValueFrom: toVal[j]})
				j++
			case j >= m || (j < m && k < len(lcs) && lcs[k] == toHashes[j].Hash):
				delta = append(delta, &types.DeltaLeaf{Op: types.REMOVE, ValueFrom: fromVal[i]})
				i++
			default: // no one equals lcs[k]
				delta = append(delta, &types.DeltaLeaf{Op: types.REPLACE, Value: toVal[j], ValueFrom: fromVal[i]})
				i++
				j++
			}
		}
		return delta, nil
	case map[string]any:
		delta := make(map[string]types.Delta)
		toVal := to.(map[string]any)
		for k, subFrom := range fromVal {
			if subTo, existSubTo := toVal[k]; existSubTo {
				subDelta, err := jd.diffWithHash(subFrom, subTo,
					fromHash.Children.(map[string]*HashTree)[k],
					toHash.Children.(map[string]*HashTree)[k])
				if err != nil {
					return nil, err
				}
				delta[k] = subDelta
			} else {
				delta[k] = &types.DeltaLeaf{Op: types.REMOVE, ValueFrom: subFrom}
			}
		}
		for k, subTo := range toVal {
			if _, existSubFrom := fromVal[k]; !existSubFrom {
				delta[k] = &types.DeltaLeaf{Op: types.ADD, Value: subTo}
			}
		}
		return delta, nil
	default:
		return nil, nil
	}
}

func longestCommonSubHashes(fromHashes, toHashes []*HashTree) []uint64 {
	n, m := len(fromHashes), len(toHashes)
	longest := make([][]int, n+1)
	for i := range longest {
		longest[i] = make([]int, m+1)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if fromHashes[i].Hash == toHashes[j].Hash {
				longest[i][j] = 1 + longest[i+1][j+1]
			} else {
				longest[i][j] = max(longest[i][j+1], longest[i+1][j])
			}
		}
	}

	var common []uint64
	i, j := 0, 0
	for i < n && j < m {
		if fromHashes[i].Hash == toHashes[j].Hash {
			common = append(common, fromHashes[i].Hash)
			i++
			j++
		} else if longest[i][j] == longest[i][j+1] {
			j++
		} else {
			i++
		}
	}
	return common
}
