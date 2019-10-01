package result

import (
	"log"

	"github.com/iv-p/apid/common/transaction"
)

// ConsoleWriter sends transaction results to stdout
type ConsoleWriter struct {
	success, failure int
}

// Header prints the header to the console
func (w *ConsoleWriter) Header() {
	log.Printf("Header\n")
}

// Write prints a result to the console
func (w *ConsoleWriter) Write(tx transaction.Result) {
	w.success++
	log.Printf("Result\n")
}

// Footer prints the footer to the console
func (w *ConsoleWriter) Footer() {
	log.Printf("success: %d failure: %d\n", w.success, w.failure)
}
