package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"snippetbox.gteruithi.com/internal/models"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
)

var DB *sql.DB

type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	var err error

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "host=localhost port=5432 user=postgre password=password dbname=snippetbox sslmode=disable", "PostgreSQL connection string")
	dbDriverName := flag.String("postgres", "postgres", "Database driver name")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	DB, err = openDB(*dbDriverName, *dsn)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer DB.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(DB)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: DB},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	logger.Info("Starting server", slog.String("addr", ":4000"))

	err = http.ListenAndServe(*addr, app.routes())

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}

func openDB(dsn, driverName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=web password=pass dbname=snippetbox sslmode=disable")

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
