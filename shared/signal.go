package shared

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func HandleSignals(receivedFn func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		for sig := range sigs {
			log.Println("Signal received:", sig)

			receivedFn()

			return
		}
	}()
}
