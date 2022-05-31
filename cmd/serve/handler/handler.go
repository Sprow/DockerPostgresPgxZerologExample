package handler

import (
	"DockerPostgreExample/internal/data"
	"DockerPostgreExample/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

type Handler struct {
	dataManager *data.Manager
}

func NewHandler(dataManager *data.Manager) *Handler {
	return &Handler{
		dataManager: dataManager,
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
		logger.Log.Error().Err(err).Msg("")
	}
	encoder := jsoniter.NewEncoder(w)
	err = encoder.Encode(d)
	if err != nil {
		logger.Log.Error().Err(err).Msg("")
	}
}

func (h *Handler) addData(w http.ResponseWriter, r *http.Request) {
	decoder := jsoniter.NewDecoder(r.Body)
	var dataObj data.Obj
	err := decoder.Decode(&dataObj)
	if err != nil { // bad request
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error().Stack().Err(err).Msg("")
	}
	err = h.dataManager.AddDataObj(r.Context(), dataObj)
	if err != nil {
		logger.Log.Error().Stack().Err(err).Msg("")
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
		logger.Log.Error().Stack().Err(err).Msg("")
	}
	err = h.dataManager.RemoveDataObj(r.Context(), id.ID)
	if err != nil {
		logger.Log.Error().Stack().Err(err).Msg("")
	}
}

func (h *Handler) updateData(w http.ResponseWriter, r *http.Request) {
	decoder := jsoniter.NewDecoder(r.Body)
	var dataObj data.Obj
	err := decoder.Decode(&dataObj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error().Stack().Err(err).Msg("")
	}
	err = h.dataManager.UpdateDataObj(r.Context(), dataObj)
	if err != nil {
		logger.Log.Error().Stack().Err(err).Msg("")
	}
}
