package diff

import (
	"goJsonDiff/pkg/types"
	"testing"
)

func TestDiff(t *testing.T) {
	jd := &JsonDiffer{}

	tests := []struct {
		name    string
		from    types.JsonValue
		to      types.JsonValue
		checkOp types.Operation
		wantErr bool
	}{
		{
			name:    "equal primitives",
			from:    "hello",
			to:      "hello",
			checkOp: types.EQUAL,
		},
		{
			name:    "replace string",
			from:    "old",
			to:      "new",
			checkOp: types.REPLACE,
		},
		{
			name:    "replace different types",
			from:    "string",
			to:      123.0,
			checkOp: types.REPLACE,
		},
		{
			name:    "equal arrays",
			from:    []any{1.0, 2.0, 3.0},
			to:      []any{1.0, 2.0, 3.0},
			checkOp: types.EQUAL,
		},
		{
			name: "array with changes",
			from: []any{1.0, 2.0},
			to:   []any{1.0, 3.0},
		},
		{
			name:    "equal maps",
			from:    map[string]any{"key": "value"},
			to:      map[string]any{"key": "value"},
			checkOp: types.EQUAL,
		},
		{
			name: "map with add",
			from: map[string]any{"a": 1.0},
			to:   map[string]any{"a": 1.0, "b": 2.0},
		},
		{
			name: "map with remove",
			from: map[string]any{"a": 1.0, "b": 2.0},
			to:   map[string]any{"a": 1.0},
		},
		{
			name: "nested objects",
			from: map[string]any{
				"user": map[string]any{
					"name": "Alice",
					"age":  30.0,
				},
			},
			to: map[string]any{
				"user": map[string]any{
					"name": "Bob",
					"age":  30.0,
				},
			},
		},
		{
			name: "deeply nested with multiple changes",
			from: map[string]any{
				"level1": map[string]any{
					"level2": map[string]any{
						"a": 1.0,
						"b": 2.0,
					},
				},
			},
			to: map[string]any{
				"level1": map[string]any{
					"level2": map[string]any{
						"a": 1.0,
						"c": 3.0,
					},
				},
			},
		},
		{
			name: "array in object",
			from: map[string]any{
				"items": []any{1.0, 2.0, 3.0},
			},
			to: map[string]any{
				"items": []any{1.0, 3.0, 4.0},
			},
		},
		{
			name: "object in array",
			from: []any{
				map[string]any{"id": 1.0, "name": "A"},
				map[string]any{"id": 2.0, "name": "B"},
			},
			to: []any{
				map[string]any{"id": 1.0, "name": "A"},
				map[string]any{"id": 3.0, "name": "C"},
			},
		},
		{
			name: "LCS array insertion",
			from: []any{1.0, 3.0},
			to:   []any{1.0, 2.0, 3.0},
		},
		{
			name: "LCS array deletion",
			from: []any{1.0, 2.0, 3.0},
			to:   []any{1.0, 3.0},
		},
		{
			name: "LCS array reorder",
			from: []any{1.0, 2.0, 3.0},
			to:   []any{3.0, 2.0, 1.0},
		},
		{
			name: "empty to populated object",
			from: map[string]any{},
			to:   map[string]any{"a": 1.0, "b": 2.0},
		},
		{
			name:    "empty to empty array",
			from:    []any{},
			to:      []any{},
			checkOp: types.EQUAL,
		},
		{
			name: "complex nested structure",
			from: map[string]any{
				"users": []any{
					map[string]any{"id": 1.0, "tags": []any{"admin", "user"}},
					map[string]any{"id": 2.0, "tags": []any{"user"}},
				},
				"metadata": map[string]any{
					"version": "1.0",
					"count":   2.0,
				},
			},
			to: map[string]any{
				"users": []any{
					map[string]any{"id": 1.0, "tags": []any{"admin", "editor"}},
					map[string]any{"id": 3.0, "tags": []any{"user"}},
				},
				"metadata": map[string]any{
					"version": "2.0",
					"count":   2.0,
				},
			},
		},
		{
			name:    "nil values",
			from:    nil,
			to:      nil,
			checkOp: types.EQUAL,
		},
		{
			name:    "null to value",
			from:    nil,
			to:      "value",
			checkOp: types.REPLACE,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delta, err := jd.Diff(tt.from, tt.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("Diff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if delta == nil {
				t.Errorf("Diff() returned nil delta")
				return
			}
			if tt.checkOp != 0 {
				if leaf, ok := delta.(*types.DeltaLeaf); ok {
					if leaf.Op != tt.checkOp {
						t.Errorf("Diff() op = %v, want %v", leaf.Op, tt.checkOp)
					}
				}
			}
		})
	}
}

