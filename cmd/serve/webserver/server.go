package webserver

import (
	"DockerPostgreExample/internal/data"
	"github.com/fasthttp/router"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"net"
)

const (
	AcceptJson  = "application/json"
	ContentJson = "application/json; charset=utf-8"
	ContentRest = "application/vnd.pgrst.object+json; charset=utf-8"
)

type webServer struct {
	config      WebConfig
	addr        string
	log         zerolog.Logger
	ln          net.Listener
	router      *router.Router
	debug       bool
	dataManager *data.Manager
}

func NewServer(cfg WebConfig, log zerolog.Logger, dataManager *data.Manager) *webServer {
	s := &webServer{
		config:      cfg,
		addr:        ServerAddr,
		log:         log,
		router:      router.New(),
		debug:       true,
		dataManager: dataManager,
	}
	return s
}

func (ws *webServer) Run() (err error) {
	ws.muxRouter()
	// reuse port

	ws.ln, err = reuseport.Listen("tcp4", ws.addr)
	if err != nil {
		return err
	}
	s := &fasthttp.Server{
		Handler:            ws.router.Handler,
		Name:               ws.config.Name,
		ReadBufferSize:     ws.config.ReadBufferSize,
		MaxConnsPerIP:      ws.config.MaxConnsPerIP,
		MaxRequestsPerConn: ws.config.MaxRequestsPerConn,
		MaxRequestBodySize: ws.config.MaxRequestBodySize, //  100 << 20, // 100MB // 1024 * 4, // MaxRequestBodySize:
		Concurrency:        ws.config.Concurrency,
	}

	return s.Serve(ws.ln)
}
