// reloadCerts - a simple programm to send a SIGHUP signal to all
// TLSWebServer processes found, in order to trigger them to
// reload their certificates.
package main

import (
	"github.com/mitchellh/go-ps"
	"log"
	"os"
	"strings"
	"syscall"
)

func main() {
	// get a list of all processes
	var process *os.Process
	processes, err := ps.Processes()
	if err != nil {
		return
	}
	// iterate process list to find TLSWebServer processes
	for _, p := range processes {
		if strings.Contains(p.Executable(), "TLSWebServer") {
			process, err = os.FindProcess(p.Pid())
			if err != nil {
				return
			}
			log.Printf("found a 'TLSWebServer' process, PID = %d\n", p.Pid())
			// send HUP signal to process
			err = process.Signal(syscall.SIGHUP)
			if err != nil {
				return
			}
			log.Printf("sent HUP signal to PID: %d\n", p.Pid())
		}
	}
}
