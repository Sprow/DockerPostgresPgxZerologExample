package handler

import (
	"DockerPostgreExample/internal/data"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
	"net/http"
)

type Handler struct {
	dataManager *data.Manager
	log         zerolog.Logger
}

func NewHandler(dataManager *data.Manager, log zerolog.Logger) *Handler {
	return &Handler{
		dataManager: dataManager,
		log:         log,
	}
}

func (h *Handler) Register(r *chi.Mux) {
	r.Get("/", h.getAllData)
	r.Post("/add_data", h.addData)
	r.Post("/remove_data", h.removeData)
	r.Post("/update_data", h.updateData)
}

func (h *Handler) getAllData(w http.ResponseWriter, r *http.Request) {
	d, err := h.dataManager.GetAllData(r.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("")
	}
	encoder := jsoniter.NewEncoder(w)
	err = encoder.Encode(d)
	if err != nil {
		h.log.Error().Err(err).Msg("")
	}
}

func (h *Handler) addData(w http.ResponseWriter, r *http.Request) {
	decoder := jsoniter.NewDecoder(r.Body)
	var dataObj data.Obj
	err := decoder.Decode(&dataObj)
	if err != nil { // bad request
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error().Stack().Err(err).Msg("")
	}
	err = h.dataManager.AddDataObj(r.Context(), dataObj)
	if err != nil {
		h.log.Error().Stack().Err(err).Msg("")
	}
}

type objID struct {
	ID uuid.UUID `json:"id"`
}

func (h *Handler) removeData(w http.ResponseWriter, r *http.Request) {
	decoder := jsoniter.NewDecoder(r.Body)
	var id objID
	err := decoder.Decode(&id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error().Stack().Err(err).Msg("")
	}
	err = h.dataManager.RemoveDataObj(r.Context(), id.ID)
	if err != nil {
		h.log.Error().Stack().Err(err).Msg("")
	}
}

func (h *Handler) updateData(w http.ResponseWriter, r *http.Request) {
	decoder := jsoniter.NewDecoder(r.Body)
	var dataObj data.Obj
	err := decoder.Decode(&dataObj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error().Stack().Err(err).Msg("")
	}
	err = h.dataManager.UpdateDataObj(r.Context(), dataObj)
	if err != nil {
		h.log.Error().Stack().Err(err).Msg("")
	}
}
