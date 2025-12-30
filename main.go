package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type myApp struct {
	app.Compo
	status string
}

func (a *myApp) OnMount(ctx app.Context) {
	ctx.ObserveState("status", &a.status)
}

func (a *myApp) onRun(ctx app.Context, e app.Event) {
	ctx.Async(func() {
		es := app.Window().Get("EventSource").New("/api/stream")
		es.Set("onmessage", app.FuncOf(func(this app.Value, args []app.Value) any {
			data := args[0].Get("data").String()
			ctx.Dispatch(func(ctx app.Context) {
				a.status += data + "\n"
				ctx.SetState("status", a.status)
			})
			return nil
		}))
		es.Set("onerror", app.FuncOf(func(this app.Value, args []app.Value) any {
			es.Call("close")
			return nil
		}))
	})
}

func (a *myApp) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("go-app Demo"),
		app.Button().Text("Run").OnClick(a.onRun),
		app.Pre().Text(a.status),
	)
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	messages := []string{"pushed", "loading", "ended"}
	for _, msg := range messages {
		fmt.Fprintf(w, "data: %s\n\n", msg)
		flusher.Flush()
		time.Sleep(1 * time.Second)
	}
}

func main() {
	app.Route("/", func() app.Composer { return &myApp{} })
	app.RunWhenOnBrowser()

	http.HandleFunc("/api/stream", sseHandler)
	http.Handle("/", &app.Handler{
		Name:        "go-app Demo",
		Description: "A simple go-app demo",
	})

	log.Println("Server starting on http://localhost:1323")
	if err := http.ListenAndServe(":1323", nil); err != nil {
		log.Fatal(err)
	}
}
