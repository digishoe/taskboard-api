package model

// Board represents a kanban board.
type Board struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Columns []Column `json:"columns,omitempty"`
}

// Column represents a column within a board.
type Column struct {
	ID       int    `json:"id"`
	BoardID  int    `json:"board_id"`
	Name     string `json:"name"`
	Position int    `json:"position"`
	Tasks    []Task `json:"tasks,omitempty"`
}

// Task represents a task card within a column.
type Task struct {
	ID          int    `json:"id"`
	ColumnID    int    `json:"column_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Position    int    `json:"position"`
}
