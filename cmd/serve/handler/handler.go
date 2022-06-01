package handler

import (
	"DockerPostgreExample/internal/data"
	"DockerPostgreExample/internal/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
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

func (h *Handler) getAllData(ctx *fasthttp.RequestCtx) {
	d, err := h.dataManager.GetAllData(ctx)
	if err != nil {
		logger.Log.Error().Stack().Err(err).Msg("Can't get all data")
		ctx.Error("Can't get all data", fasthttp.StatusInternalServerError)
		return
	}

	jsonData, err := jsoniter.Marshal(d)
	if err != nil {
		return
	}
	_, err = ctx.Write(jsonData)
	if err != nil {
		return
	}
	ctx.SetContentType("application/json")
}

func (h *Handler) addData(ctx *fasthttp.RequestCtx) {
	var dataObj data.Obj
	err := jsoniter.Unmarshal(ctx.Request.Body(), &dataObj)
	if err != nil { // bad request
		logger.Log.Error().Stack().Err(err).Msg("can't unmarshal req.body")
		ctx.Error("can't unmarshal req.body", fasthttp.StatusBadRequest)
		return
	}
	if err = dataObj.IsValid(); err != nil {
		logger.Log.Error().Err(err).Msg("invalid data")
		//example log with stack
		//logger.Log.Error().Stack().Err(logger.Stack(err)).Msg("invalid data")
		ctx.Error("invalid data", fasthttp.StatusBadRequest)
		return
	}
	err = h.dataManager.AddDataObj(ctx, dataObj)
	if err != nil {
		logger.Log.Error().Stack().Err(err).Msg("")
		ctx.Error("update failed", fasthttp.StatusBadRequest)
		return
	}
	ctx.SetStatusCode(http.StatusCreated)
	return
}

type objID struct {
	ID int `json:"id"`
}

func (h *Handler) removeData(ctx *fasthttp.RequestCtx) {
	var id objID
	err := jsoniter.Unmarshal(ctx.Request.Body(), &id)
	if err != nil { // bad request
		logger.Log.Error().Stack().Err(err).Msg("can't unmarshal req.body")
		ctx.Error("can't unmarshal req.body", fasthttp.StatusBadRequest)
		return
	}
	err = h.dataManager.RemoveDataObj(ctx, id.ID)
	if err != nil {
		logger.Log.Error().Stack().Err(err).Msg("")
		ctx.Error("delete failed", fasthttp.StatusNotFound)
		return
	}
	logger.Log.Info().Msgf("remove Object id=%d", id.ID)
	ctx.SetStatusCode(fasthttp.StatusOK)
	return
}

func (h *Handler) updateData(ctx *fasthttp.RequestCtx) {
	var dataObj data.Obj
	err := jsoniter.Unmarshal(ctx.Request.Body(), &dataObj)
	if err != nil {
		logger.Log.Error().Stack().Err(err).Msg("can't unmarshal")
		ctx.Error("delete failed", fasthttp.StatusBadRequest)
		return
	}

	if err = dataObj.IsValid(); err != nil {
		logger.Log.Error().Stack().Err(err).Msg("not valid data")
		ctx.Error("not valid data", fasthttp.StatusBadRequest)
		return
	}

	err = h.dataManager.UpdateDataObj(ctx, dataObj)
	if err != nil {
		logger.Log.Error().Stack().Err(err).Msg("can't update")
		ctx.Error("update failed", fasthttp.StatusBadRequest)
		return
	}
	ctx.SetStatusCode(http.StatusOK)
	return
}
