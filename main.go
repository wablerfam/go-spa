package main

import (
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
		ctx.SetState("status", "pushed")
		time.Sleep(1 * time.Second)
		ctx.SetState("status", "pushed\nloading")
		time.Sleep(1 * time.Second)
		ctx.SetState("status", "pushed\nloading\nended")
	})
}

func (a *myApp) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("go-app Demo"),
		app.Button().Text("Run").OnClick(a.onRun),
		app.Pre().Text(a.status),
	)
}

func main() {
	app.Route("/", func() app.Composer { return &myApp{} })
	app.RunWhenOnBrowser()

	http.Handle("/", &app.Handler{
		Name:        "go-app Demo",
		Description: "A simple go-app demo",
	})

	log.Println("Server starting on http://localhost:1323")
	if err := http.ListenAndServe(":1323", nil); err != nil {
		log.Fatal(err)
	}
}
