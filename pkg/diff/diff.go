package diff

import (
	"encoding/json"
	"fmt"
	. "goJsonDiff/pkg/types"
	"reflect"
)

type JsonDiffer struct {
	// to put config/metadata
}

func (jd *JsonDiffer) DiffJsonStr(fromStr, toStr string) (Delta, error) {
	var from, to JsonValue
	err := json.Unmarshal([]byte(fromStr), &from)
	if err != nil {
		return nil, fmt.Errorf("failed to Unmarshal(fromStr): %w", err)
	}

	err = json.Unmarshal([]byte(toStr), &to)
	if err != nil {
		return nil, fmt.Errorf("failed to Unmarshal(toStr): %w", err)
	}

	return jd.Diff(from, to)
}

func (jd *JsonDiffer) Diff(from, to JsonValue) (Delta, error) {
	hb := &HashBuilder{}
	fromHash, err := hb.buildHashTree(from, nil)
	if err != nil {
		return nil, err
	}
	toHash, err := hb.buildHashTree(to, nil)
	if err != nil {
		return nil, err
	}
	return jd.diff(from, to, fromHash, toHash)
}

func (jd *JsonDiffer) diff(from, to JsonValue, frHash, toHash *HashTree) (Delta, error) {
	if reflect.TypeOf(from) != reflect.TypeOf(to) {
		return &DeltaLeaf{
			Op:        REPLACE,
			Value:     to,
			ValueFrom: from,
		}, nil
	}

	if frHash.Hash == toHash.Hash {
		return &DeltaLeaf{
			Op:    EQUAL,
			Value: to,
		}, nil
	}

	switch fromVal := from.(type) {
	case string, float64, bool, nil:
		return &DeltaLeaf{
			Op:        REPLACE,
			Value:     to,
			ValueFrom: from,
		}, nil

	case []any: // LCS for array case
		delta, toVal := make([]ArrayItemDelta, 0), to.([]any)
		frHashes, toHashes := frHash.Children.([]*HashTree), toHash.Children.([]*HashTree)
		lcs := longestCommonSubHashes(frHashes, toHashes)
		i, j, k, n, m := 0, 0, 0, len(frHashes), len(toHashes)
		for i < n || j < m {
			switch {
			case i < n && j < m && frHashes[i].Hash == toHashes[j].Hash:
				delta = append(delta,
					ArrayItemDelta{
						Delta:     &DeltaLeaf{Op: EQUAL, Value: fromVal[i]},
						IndexFrom: i,
						Index:     j,
					},
				)
				k++
				i++
				j++
			case i >= n || (i < n && k < len(lcs) && lcs[k] == frHashes[i].Hash):
				delta = append(delta,
					ArrayItemDelta{
						Delta:     &DeltaLeaf{Op: ADD, Value: toVal[j]},
						IndexFrom: -1,
						Index:     j,
					},
				)
				j++
			case j >= m || (j < m && k < len(lcs) && lcs[k] == toHashes[j].Hash):
				delta = append(delta,
					ArrayItemDelta{
						Delta:     &DeltaLeaf{Op: REMOVE, Value: fromVal[i]},
						IndexFrom: i,
						Index:     -1,
					},
				)
				i++
			default: // when both are not lcs[k], replace(recur diff)
				subDelta, err := jd.diff(fromVal[i], toVal[j], frHashes[i], toHashes[j])
				if err != nil {
					return nil, err
				}
				delta = append(delta,
					ArrayItemDelta{
						Delta:     subDelta,
						IndexFrom: i,
						Index:     j,
					},
				)
				i++
				j++
			}
		}
		return delta, nil

	case map[string]any:
		delta := make(map[string]Delta)
		toVal := to.(map[string]any)
		for k, subFrom := range fromVal {
			if subTo, existSubTo := toVal[k]; existSubTo { // when same key, recur diff
				subDelta, err := jd.diff(subFrom, subTo,
					frHash.Children.(map[string]*HashTree)[k],
					toHash.Children.(map[string]*HashTree)[k])
				if err != nil {
					return nil, err
				}
				delta[k] = subDelta
			} else {
				delta[k] = &DeltaLeaf{Op: REMOVE, Value: subFrom}
			}
		}
		for k, subTo := range toVal {
			if _, existSubFrom := fromVal[k]; !existSubFrom {
				delta[k] = &DeltaLeaf{Op: ADD, Value: subTo}
			}
		}
		return delta, nil
	default:
		return nil, nil
	}
}
