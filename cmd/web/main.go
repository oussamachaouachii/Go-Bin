package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"snippetbox.oussama.com/internal/models"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
	template map[string]*template.Template
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	user := "root"
	password := "root"
	host := "localhost"
	port := 3306
	dbName := "snippetbox"

	// Constructing the connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", user, password, host, port, dbName)

	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	template, err := cache()
	if err != nil {
		return
	}

	app := application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
		template: template,
	}

	srv := &http.Server{
		Addr:     ":9000",
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Println("Starting server on http://localhost:9000")
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
