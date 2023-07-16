package models

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func HandlerSig(app *App) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-sig
	fmt.Println("\nClose Program...")
	killPid(&app.Program)
	fmt.Println("Thank see you next time...")
	os.Exit(1)
}
