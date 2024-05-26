package main

import (
	"Web-Server-on-Go/structures"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	books   map[string]structures.Book = make(map[string]structures.Book)
	readers map[int]structures.Reader  = make(map[int]structures.Reader)
	mu      sync.Mutex
	nextID  int = 1
)

// Handlers for books
func getBooks(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	var bookList []structures.Book
	for _, book := range books {
		bookList = append(bookList, book)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookList)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	isbn := r.URL.Query().Get("isbn")
	mu.Lock()
	defer mu.Unlock()
	book, exists := books[isbn]
	if !exists {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book structures.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	books[book.ISBN] = book
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	isbn := r.URL.Query().Get("isbn")
	mu.Lock()
	defer mu.Unlock()
	if _, exists := books[isbn]; !exists {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	var book structures.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	books[isbn] = book
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	isbn := r.URL.Query().Get("isbn")
	mu.Lock()
	defer mu.Unlock()
	if _, exists := books[isbn]; !exists {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	delete(books, isbn)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Book deleted"})
}

// Handlers for readers
func getReaders(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	var readerList []structures.Reader
	for _, reader := range readers {
		readerList = append(readerList, reader)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(readerList)
}

func getReader(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	readerID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	reader, exists := readers[readerID]
	if !exists {
		http.Error(w, "Reader not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reader)
}

func createReader(w http.ResponseWriter, r *http.Request) {
	var reader structures.Reader
	if err := json.NewDecoder(r.Body).Decode(&reader); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	reader.RegistrationDate = time.Now()
	mu.Lock()
	defer mu.Unlock()
	reader.ID = nextID
	nextID++
	readers[reader.ID] = reader
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reader)
}

func updateReader(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	readerID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var reader structures.Reader
	if err := json.NewDecoder(r.Body).Decode(&reader); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	if _, exists := readers[readerID]; !exists {
		http.Error(w, "Reader not found", http.StatusNotFound)
		return
	}
	readers[readerID] = reader
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reader)
}

func deleteReader(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	readerID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	if _, exists := readers[readerID]; !exists {
		http.Error(w, "Reader not found", http.StatusNotFound)
		return
	}
	delete(readers, readerID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Reader deleted"})
}

// Main function
func main() {
	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getBooks(w, r)
		case http.MethodPost:
			createBook(w, r)
		case http.MethodPut:
			updateBook(w, r)
		case http.MethodDelete:
			deleteBook(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/books/book", func(w http.ResponseWriter, r *http.Request) {
		getBook(w, r)
	})

	http.HandleFunc("/readers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getReaders(w, r)
		case http.MethodPost:
			createReader(w, r)
		case http.MethodPut:
			updateReader(w, r)
		case http.MethodDelete:
			deleteReader(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/readers/reader", func(w http.ResponseWriter, r *http.Request) {
		getReader(w, r)
	})

	log.Printf("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
