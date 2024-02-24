package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stackus/errors"
	"github.com/stackus/hxgo"

	"github.com/stackus/goht-todos/domain"
	"github.com/stackus/goht-todos/templates/pages"
	"github.com/stackus/goht-todos/templates/partials"
)

func handleHome(todos domain.TodosStore) http.HandlerFunc {
	return errorHandler(func(w http.ResponseWriter, r *http.Request) error {
		todos, err := todos.List(r.Context())
		if err != nil {
			return err
		}

		return render(r.Context(), pages.HomePage(todos), w)
	})
}

func handleCreate(todos domain.TodosStore) http.HandlerFunc {
	return errorHandler(func(w http.ResponseWriter, r *http.Request) error {
		_, form, err := parseRequest(r)
		if err != nil {
			return err
		}
		description := form.Get("description")

		todo, err := todos.Create(r.Context(), domain.NewTodo(description))
		if err != nil {
			return err
		}

		if hx.IsHtmx(r) {
			return render(r.Context(), partials.RenderTodo(todo), w)
		}

		http.Redirect(w, r, "/", http.StatusFound)

		return nil
	})
}

func handleUpdate(todos domain.TodosStore) http.HandlerFunc {
	return errorHandler(func(w http.ResponseWriter, r *http.Request) error {
		_, form, err := parseRequest(r)
		if err != nil {
			return err
		}
		id := chi.URLParam(r, "todoId")
		completed := form.Get("completed") == "true"
		description := form.Get("description")

		todo, err := todos.Update(r.Context(), id, domain.Todo{
			Description: description,
			Completed:   completed,
		})
		if err != nil {
			return err
		}

		if hx.IsHtmx(r) {
			return render(r.Context(), partials.RenderTodo(todo), w)
		}

		http.Redirect(w, r, "/", http.StatusFound)

		return nil
	})
}

func handleGet(todos domain.TodosStore) http.HandlerFunc {
	return errorHandler(func(w http.ResponseWriter, r *http.Request) error {
		id := chi.URLParam(r, "todoId")

		todo, err := todos.Get(r.Context(), id)
		if err != nil {
			return err
		}

		if hx.IsHtmx(r) {
			return render(r.Context(), partials.EditTodoForm(todo), w)
		}

		return render(r.Context(), pages.TodoPage(todo), w)
	})
}

func handleDelete(todos domain.TodosStore) http.HandlerFunc {
	return errorHandler(func(w http.ResponseWriter, r *http.Request) error {
		id := chi.URLParam(r, "todoId")

		if err := todos.Delete(r.Context(), id); err != nil {
			return err
		}

		if hx.IsHtmx(r) {
			w.WriteHeader(http.StatusOK)
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
		}

		return nil
	})
}

func handleSort(todos domain.TodosStore) http.HandlerFunc {
	return errorHandler(func(w http.ResponseWriter, r *http.Request) error {
		_, form, err := parseRequest(r)
		if err != nil {
			return err
		}
		if form["id"] == nil {
			return errors.ErrBadRequest.Msg("missing ids")
		}

		if err := todos.Reorder(r.Context(), form["id"]); err != nil {
			return err
		}

		if hx.IsHtmx(r) {
			w.WriteHeader(http.StatusNoContent)
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
		}

		return nil
	})
}

func handleSearch(todos domain.TodosStore) http.HandlerFunc {
	return errorHandler(func(w http.ResponseWriter, r *http.Request) error {
		search := r.URL.Query().Get("search")
		todos, err := todos.Filter(r.Context(), search)
		if err != nil {
			return err
		}

		if hx.IsHtmx(r) {
			return render(r.Context(), partials.RenderTodos(todos), w)
		}

		return render(r.Context(), pages.TodosPage(todos, search), w)
	})
}
