package server

import (
	"goquickstart/config"
	"goquickstart/database"
	"goquickstart/logger"
	"goquickstart/server/handler"
	"net/http"
	"time"
)

type Server struct {
	dataBase database.DataBase
	lgr      logger.Logger
	config   *config.SvrConfig
}

func (svr Server) Database() database.DataBase {
	return svr.dataBase
}

func NewServer() *Server {

	svr := Server{}
	svr.config = config.GetConfig()
	svr.lgr = logger.GetLogger(svr.config.LogFileDir)
	svr.dataBase.Init(svr.config, svr.lgr)
	return &svr
}

func (svr *Server) Run() {
	go svr.lgr.Run()
	myServer := http.Server{
		Addr:         svr.config.SvrHostPort,
		Handler:      handler.NewHandler(svr.lgr, svr.dataBase),
		WriteTimeout: time.Duration(svr.config.WriteTimeout),
		ReadTimeout:  time.Duration(svr.config.ReadTimeout),
	}
	svr.lgr.Write("server ready for listening")
	myServer.ListenAndServe()
}
