package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/vitorbass93/challenge1/server/provider"
	"github.com/vitorbass93/challenge1/server/repository"
)

func main() {
	db, err := getDatabaseConnection("file:challenge1.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatalf("getDatabaseConnection() failed with err: %+v\n", err)
	}
	defer db.Close()
	dolarRepository := repository.NewDolarRepository(db)
	if err != nil {
		log.Fatalf("repository.NewDolarRepository() failed with err: %+v\n", err)
	}
	dolarProvider := provider.NewDolarProvider("https://economia.awesomeapi.com.br/json/last/USD-BRL", dolarRepository)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health-check", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("GET /cotacao", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		dolar, err := dolarProvider.GetDolar()
		if err != nil {
			log.Printf("dolarProvider.GetDolar() failed with err: %+v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dolar)
	})
	startServer(":8080", mux)
}

func getDatabaseConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("DROP TABLE IF EXISTS usd_brl_quotations")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS usd_brl_quotations (id INTEGER PRIMARY KEY, code TEXT, codein TEXT, name TEXT, high TEXT, low TEXT, var_bid TEXT, pct_change TEXT, bid TEXT, ask TEXT, timestamp TEXT, create_date TEXT)")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func startServer(addr string, mux *http.ServeMux) {
	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}
	go func() {
		log.Println("serving new connections on :8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server.ListenAndServe() failed with err: %+v\n", err)
		}
		log.Println("stopped serving new connections")
	}()
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
	<-shutdownChan
	log.Println("server shutting down")
	server.SetKeepAlivesEnabled(false)
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown failed with err: %+v\n", err)
	}
	log.Println("server stopped")
}
