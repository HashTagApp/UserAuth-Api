package http

import (
	"context"
	"fmt"
	"github.com/HashTagApp/UserAuth-Api/internal/app/server"
	log "github.com/HashTagApp/log"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

var httpServer *http.Server

func Init(ctx context.Context) {
	r := mux.NewRouter()
	port := 8080
	// Ping
	r.Handle(`/v1/ping`, server.Ping()).Methods(http.MethodGet)

	StartServer(ctx, port, r)

}

func StartServer(ctx context.Context, port int, r http.Handler) {
	running := make(chan interface{}, 1)

	httpServer = &http.Server{
		Addr:         fmt.Sprintf(`:%d`, port),
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 60,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func(ctx context.Context) {
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Fatal(ctx, `Cannot start web server`, err)
		}
		running <- `done`
	}(ctx)

	log.InfoContext(ctx, fmt.Sprintf(`HTTP router started on port [%d]`, port))

	<-running
}

func methodControl(method string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			h.ServeHTTP(w, r)
		} else {
			http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
		}
	})
}

func StopServer(ctx context.Context) {
	if err := httpServer.Shutdown(context.Background()); err != nil {

	}

}
