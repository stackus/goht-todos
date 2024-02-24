package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stackus/errors"

	"github.com/stackus/goht-todos/assets"
	"github.com/stackus/goht-todos/domain"
	"github.com/stackus/goht-todos/stores"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// begin flags
	// port to listen on (assume it will have the format of ":3000")
	var port string
	fs := flag.NewFlagSet("todos", flag.ExitOnError)
	fs.StringVar(&port, "port", ":3000", "port to listen on")
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: %s [OPTIONS]\n\nOPTIONS:\n", fs.Name())
		fs.PrintDefaults()
	}
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	// ensure port is prefixed with a colon
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	// end flags

	router := chi.NewRouter()
	todoStore := stores.NewTodosStore()
	_, _ = todoStore.Create(ctx, domain.NewTodo("Bake a cake"))
	_, _ = todoStore.Create(ctx, domain.NewTodo("Feed the cat"))
	_, _ = todoStore.Create(ctx, domain.NewTodo("Take out the trash"))

	// home
	router.Get("/", handleHome(todoStore))
	// todos
	router.Route("/todos", func(r chi.Router) {
		r.Get("/", handleSearch(todoStore))
		r.Post("/", handleCreate(todoStore))
		r.Route("/{todoId}", func(r chi.Router) {
			r.Patch("/", handleUpdate(todoStore))
			r.Post("/edit", handleUpdate(todoStore))
			r.Get("/", handleGet(todoStore))
			r.Delete("/", handleDelete(todoStore))
			r.Post("/delete", handleDelete(todoStore))
		})
		r.Post("/sort", handleSort(todoStore))
	})
	// assets
	assets.Mount(router)

	server := &http.Server{
		Addr:    port,
		Handler: http.TimeoutHandler(router, 30*time.Second, "request timed out"),
	}

	// Display the localhost address and port
	fmt.Printf("Listening on http://localhost%s\n", port)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