func TestDiff_ComplexValidation(t *testing.T) {
	jd := &JsonDiffer{}

	t.Run("nested map operations", func(t *testing.T) {
		from := map[string]any{
			"keep":   "same",
			"remove": "value",
			"update": "old",
		}
		to := map[string]any{
			"keep":   "same",
			"add":    "new",
			"update": "new",
		}

		delta, err := jd.Diff(from, to)
		if err != nil {
			t.Fatalf("Diff() error = %v", err)
		}

		deltaMap, ok := delta.(map[string]types.Delta)
		if !ok {
			t.Fatalf("Expected map[string]Delta, got %T", delta)
		}

		// Check EQUAL operation
		if leaf, ok := deltaMap["keep"].(*types.DeltaLeaf); ok {
			if leaf.Op != types.EQUAL {
				t.Errorf("keep: expected EQUAL, got %v", leaf.Op)
			}
		}

		// Check REMOVE operation
		if leaf, ok := deltaMap["remove"].(*types.DeltaLeaf); ok {
			if leaf.Op != types.REMOVE {
				t.Errorf("remove: expected REMOVE, got %v", leaf.Op)
			}
		}

		// Check REPLACE operation
		if leaf, ok := deltaMap["update"].(*types.DeltaLeaf); ok {
			if leaf.Op != types.REPLACE {
				t.Errorf("update: expected REPLACE, got %v", leaf.Op)
			}
		}

		// Check ADD operation
		if leaf, ok := deltaMap["add"].(*types.DeltaLeaf); ok {
			if leaf.Op != types.ADD {
				t.Errorf("add: expected ADD, got %v", leaf.Op)
			}
		}
	})

	t.Run("array LCS operations", func(t *testing.T) {
		from := []any{1.0, 2.0, 4.0}
		to := []any{1.0, 3.0, 4.0}

		delta, err := jd.Diff(from, to)
		if err != nil {
			t.Fatalf("Diff() error = %v", err)
		}

		deltaSlice, ok := delta.([]types.ArrayItemDelta)
		if !ok {
			t.Fatalf("Expected []Delta, got %T", delta)
		}

		if len(deltaSlice) == 0 {
			t.Error("Expected non-empty delta slice")
		}
	})

	t.Run("mixed nested structure", func(t *testing.T) {
		from := map[string]any{
			"data": map[string]any{
				"items": []any{
					map[string]any{"id": 1.0, "active": true},
					map[string]any{"id": 2.0, "active": false},
				},
			},
		}
		to := map[string]any{
			"data": map[string]any{
				"items": []any{
					map[string]any{"id": 1.0, "active": false},
					map[string]any{"id": 2.0, "active": false},
				},
			},
		}

		delta, err := jd.Diff(from, to)
		if err != nil {
			t.Fatalf("Diff() error = %v", err)
		}

		deltaMap, ok := delta.(map[string]types.Delta)
		if !ok {
			t.Fatalf("Expected map[string]Delta, got %T", delta)
		}

		// Navigate to nested data.items
		dataMap, ok := deltaMap["data"].(map[string]types.Delta)
		if !ok {
			t.Fatalf("Expected data to be map[string]Delta")
		}

		itemsSlice, ok := dataMap["items"].([]types.ArrayItemDelta)
		if !ok {
			t.Fatalf("Expected items to be []Delta")
		}

		if len(itemsSlice) == 0 {
			t.Error("Expected items to have deltas")
		}
	})
}
