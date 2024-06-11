package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Book struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

func (db *AppHandler) GetBooks() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var books []Book
		query := "SELECT id, title, author, description FROM books"
		rows, err := db.DB.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var book Book
			if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Description); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			books = append(books, book)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(books)
	})
}

func (db *AppHandler) GetBook() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var book Book
		query := "SELECT id, title, author, description FROM books WHERE id = ?"
		if err := db.DB.QueryRow(query, id).Scan(&book.ID, &book.Title, &book.Author, &book.Description); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)
	})
}

func (db *AppHandler) CreateBook() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var book Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		res, err := db.DB.Exec("INSERT INTO books (title, author, description) VALUES (?, ?, ?)", book.Title, book.Author, book.Description)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lastInserId, err := res.LastInsertId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		book.ID = int(lastInserId)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)
	})
}

func (db *AppHandler) UpdateBook() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var book Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err := db.DB.Exec("UPDATE books SET title = ?, author = ?, description = ? WHERE id = ?", book.Title, book.Author, book.Description, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(book)
	})
}

func (db *AppHandler) DeleteBook() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		_, err := db.DB.Exec("DELETE FROM books WHERE id = ?", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "kitap silindi"})
	})
}
