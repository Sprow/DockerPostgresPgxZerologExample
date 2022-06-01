package handler

import (
	"DockerPostgreExample/internal/data"
	"DockerPostgreExample/internal/logger"
	"github.com/go-chi/chi/v5"
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
		return
	}
	if err = dataObj.IsValid(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error().Err(err).Msg("")
		return
	}
	err = h.dataManager.AddDataObj(r.Context(), dataObj)
	if err != nil {
		logger.Log.Error().Stack().Err(err).Msg("")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

type objID struct {
	ID int `json:"id"`
}

func (h *Handler) removeData(w http.ResponseWriter, r *http.Request) {
	decoder := jsoniter.NewDecoder(r.Body)
	var id objID
	err := decoder.Decode(&id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error().Stack().Err(err).Msg("can't decode")
		return
	}
	err = h.dataManager.RemoveDataObj(r.Context(), id.ID)
	if err != nil {
		logger.Log.Error().Stack().Err(err).Msg("can's remove data")
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (h *Handler) updateData(w http.ResponseWriter, r *http.Request) {
	decoder := jsoniter.NewDecoder(r.Body)
	var dataObj data.Obj
	err := decoder.Decode(&dataObj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error().Stack().Err(err).Msg("can't decode")
		return
	}

	if err = dataObj.IsValid(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error().Err(err).Msg("not valid data")
		return
	}

	err = h.dataManager.UpdateDataObj(r.Context(), dataObj)
	if err != nil {
		logger.Log.Error().Stack().Err(err).Msg("can't update")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
