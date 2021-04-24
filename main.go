package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func handler(a, b string) {
	in := map[string]func(a string){
		"a": func(a string) {
			if a == "a" {
				fmt.Println("a")
			}
			if a == "b" {
				fmt.Println("b")
			}
		},
	}
	fn, exists := in[a]
	if !exists {
		fmt.Println("not exists")
	} else {
		fn(b)
	}

}

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	handleTerm(cancel)

	mux := http.NewServeMux()
	for path, route := range routes {
		mux.HandleFunc(path, route)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	httpServer := &http.Server{
		Addr:        "127.0.0.1:" + port,
		Handler:     mux,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}
	httpServer.RegisterOnShutdown(cancel)
	go func() {
		<-ctx.Done()
		gracefullCtx, cancelShutdown := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancelShutdown()

		if err := httpServer.Shutdown(gracefullCtx); err != nil {
			fmt.Printf("shutdown error: %v\n", err)
			defer os.Exit(1)
			return
		} else {
			fmt.Printf("gracefully stopped\n")
		}

	}()

	fmt.Printf("Listening on port %s", port)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	wg.Wait()
}

var routes = map[string]http.HandlerFunc{
	"/differentCase": func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Query().Get("foo") == "bar" {
			fmt.Fprint(w, "foobar")
			return
		}
		if req.URL.Query().Get("foo") == "bazzzz" {
			fmt.Fprint(w, "bazzzz")
			return
		}
		fmt.Fprint(w, "case2")
	},
	"/onlyPost": func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			fmt.Fprint(w, "only POST")
			return
		}
		if req.Method == http.MethodHead {
			fmt.Fprint(w, "only HEAD")
			return
		}
		fmt.Fprint(w, "hello")
	},
}

func handleTerm(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-signals
		fmt.Println("Caught", s)
		cancel()

		gracefulTimeout := 5 * time.Second
		select {
		case s = <-signals:
			fmt.Println("Exiting!", s)
		case <-time.After(gracefulTimeout):
			fmt.Println("Timeout on graceful leave. Exiting")
		}
	}()
}
