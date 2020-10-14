package main

import (
	"fmt"
	"github.com/jhdriver/UWORLD/config"
	"github.com/jhdriver/UWORLD/node"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

// interruptSignals defines the default signals to catch in order to do a proper
// shutdown.  This may be modified during init depending on the platform.
var interruptSignals = []os.Signal{
	os.Interrupt,
	os.Kill,
	syscall.SIGINT,
	syscall.SIGTERM,
}

func main() {
	// Initialize the goroutine count,  Use all processor cores.
	runtime.GOMAXPROCS(runtime.NumCPU())
	// Work around defer not working after os.Exit()
	if err := UWDMain(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// main start the UWD node function
func UWDMain() error {
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Load configuration and parse command line.  This function also
	// initializes logging and configures it accordingly.
	config, err := config.LoadConfig()
	if err != nil {
		return err
	}
	//Initialize the UWD node
	node, err := node.NewNode(config)
	if err != nil {
		return err
	}

	//Start UWD node
	if err := node.Start(); err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, interruptSignals...)

	// Listen for initial shutdown signal and close the returned
	// channel to notify the caller.
	go func() {
		<-c
		node.Stop()
		close(c)
		os.Exit(0)
	}()
	wg.Wait()
	return nil
}
