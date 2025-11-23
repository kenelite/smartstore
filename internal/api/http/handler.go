package apihttp

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/kenelite/smartstore/internal/storage/smart"
)

type Handler struct {
	svc *smart.Service
}

func NewHandler(svc *smart.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Put("/v1/{env}/{region}/{bucket}/*", h.PutObject)
	r.Get("/v1/{env}/{region}/{bucket}/*", h.GetObject)
}

func (h *Handler) PutObject(w http.ResponseWriter, r *http.Request) {
	env := chi.URLParam(r, "env")
	region := chi.URLParam(r, "region")
	bucket := chi.URLParam(r, "bucket")
	key := chi.URLParam(r, "*")

	ct := r.Header.Get("Content-Type")
	sizeStr := r.Header.Get("Content-Length")
	size, _ := strconv.ParseInt(sizeStr, 10, 64)

	req := &smart.PutRequest{
		Env:           env,
		LogicalRegion: region,
		Bucket:        bucket,
		Key:           key,
		ContentType:   ct,
		Size:          size,
		Body:          r.Body,
		StorageClass:  r.Header.Get("X-Storage-Class"),
	}

	resp, err := h.svc.Put(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetObject(w http.ResponseWriter, r *http.Request) {
	env := chi.URLParam(r, "env")
	region := chi.URLParam(r, "region")
	bucket := chi.URLParam(r, "bucket")
	key := chi.URLParam(r, "*")

	resp, err := h.svc.Get(r.Context(), &smart.GetRequest{
		Env:           env,
		LogicalRegion: region,
		Bucket:        bucket,
		Key:           key,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.ContentType != "" {
		w.Header().Set("Content-Type", resp.ContentType)
	}
	if resp.Size > 0 {
		w.Header().Set("Content-Length", strconv.FormatInt(resp.Size, 10))
	}
	if _, err := io.Copy(w, resp.Body); err != nil {
		return
	}
}
