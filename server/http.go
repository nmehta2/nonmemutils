package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pressly/chi"
)

func (ms *ModelStore) handleGetModelID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	model, ok := ctx.Value("model").(*Model)

	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	w.Write([]byte(fmt.Sprintf("modelID:%v", model.ID)))
}

// ModelCtx is the context
func (ms *ModelStore) ModelCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		modelID := chi.URLParam(r, "modelID")
		mid, _ := strconv.ParseInt(modelID, 10, 64)
		model, err := ms.GetModel(int(mid))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "model", model)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}