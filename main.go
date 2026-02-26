package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./vulnerable.sqlite")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connection established")
}

func main() {
	// Serve static files
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/demo/auth", serveAuthDemo)
	http.HandleFunc("/demo/union", serveUnionDemo)
	http.HandleFunc("/demo/error", serveErrorDemo)
	http.HandleFunc("/demo/exfil", serveExfilDemo)

	// API endpoints
	http.HandleFunc("/api/login", handleLogin)
	http.HandleFunc("/api/search", handleSearch)
	http.HandleFunc("/api/profile", handleProfile)
	http.HandleFunc("/api/product", handleProduct)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	html, err := os.ReadFile("templates/home.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(html)
}

func serveAuthDemo(w http.ResponseWriter, r *http.Request) {
	html, err := os.ReadFile("templates/auth.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(html)
}

func serveUnionDemo(w http.ResponseWriter, r *http.Request) {
	html, err := os.ReadFile("templates/union.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(html)
}

func serveErrorDemo(w http.ResponseWriter, r *http.Request) {
	html, err := os.ReadFile("templates/error.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(html)
}

func serveExfilDemo(w http.ResponseWriter, r *http.Request) {
	html, err := os.ReadFile("templates/exfil.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(html)
}

// API Endpoints

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// INTENTIONALLY VULNERABLE - demonstrates SQL injection
	query := fmt.Sprintf("SELECT username, password FROM users WHERE username = '%s' AND password = '%s'", username, password)
	row := db.QueryRow(query)

	var dbUsername, dbPassword string
	err := row.Scan(&dbUsername, &dbPassword)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		fmt.Fprintf(w, `{"success": false, "error": "Invalid credentials"}`)
	} else {
		fmt.Fprintf(w, `{"success": true, "username": "%s"}`, dbUsername)
	}
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")

	// INTENTIONALLY VULNERABLE - demonstrates SQL injection
	query := fmt.Sprintf("SELECT product_id, name FROM products WHERE name LIKE '%s%%'", searchTerm)
	rows, err := db.Query(query)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		fmt.Fprintf(w, `{"error": "%s"}`, err.Error())
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			continue
		}
		results = append(results, map[string]interface{}{"id": id, "name": name})
	}

	if len(results) == 0 {
		fmt.Fprintf(w, `{"results": []}`)
	} else {
		// Simple JSON marshaling
		fmt.Fprint(w, `{"results": [`)
		for i, r := range results {
			if i > 0 {
				fmt.Fprint(w, `,`)
			}
			fmt.Fprintf(w, `{"id": %d, "name": "%s"}`, r["id"], r["name"])
		}
		fmt.Fprint(w, `]}`)
	}
}

func handleProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")

	// INTENTIONALLY VULNERABLE - demonstrates SQL injection
	query := fmt.Sprintf("SELECT id, username, email FROM users WHERE id = %s", userID)
	row := db.QueryRow(query)

	var id int
	var username, email string
	err := row.Scan(&id, &username, &email)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		fmt.Fprintf(w, `{"error": "%s"}`, err.Error())
	} else {
		fmt.Fprintf(w, `{"id": %d, "username": "%s", "email": "%s"}`, id, username, email)
	}
}

func handleProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("id")

	// INTENTIONALLY VULNERABLE - demonstrates SQL injection
	query := fmt.Sprintf("SELECT product_id, name, price FROM products WHERE product_id = %s", productID)
	rows, err := db.Query(query)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		fmt.Fprintf(w, `{"error": "%s"}`, err.Error())
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var product_id int
		var name string
		var price float64
		if err := rows.Scan(&product_id, &name, &price); err != nil {
			continue
		}
		results = append(results, map[string]interface{}{"id": product_id, "name": name, "price": price})
	}

	if len(results) == 0 {
		fmt.Fprint(w, `{"results": []}`)
	} else {
		fmt.Fprint(w, `{"results": [`)
		for i, r := range results {
			if i > 0 {
				fmt.Fprint(w, `,`)
			}
			fmt.Fprintf(w, `{"id": %d, "name": "%s", "price": %v}`, r["id"], r["name"], r["price"])
		}
		fmt.Fprint(w, `]}`)
	}
}
