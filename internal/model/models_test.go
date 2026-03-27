package model

import (
	"encoding/json"
	"testing"
)

func TestBoardJSONOmitsEmptyColumns(t *testing.T) {
	b := Board{ID: 1, Name: "Test"}
	data, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	json.Unmarshal(data, &m)
	if _, ok := m["columns"]; ok {
		t.Error("expected columns to be omitted when empty")
	}
}

func TestBoardJSONIncludesColumns(t *testing.T) {
	b := Board{
		ID:   1,
		Name: "Test",
		Columns: []Column{
			{ID: 1, BoardID: 1, Name: "To Do", Position: 0},
		},
	}
	data, _ := json.Marshal(b)
	var m map[string]interface{}
	json.Unmarshal(data, &m)
	cols, ok := m["columns"]
	if !ok {
		t.Fatal("expected columns in JSON")
	}
	arr := cols.([]interface{})
	if len(arr) != 1 {
		t.Errorf("expected 1 column, got %d", len(arr))
	}
}

func TestColumnJSONOmitsEmptyTasks(t *testing.T) {
	c := Column{ID: 1, BoardID: 1, Name: "Col", Position: 0}
	data, _ := json.Marshal(c)
	var m map[string]interface{}
	json.Unmarshal(data, &m)
	if _, ok := m["tasks"]; ok {
		t.Error("expected tasks to be omitted when empty")
	}
}

func TestTaskJSONRoundTrip(t *testing.T) {
	tk := Task{ID: 5, ColumnID: 2, Title: "Do thing", Description: "Details", Position: 3}
	data, _ := json.Marshal(tk)
	var out Task
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.ID != tk.ID || out.ColumnID != tk.ColumnID || out.Title != tk.Title ||
		out.Description != tk.Description || out.Position != tk.Position {
		t.Errorf("round-trip mismatch: got %+v", out)
	}
}
