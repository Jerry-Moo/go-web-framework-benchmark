package main

import (
	"flag"
	"fmt"
	"github.com/savsgio/atreugo"
	"go-web-framework-benchmark/pow"
	"gopkg.in/baa.v1"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

var (
	port                  int
	sleepTime             int
	cpuBound              bool
	target                = 15
	sleepTimeDuration     time.Duration
	samplingPointDuration time.Duration
	message               = []byte("hello world")
	messageStr            = "hello world"
	// seconds
	samplingPoint int
	// run web framework
	webFramework string
)

func init() {
	// web-framework
	flag.StringVar(&webFramework, "wf", "default", "set test `web-framework`")
	flag.IntVar(&sleepTime, "s", 0, "set `sleep time`")
	flag.IntVar(&port, "p", 8080, "set `web` port")
	flag.IntVar(&samplingPoint, "sp", 20, "set `sampling point`")

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
	flag.Parse()
	if sleepTime == -1 {
		cpuBound = true
		sleepTime = 0
	}
	sleepTimeDuration = time.Duration(sleepTime) * time.Millisecond
	samplingPointDuration = time.Duration(samplingPoint) * time.Second

	switch webFramework {
	case "default":
		startDefaultMux()
	case "atreugo":
		startAtreugo()
	case "baa":
		startBaa()
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

// atreugo
func atreugoHandler(ctx *atreugo.RequestCtx) error {
	if cpuBound {
		pow.Pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	return ctx.TextResponseBytes(message)
}

func startAtreugo() {
	mux := atreugo.New(&atreugo.Config{Host: "127.0.0.1", Port: port})
	mux.Path("GET", "/hello", atreugoHandler)
	mux.ListenAndServe()
}

// baa
func baaHandler(ctx *baa.Context) {
	if cpuBound {
		pow.Pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	ctx.Text(http.StatusOK, message)
}
func startBaa() {
	mux := baa.New()
	mux.Get("/hello", baaHandler)
	mux.Run(":" + strconv.Itoa(port))
}
