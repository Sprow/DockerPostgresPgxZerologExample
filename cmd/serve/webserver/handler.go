package webserver

import (
	"DockerPostgreExample/internal/data"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"net/http"
)

func (ws *webServer) getAllData() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		ws.log.Debug().Str("request", ctx.String()).Msg("")
		allData, err := ws.dataManager.GetAllData(ctx)
		if err != nil {
			ws.log.Error().Err(err).Msg("")
		}

		ctx.SetContentType(ContentJson)
		ctx.SetStatusCode(http.StatusOK)

		jsonData, err := jsoniter.Marshal(allData)
		if err != nil {
			ws.log.Error().Err(err).Msg("")
		}
		ctx.SetBody(jsonData)

		ws.log.Debug().Int("len", len(allData)).Msg("get all data")
		return
	}
}

func (ws *webServer) addData() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		ws.log.Debug().Str("request", ctx.String()).Msg("")
		var dataObj data.Obj
		if err := jsoniter.Unmarshal(ctx.Request.Body(), &dataObj); err != nil {
			ws.log.Error().Stack().Err(err).Msg("can't unmarshal req.body")
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		ctx.SetContentType(ContentJson)
		err := ws.dataManager.AddDataObj(ctx, dataObj)
		if err != nil {
			ctx.SetStatusCode(http.StatusConflict)
			//ctx.SetBody([]byte(`{"Done": "false"}`))
			ws.log.Error().Stack().Err(err).Msg("")
			return
		}

		ctx.SetStatusCode(http.StatusCreated)
		//ctx.SetBody([]byte(`{"Done": "true"}`))
		return
	}
}

type objID struct {
	ID uuid.UUID `json:"id"`
}

func (ws *webServer) removeData() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		ws.log.Debug().Str("request", ctx.String()).Msg("")
		var id objID

		if err := jsoniter.Unmarshal(ctx.Request.Body(), &id); err != nil {
			ws.log.Error().Stack().Err(err).Msg("can't unmarshal req.body")
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		ctx.SetContentType(ContentJson)
		err := ws.dataManager.RemoveDataObj(ctx, id.ID)
		if err != nil {
			ctx.SetStatusCode(http.StatusNotFound)
			ws.log.Error().Stack().Err(err).Msg("")
			return
		}

		ctx.SetStatusCode(http.StatusOK)
		return
	}
}

func (ws *webServer) updateData() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		ws.log.Debug().Str("request", ctx.String()).Msg("")
		var dataObj data.Obj

		if err := jsoniter.Unmarshal(ctx.Request.Body(), &dataObj); err != nil {
			ws.log.Error().Stack().Err(err).Msg("can't unmarshal update request body")
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		ctx.SetContentType(ContentJson)
		err := ws.dataManager.UpdateDataObj(ctx, dataObj)
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			ws.log.Error().Stack().Err(err).Msg("")
			return
		}

		ctx.SetStatusCode(http.StatusOK)
		return
	}
}
