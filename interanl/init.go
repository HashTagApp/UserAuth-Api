package interanl

import (
	"os"
	"os/signal"
	"syscall"
)

func Init() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//	tCtx := traceable_context.Background()
	//ctx, cancel := context.WithCancel(tCtx)
	//defer cancel()

}
