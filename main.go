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

var h Handler

func init() {
	cnf := zap.NewProductionConfig()
	cnf.OutputPaths = []string{"errors.log"}
	logger, err := cnf.Build()
	if err != nil {
		log.Println("error in build zap")
	}
	data := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", username, password, dbname)
	db, err := sql.Open("postgres", data)
	if err != nil {
		logger.Warn("error in sql.Open",
			zap.Error(err))
	}
	err = db.Ping()
	if err != nil {
		logger.Warn("error in db.Ping",
			zap.Error(err))
	}
	h.db = db
	h.log = logger
}

func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hello, World!")
	http.Redirect(w, r, "/login", http.StatusFound)
	// добавляем этот код, чтобы сохранить заголовок
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

func main() {
	router := httprouter.New()
	router.GET("/", h.MainPage)
	router.GET("/login", h.LoginPage)
	hs := make(HostSwitch)
	hs["localhost:8080"] = router
	fmt.Println("start server", http.ListenAndServe("localhost:8080", hs))
}
