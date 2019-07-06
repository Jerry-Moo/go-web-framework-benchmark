package main

import (
	"fmt"
	// "io"
	"go-web-framework-benchmark/pow"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

var (
	port                  = 8080
	sleepTime             = 0
	cpuBound              bool
	target                = 15
	sleepTimeDuration     time.Duration
	samplingPointDuration time.Duration
	message               = []byte("hello world")
	messageStr            = "hello world"
	// seconds
	samplingPoint = 20
	// run web framework
	webFramework = "default"
)

func init() {
	args := os.Args
	argsLen := len(args)
	// args[1] web-framework
	if argsLen > 1 {
		webFramework = args[1]
	}
	// args[2] Processing Time
	if argsLen > 2 {
		sleepTime, _ = strconv.Atoi(args[2])
		if sleepTime == -1 {
			cpuBound = true
			sleepTime = 0
		}
	}

	if argsLen > 3 {
		port, _ = strconv.Atoi(args[3])
	}

	if argsLen > 4 {
		samplingPoint, _ = strconv.Atoi(args[4])
	}
	sleepTimeDuration = time.Duration(sleepTime) * time.Millisecond
	samplingPointDuration = time.Duration(samplingPoint) * time.Second

	go func() {
		time.Sleep(samplingPointDuration)
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		var u uint64 = 1024 * 1024
		fmt.Printf("TotalAlloc: %d\n", mem.TotalAlloc/u)
		fmt.Printf("Alloc: %d\n", mem.Alloc/u)
		fmt.Printf("HeapAlloc: %d\n", mem.HeapAlloc/u)
		fmt.Printf("HeapSys: %d\n", mem.HeapSys/u)
	}()
}

// server [default] [10] [8080]
func main() {
	switch webFramework {
	case "default":
		startDefaultMux()
	}
}

// default mux
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if cpuBound {
		pow.Pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	w.Write(message)
}
func startDefaultMux() {
	http.HandleFunc("/", helloHandler)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
