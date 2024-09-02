package main

import (
	"database/sql"
	"encoding/json"
  "fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	. "github.com/workspace-9/erk"
  "git.jordanbonecutter.com/bibleplan/backend/calendar"
)

func main() {
  db := Try(sql.Open("sqlite3", "bibleplan.db"))
  setupDb(db)

  mux := &http.ServeMux{}
  mux.Handle("/", http.FileServer(http.Dir("./static/")))

  api := &http.ServeMux{}
  mux.Handle("/api/v1/", api)
  api.HandleFunc("PUT /api/v1/plan", func(w http.ResponseWriter, r *http.Request) {
    var request struct {
      StartDay time.Time `json:"startDay"`
      Email string `json:"email"`
    }

    Must(BadRequest(json.NewDecoder(r.Body).Decode(&request)))
    log.Println(request.StartDay, request.StartDay.UnixMilli())
    start := float64(request.StartDay.UnixMilli())/1000
    Try(db.Exec(`
      INSERT INTO subscribers(
        email, start_time
      ) VALUES ($1, $2) ON CONFLICT(email) DO UPDATE SET start_time = $3;
    `, request.Email, start, start))
  })
  api.HandleFunc("DELETE /api/v1/plan", func(w http.ResponseWriter, r *http.Request) {
    var request struct {
      Email string `json:"email"`
    }

    Must(BadRequest(json.NewDecoder(r.Body).Decode(&request)))
    Try(db.Exec(`
      DELETE FROM subscribers WHERE email = $1;
    `, request.Email))
  })
  api.HandleFunc("GET /api/v1/reading", func(w http.ResponseWriter, r *http.Request) {
    startDayStr := r.URL.Query().Get("start")
    log.Println(startDayStr)
    var startDay time.Time
    Must(BadRequest(json.Unmarshal([]byte(`"` + startDayStr + `"`), &startDay)))
    p, ok := calendar.MCheyne.On(startDay)
    if !ok {
      panic(BadRequest(fmt.Errorf("Enter a day within the last year")))
    }
    Must(json.NewEncoder(w).Encode(p))
  })
  Must(http.ListenAndServe(":7771", PanicHandler(cors{mux})))
}

type cors struct {
  inner http.Handler
}

func (c cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  w.Header().Add("Access-Control-Allow-Origin", "*")
  w.Header().Add("Access-Control-Allow-Headers", "*")
  w.Header().Add("Access-Control-Allow-Methods", "OPTIONS")
  w.Header().Add("Access-Control-Allow-Methods", "GET")
  w.Header().Add("Access-Control-Allow-Methods", "PUT")
  if r.Method == http.MethodOptions {
    w.WriteHeader(http.StatusNoContent)
    return
  }
  
  c.inner.ServeHTTP(w, r)
}

func setupDb(db *sql.DB) {
  Try(db.Exec(`
    CREATE TABLE IF NOT EXISTS subscribers(
      email varchar(255) unique primary key,
      start_time real
    );
  `))
}
