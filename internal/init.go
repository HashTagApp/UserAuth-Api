package internal

import (
	"context"
	"github.com/HashTagApp/UserAuth-Api/transport/http"
	"github.com/HashTagApp/hashlibry"
	"os"
	"os/signal"
	"syscall"
)

func Init() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	tCtx := hashlibry.Background()
	ctx, cancel := context.WithCancel(tCtx)
	defer cancel()
	println(ctx)

	http.Init(ctx)

	select {
	case <-sigs:
		http.StopServer(ctx)
	}

}
