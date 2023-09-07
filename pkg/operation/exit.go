package operation

import (
	"fmt"
	"os"
	"time"

	"github.com/inconshreveable/mousetrap"
)

// Exit closes the client
func Exit(sig int) {
	if !mousetrap.StartedByExplorer() {
		os.Exit(sig)
	}
	fmt.Println("Close this window or press CTRL+C to exit.")
	for {
		time.Sleep(1 * time.Second)
	}
}
