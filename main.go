package main

import (
	"log"
	"net/http"
	"os"

	"github.com/digishoe/taskboard-api/internal/handler"
	"github.com/digishoe/taskboard-api/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "taskboard.db"
	}

	s, err := store.New(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer s.Close()

	boards := handler.NewBoardHandler(s)
	columns := handler.NewColumnHandler(s)
	tasks := handler.NewTaskHandler(s)

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Board routes
	mux.HandleFunc("GET /api/boards", boards.List)
	mux.HandleFunc("POST /api/boards", boards.Create)
	mux.HandleFunc("GET /api/boards/{id}", boards.Get)

	// Column routes
	mux.HandleFunc("POST /api/boards/{id}/columns", columns.Create)

	// Task routes
	mux.HandleFunc("POST /api/columns/{id}/tasks", tasks.Create)
	mux.HandleFunc("PUT /api/tasks/{id}", tasks.Update)
	mux.HandleFunc("DELETE /api/tasks/{id}", tasks.Delete)

	log.Printf("taskboard-api listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsMiddleware(mux)))
}

// corsMiddleware adds CORS headers for local frontend development.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
