package main

import (
	"database/sql"
	"fmt"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"net/http"
)

type Handler struct {
	db  *sql.DB
	log *zap.Logger
}

func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hello, World!")

}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("autorizon", "success")
	fmt.Fprintf(w, "success")
}

type HostSwitch map[string]http.Handler

func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if i := recover(); i != nil {
		fmt.Println("Recover panic")
	}
	if handler := hs[r.Host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		// Handle host names for which no handler is registered
		http.Error(w, "Forbidden", 403) // Or Redirect?
	}
}

func (h *Handler) connect_db() {
	data := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", username, password, dbname)
	db, err := sql.Open("postgres", data)
	if err != nil {
		h.log.Warn("error in sql.Open",
			zap.Error(err))
	}
	err = db.Ping()
	if err != nil {
		h.log.Warn("error in db.Ping",
			zap.Error(err))
	}
	h.db = db
}

func (h *Handler) create_logger() {
	cnf := zap.NewProductionConfig()
	cnf.OutputPaths = []string{"errors.log"}
	logger, err := cnf.Build()
	if err != nil {
		log.Println("error in build zap")
	}
	h.log = logger
}

func main() {
	var h Handler
	h.create_logger()
	h.connect_db()
	router := httprouter.New()
	router.GET("/", h.MainPage)
	router.GET("/login", h.LoginPage)
	hs := make(HostSwitch)
	hs["localhost:8080"] = router
	fmt.Println("starting server at :8080")
	http.ListenAndServe("localhost:8080", hs)
}
