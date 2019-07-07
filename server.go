package main

import (
	"flag"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/buaazp/fasthttprouter"
	bx "github.com/claygod/Bxog"
	"github.com/dimfeld/httptreemux"
	"github.com/dinever/golf"
	restful "github.com/emicklei/go-restful"
	fasthttpSlashRouter "github.com/fasthttp/router"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi"
	ozzo "github.com/go-ozzo/ozzo-routing"
	"github.com/go-zoo/bone"
	"github.com/gocraft/web"
	"github.com/gorilla/mux"
	gowwwrouter "github.com/gowww/router"
	"github.com/gramework/gramework"
	"github.com/julienschmidt/httprouter"
	echov4 "github.com/labstack/echo/v4"
	"github.com/naoina/denco"
	"github.com/oxequa/fresh"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/razonyang/fastrouter"
	"github.com/savsgio/atreugo"
	"github.com/teambition/gear"
	"github.com/valyala/fasthttp"
	"github.com/xujiajun/gorouter"
	"go-web-framework-benchmark/pow"
	goji "goji.io"
	gojipat "goji.io/pat"
	"gopkg.in/baa.v1"
	"io"
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
	case "beego":
		startBeego()
	case "bone":
		startBone()
	case "bxog":
		startBxog()
	case "chi":
		startChi()
	case "denco":
		startDenco()
	case "echov4":
		startEchoV4()
	case "fasthttp":
		startFasthttp()
	case "fasthttprouter":
		startFastHTTPRouter()
	case "fasthttp/router":
		startFastHTTPSlashRouter()
	case "fasthttp-routing":
		startFastHTTPRouting()
	case "fastrouter":
		startFastRouter()
	case "fresh":
		startFresh()
	case "gear":
		startGear()
	case "gin":
		startGin()
	case "gocraftWeb":
		startGocraftWeb()
	case "goji":
		startGoji()
	case "gojsonrest":
		startGoJSONRest()
	case "golf":
		startGolf()
	case "gorestful":
		startGoRestful()
	case "gorilla":
		startGorilla()
	case "gorouter":
		startGorouter()
	case "go-ozzo":
		startGoozzo()
	case "gowww":
		startGowww()
	case "gramework":
		startGramework()
	case "httprouter":
		startHTTPRouter()
	case "httptreemux":
		starthttpTreeMux()
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
	http.HandleFunc("/hello", helloHandler)
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

// beego
func beegoHandler(ctx *context.Context) {
	if cpuBound {
		pow.Pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	ctx.WriteString(messageStr)

}
func startBeego() {
	beego.BConfig.RunMode = beego.PROD
	beego.BeeLogger.Close()
	mux := beego.NewControllerRegister()
	mux.Get("/hello", beegoHandler)
	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}

// bone
func startBone() {
	mux := bone.New()
	mux.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}

// bxog
func bxogHandler(w http.ResponseWriter, req *http.Request, r *bx.Router) {
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
func startBxog() {
	mux := bx.New()
	mux.Add("/hello", bxogHandler)
	mux.Start(":" + strconv.Itoa(port))
}

// chi
func startChi() {
	// Create a router instance.
	r := chi.NewRouter()
	// Register route handler.
	r.Get("/hello", helloHandler)
	// Start Chi.
	http.ListenAndServe(":"+strconv.Itoa(port), r)
}

// denco
func dencoHandler(w http.ResponseWriter, r *http.Request, params denco.Params) {
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
func startDenco() {
	mux := denco.NewMux()
	handler, _ := mux.Build([]denco.Handler{mux.GET("/hello", denco.HandlerFunc(dencoHandler))})
	http.ListenAndServe(":"+strconv.Itoa(port), handler)
}

// echov4-standard
func echov4Handler(c echov4.Context) error {
	if cpuBound {
		pow.Pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	return c.String(http.StatusOK, messageStr)
}
func startEchoV4() {
	e := echov4.New()
	e.GET("/hello", echov4Handler)
	e.Start(":" + strconv.Itoa(port))
}

// fasthttp
func fastHTTPRawHandler(ctx *fasthttp.RequestCtx) {
	if string(ctx.Method()) == "GET" {
		switch string(ctx.Path()) {
		case "/hello":
			if cpuBound {
				pow.Pow(target)
			} else {
				if sleepTime > 0 {
					time.Sleep(sleepTimeDuration)
				} else {
					runtime.Gosched()
				}
			}
			ctx.Write(message)
		default:
			ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		}
		return
	}
	ctx.Error("Unsupported method", fasthttp.StatusMethodNotAllowed)
}
func startFasthttp() {
	fasthttp.ListenAndServe(":"+strconv.Itoa(port), fastHTTPRawHandler)
}

// fasthttprouter
func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	if cpuBound {
		pow.Pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	ctx.Write(message)
}
func startFastHTTPRouter() {
	mux := fasthttprouter.New()
	mux.GET("/hello", fastHTTPHandler)
	fasthttp.ListenAndServe(":"+strconv.Itoa(port), mux.Handler)
}

// fasthttp/rrouter
func startFastHTTPSlashRouter() {
	mux := fasthttpSlashRouter.New()
	mux.GET("/hello", fastHTTPHandler)
	fasthttp.ListenAndServe(":"+strconv.Itoa(port), mux.Handler)
}

// fasthttp-routing
func fastHTTPRoutingHandler(c *routing.Context) error {
	if cpuBound {
		pow.Pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	c.Write(message)
	return nil
}
func startFastHTTPRouting() {
	mux := routing.New()
	mux.Get("/hello", fastHTTPRoutingHandler)
	fasthttp.ListenAndServe(":"+strconv.Itoa(port), mux.HandleRequest)
}

// fastrouter
func startFastRouter() {
	mux := fastrouter.New()
	mux.Get("/hello", helloHandler)
	mux.Prepare()
	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}

// fresh
func freshHandler(c fresh.Context) error {
	if cpuBound {
		pow.Pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	c.Response().Text(http.StatusOK, messageStr)
	return nil
}
func startFresh() {
	f := fresh.New()
	f.Config().Port = port
	f.GET("/hello", freshHandler)
	f.Start()
}

// gear
func startGear() {
	app := gear.New()
	router := gear.NewRouter()

	router.Get("/hello", func(c *gear.Context) error {
		if cpuBound {
			pow.Pow(target)
		} else {
			if sleepTime > 0 {
				time.Sleep(sleepTimeDuration)
			} else {
				runtime.Gosched()
			}
		}
		return c.HTML(http.StatusOK, messageStr)
	})

	app.UseHandler(router)
	app.Listen(":" + strconv.Itoa(port))
}

// gin
func ginHandler(c *gin.Context) {
	if cpuBound {
		pow.Pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	c.Writer.Write(message)
}
func startGin() {
	gin.SetMode(gin.ReleaseMode)
	mux := gin.New()
	mux.GET("/hello", ginHandler)
	mux.Run(":" + strconv.Itoa(port))
}

// gocraftWeb
type gocraftWebContext struct{}

func gocraftWebHandler(w web.ResponseWriter, r *web.Request) {
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
func startGocraftWeb() {
	mux := web.New(gocraftWebContext{})
	mux.Get("/hello", gocraftWebHandler)
	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}

// goji
func startGoji() {
	mux := goji.NewMux()
	mux.HandleFunc(gojipat.Get("/hello"), helloHandler)
	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}

// goJsonRest
func goJSONRestHandler(w rest.ResponseWriter, req *rest.Request) {
	if cpuBound {
		pow.Pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	iow := w.(io.Writer)
	iow.Write(message)
}
func startGoJSONRest() {
	api := rest.NewApi()
	router, _ := rest.MakeRouter(
		&rest.Route{HttpMethod: "GET", PathExp: "/hello", Func: goJSONRestHandler},
	)
	api.SetApp(router)
	http.ListenAndServe(":"+strconv.Itoa(port), api.MakeHandler())
}

// golf
func golfHandler(ctx *golf.Context) {
	if cpuBound {
		pow.Pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	ctx.Send(messageStr)
}
func startGolf() {
	app := golf.New()
	app.Get("/hello", golfHandler)
	app.Run(":" + strconv.Itoa(port))
}

// goRestful
func goRestfulHandler(r *restful.Request, w *restful.Response) {
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
func startGoRestful() {
	wsContainer := restful.NewContainer()
	ws := new(restful.WebService)
	ws.Route(ws.GET("/hello").To(goRestfulHandler))
	wsContainer.Add(ws)
	http.ListenAndServe(":"+strconv.Itoa(port), wsContainer)
}

// gorilla
func startGorilla() {
	mux := mux.NewRouter()
	mux.HandleFunc("/hello", helloHandler).Methods("GET")
	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}

// gorouter
func startGorouter() {
	mux := gorouter.New()
	mux.GET("/hello", helloHandler)
	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}

// go-ozzo
func ozzoHandler(c *ozzo.Context) error {
	if cpuBound {
		pow.Pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	c.Write(message)
	return nil
}
func startGoozzo() {
	r := ozzo.New()
	r.Get("/hello", ozzoHandler)
	http.ListenAndServe(":"+strconv.Itoa(port), r)
}

// gowww
func startGowww() {
	rt := gowwwrouter.New()
	rt.Handle("GET", "/hello", http.HandlerFunc(helloHandler))
	http.ListenAndServe(":"+strconv.Itoa(port), rt)
}

// Gramework
func grameworkHandler(ctx *gramework.Context) {
	if cpuBound {
		pow.Pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	ctx.WriteString(messageStr)
}
func startGramework() {
	gramework.SetEnv(gramework.PROD)
	app := gramework.New()
	app.GET("/hello", grameworkHandler)
	app.ListenAndServe(":" + strconv.Itoa(port))
}

// httprouter
func httpRouterHandler(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
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
func startHTTPRouter() {
	mux := httprouter.New()
	mux.GET("/hello", httpRouterHandler)
	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}

// httpTreeMux
func httpTreeMuxHandler(w http.ResponseWriter, _ *http.Request, vars map[string]string) {
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
func starthttpTreeMux() {
	mux := httptreemux.New()
	mux.GET("/hello", httpTreeMuxHandler)
	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}
