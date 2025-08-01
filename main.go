package main

import (
	"fmt"
	"net/http"
	"sync"
)

var (
	tasks []string
	mu    sync.Mutex
)

func main() {
	fmt.Println("Starting Todo List App at http://localhost:8080")
	http.HandleFunc("/", showHome)
	http.HandleFunc("/add-task", addTask)
	http.HandleFunc("/clear-tasks", clearTasks)
	http.ListenAndServe(":8080", nil)
}

func showHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, `<html><body>`)
	fmt.Fprintln(w, `<h1>My Todo List</h1>`)

	mu.Lock()
	if len(tasks) == 0 {
		fmt.Fprintln(w, `<p><em>No tasks yet.</em></p>`)
	} else {
		fmt.Fprintln(w, `<ul>`)
		for _, task := range tasks {
			fmt.Fprintf(w, "<li>%s</li>", task)
		}
		fmt.Fprintln(w, `</ul>`)
	}
	mu.Unlock()

	// Form to add new task
	fmt.Fprintln(w, `
		<form action="/add-task" method="POST">
			<input type="text" name="task" placeholder="Enter new task" required>
			<input type="submit" value="Add Task">
		</form>
		<br>
		<form action="/clear-tasks" method="POST">
			<input type="submit" value="Clear All Tasks" style="color: red;">
		</form>
	</body></html>`)
}

func addTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		task := r.FormValue("task")
		if task != "" {
			mu.Lock()
			tasks = append(tasks, task)
			mu.Unlock()
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func clearTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		mu.Lock()
		tasks = nil
		mu.Unlock()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
