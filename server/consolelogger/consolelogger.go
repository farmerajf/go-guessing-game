package consolelogger

import "fmt"

// Consolelogger represents an instance of a logger that writes to the console
type Consolelogger struct {
}

// Log prints the log message to the console
func (l *Consolelogger) Log(message string) {
	fmt.Println(message)
}
