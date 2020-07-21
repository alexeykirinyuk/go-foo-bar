package bar

import (
	"fmt"
	"github.com/alexeykirinyuk/tech-task-go/components/auth"
	"github.com/alexeykirinyuk/tech-task-go/data"
	"github.com/alexeykirinyuk/tech-task-go/libs"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/volatiletech/authboss/v3"
	"net/http"
	"time"
)

const barTemplateBasePath = "components/bar/templates"

func ConfigureRouter(mux *chi.Mux, boss *authboss.Authboss, dbProvider data.IDatabaseProvider) {
	mux.Group(func(r chi.Router) {
		libs.ConfigureAuthMiddleware(r, boss, auth.RoleMember, auth.RoleAdmin)
		configure(r, dbProvider)
	})
}

func configure(r chi.Router, dbProvider data.IDatabaseProvider) {
	r.MethodFunc("GET", "/templates", func(w http.ResponseWriter, r *http.Request) {
		barStorage := newStorage(dbProvider)

		items, err := barStorage.getAll()
		if err != nil {
			panic(err)
		}

		libs.Render(w, r, barTemplateBasePath+"view.tpl", items)
	})
	r.MethodFunc("GET", "/templates/create", func(w http.ResponseWriter, r *http.Request) {
		libs.Render(w, r, "views/templates/create.tpl", nil)
	})
	r.MethodFunc("POST", "/templates/create", func(w http.ResponseWriter, r *http.Request) {
		service := newService(dbProvider)

		bar := extractBarFromFormData(r, uuid.New())

		errs, ok := service.create(bar)
		if !ok {
			libs.BadRequest(w, r, libs.ToResponse(errs))
			return
		}

		redirectToAllBar(w, r)
	})
	r.MethodFunc("POST", "/templates/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		service := newService(dbProvider)

		id, ok := extractIdFromRouteParameters(w, r)
		if !ok {
			return
		}

		errs, ok := service.delete(id)
		if !ok {
			libs.BadRequest(w, r, libs.ToResponse(errs))
			return
		}

		redirectToAllBar(w, r)
	})
	r.MethodFunc("GET", "/templates/update/{id}", func(w http.ResponseWriter, r *http.Request) {
		service := newService(dbProvider)

		id, ok := extractIdFromRouteParameters(w, r)
		if !ok {
			return
		}

		foo, errs, ok := service.get(id)
		if !ok {
			libs.BadRequest(w, r, libs.ToResponse(errs))
			return
		}

		libs.Render(w, r, barTemplateBasePath+"update.tpl", foo)
	})
	r.MethodFunc("POST", "/templates/update/{id}", func(w http.ResponseWriter, r *http.Request) {
		service := newService(dbProvider)

		id, ok := extractIdFromRouteParameters(w, r)
		if !ok {
			return
		}

		bar := extractBarFromFormData(r, id)
		bar.Id = id

		errs, ok := service.update(bar)
		if !ok {
			libs.BadRequest(w, r, libs.ToResponse(errs))
			return
		}

		redirectToAllBar(w, r)
	})
}

func extractBarFromFormData(r *http.Request, id uuid.UUID) bar {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	openingDate, err := time.Parse("test-layout", r.FormValue("opening_date"))
	if err != nil {
		panic(err)
	}

	return bar{
		Id:          id,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Address:     r.FormValue("address"),
		OpeningDate: openingDate,
	}
}

func extractIdFromRouteParameters(w http.ResponseWriter, r *http.Request) (id uuid.UUID, ok bool) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		libs.BadRequest(w, r, fmt.Sprintf("can't parse id '%s'", idStr))
		return uuid.UUID{}, false
	}

	return id, true
}

func redirectToAllBar(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/templates", http.StatusMovedPermanently)
}
