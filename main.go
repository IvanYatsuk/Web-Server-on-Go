package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Book struct {
	ISBN      string `json:"isbn"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Year      int    `json:"year"`
	Pages     int    `json:"pages"`
	Publisher string `json:"publisher"`
	Language  string `json:"language"`
	Genre     string `json:"genre"` // Один жанр
}

type Reader struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	PhoneNumber  string    `json:"phone_number"`
	Registration time.Time `json:"registration"`
	Notes        string    `json:"notes"`
}

var db *sql.DB

func initDB() {
	var err error
	dbURL := "postgres://postgres:ivanyatsuk@db:5432/library"

	db, err = sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Unable to reach the database: %v\n", err)
	}

	_, err = db.ExecContext(ctx, `
    CREATE TABLE IF NOT EXISTS books (
        isbn TEXT PRIMARY KEY,
        title TEXT,
        author TEXT,
        year INT,
        pages INT,
        publisher TEXT,
        language TEXT,
        genre TEXT
    );

    CREATE TABLE IF NOT EXISTS readers (
        id SERIAL PRIMARY KEY,
        first_name TEXT,
        last_name TEXT,
        phone_number TEXT,
        registration TIMESTAMP,
        notes TEXT
    );
    `)
	if err != nil {
		log.Fatalf("Unable to create tables: %v\n", err)
	}
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("1.0.0"))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT isbn, title, author, year, pages, publisher, language, genre FROM books")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ISBN, &book.Title, &book.Author, &book.Year, &book.Pages, &book.Publisher, &book.Language, &book.Genre); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	if err := json.NewEncoder(w).Encode(books); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getBook(w http.ResponseWriter, r *http.Request) {
	isbn := r.URL.Query().Get("isbn")
	if isbn == "" {
		http.Error(w, "Missing isbn parameter", http.StatusBadRequest)
		return
	}

	var book Book
	err := db.QueryRow("SELECT isbn, title, author, year, pages, publisher, language, genre FROM books WHERE isbn = $1", isbn).Scan(
		&book.ISBN, &book.Title, &book.Author, &book.Year, &book.Pages, &book.Publisher, &book.Language, &book.Genre)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO books (isbn, title, author, year, pages, publisher, language, genre) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		book.ISBN, book.Title, book.Author, book.Year, book.Pages, book.Publisher, book.Language, book.Genre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE books SET title = $2, author = $3, year = $4, pages = $5, publisher = $6, language = $7, genre = $8 WHERE isbn = $1",
		book.ISBN, book.Title, book.Author, book.Year, book.Pages, book.Publisher, book.Language, book.Genre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	isbn := r.URL.Query().Get("isbn")
	if isbn == "" {
		http.Error(w, "Missing isbn parameter", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM books WHERE isbn = $1", isbn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getReaders(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, first_name, last_name, phone_number, registration, notes FROM readers")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var readers []Reader
	for rows.Next() {
		var reader Reader
		if err := rows.Scan(&reader.ID, &reader.FirstName, &reader.LastName, &reader.PhoneNumber, &reader.Registration, &reader.Notes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		readers = append(readers, reader)
	}

	if err := json.NewEncoder(w).Encode(readers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getReader(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	var reader Reader
	err = db.QueryRow("SELECT id, first_name, last_name, phone_number, registration, notes FROM readers WHERE id = $1", id).Scan(
		&reader.ID, &reader.FirstName, &reader.LastName, &reader.PhoneNumber, &reader.Registration, &reader.Notes)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(reader); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createReader(w http.ResponseWriter, r *http.Request) {
	var reader Reader
	if err := json.NewDecoder(r.Body).Decode(&reader); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO readers (first_name, last_name, phone_number, registration, notes) VALUES ($1, $2, $3, $4, $5)",
		reader.FirstName, reader.LastName, reader.PhoneNumber, time.Now(), reader.Notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func updateReader(w http.ResponseWriter, r *http.Request) {
	var reader Reader
	if err := json.NewDecoder(r.Body).Decode(&reader); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE readers SET first_name = $2, last_name = $3, phone_number = $4, notes = $5 WHERE id = $1",
		reader.ID, reader.FirstName, reader.LastName, reader.PhoneNumber, reader.Notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteReader(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM readers WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/version", getVersion)
	http.HandleFunc("/books", getBooks)
	http.HandleFunc("/books/get", getBook)
	http.HandleFunc("/books/create", createBook)
	http.HandleFunc("/books/update", updateBook)
	http.HandleFunc("/books/delete", deleteBook)

	http.HandleFunc("/readers", getReaders)
	http.HandleFunc("/readers/get", getReader)
	http.HandleFunc("/readers/create", createReader)
	http.HandleFunc("/readers/update", updateReader)
	http.HandleFunc("/readers/delete", deleteReader)

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
