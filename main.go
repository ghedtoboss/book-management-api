package main

import (
	"book-management-api/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB(dataSourceName string) *sql.DB {
	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	return db
}

func main() {
	initDB("root:Eses147852@tcp(127.0.0.1:3306)/book_management?parseTime=true")
	defer db.Close()
	fmt.Println("Veritabanına bağlanıldı.")

	r := mux.NewRouter()

	appHandler := &handlers.AppHandler{DB: db}

	//routes
	r.Handle("/books", appHandler.GetBooks()).Methods("GET")
	r.Handle("/books/{id}", appHandler.GetBook()).Methods("GET")
	r.Handle("/books", appHandler.CreateBook()).Methods("POST")
	r.Handle("/books/{id}", appHandler.UpdateBook()).Methods("PUT")
	r.Handle("/books/{id}", appHandler.DeleteBook()).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
