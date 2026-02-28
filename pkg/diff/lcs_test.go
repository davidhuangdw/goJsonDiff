package diff

import (
	"reflect"
	"testing"
)

func TestLongestCommonSublist(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []int
		want []int
	}{
		{
			name: "empty slices",
			a:    []int{},
			b:    []int{},
			want: []int{},
		},
		{
			name: "no common elements",
			a:    []int{1, 2, 3},
			b:    []int{4, 5, 6},
			want: []int{},
		},
		{
			name: "all common in order",
			a:    []int{1, 2, 3},
			b:    []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			name: "partial common",
			a:    []int{1, 2, 3, 4},
			b:    []int{2, 4},
			want: []int{2, 4},
		},
		{
			name: "reordered",
			a:    []int{1, 2, 3},
			b:    []int{3, 2, 1},
			want: []int{3},
		},
		{
			name: "single element match",
			a:    []int{1},
			b:    []int{1},
			want: []int{1},
		},
		{
			name: "one empty",
			a:    []int{1, 2, 3},
			b:    []int{},
			want: []int{},
		},
		{
			name: "classic LCS example",
			a:    []int{1, 2, 3, 4, 5},
			b:    []int{2, 3, 5},
			want: []int{2, 3, 5},
		},
		{
			name: "interleaved common",
			a:    []int{1, 3, 5, 7, 9},
			b:    []int{2, 3, 4, 7, 8, 9},
			want: []int{3, 7, 9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := longestCommonSublist(tt.a, tt.b)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("longestCommonSublist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLongestCommonSublist_Strings(t *testing.T) {
	tests := []struct {
		name string
		a    []string
		b    []string
		want []string
	}{
		{
			name: "string sequence",
			a:    []string{"a", "b", "c", "d"},
			b:    []string{"b", "d"},
			want: []string{"b", "d"},
		},
		{
			name: "words",
			a:    []string{"hello", "world", "foo"},
			b:    []string{"world", "bar", "foo"},
			want: []string{"world", "foo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := longestCommonSublist(tt.a, tt.b)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("longestCommonSublist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLongestCommonSubHashes(t *testing.T) {
	tests := []struct {
		name      string
		fromNodes []*HashTree
		toNodes   []*HashTree
		want      []uint64
	}{
		{
			name:      "empty",
			fromNodes: []*HashTree{},
			toNodes:   []*HashTree{},
			want:      []uint64{},
		},
		{
			name: "all matching",
			fromNodes: []*HashTree{
				{Hash: 1},
				{Hash: 2},
				{Hash: 3},
			},
			toNodes: []*HashTree{
				{Hash: 1},
				{Hash: 2},
				{Hash: 3},
			},
			want: []uint64{1, 2, 3},
		},
		{
			name: "partial match",
			fromNodes: []*HashTree{
				{Hash: 1},
				{Hash: 2},
				{Hash: 3},
			},
			toNodes: []*HashTree{
				{Hash: 2},
				{Hash: 3},
			},
			want: []uint64{2, 3},
		},
		{
			name: "no match",
			fromNodes: []*HashTree{
				{Hash: 1},
				{Hash: 2},
			},
			toNodes: []*HashTree{
				{Hash: 3},
				{Hash: 4},
			},
			want: []uint64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := longestCommonSubHashes(tt.fromNodes, tt.toNodes)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("longestCommonSubHashes() = %v, want %v", got, tt.want)
			}
		})
	}
}
