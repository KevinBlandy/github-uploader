package server

import (
	"context"
	"fmt"
	"github-uploader/config"
	"net/http"
)

var httpServer *http.Server

func Run (handler http.Handler) error {
	var serverConfig = config.App.Server
	httpServer = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port),
		Handler:           handler,
	}
	if serverConfig.ReadTimeout > 0 {
		httpServer.ReadTimeout = serverConfig.ReadTimeout
	}
	if serverConfig.ReadHeaderTimeout > 0 {
		httpServer.ReadHeaderTimeout = serverConfig.ReadHeaderTimeout
	}
	if serverConfig.WriteTimeout > 0 {
		httpServer.WriteTimeout = serverConfig.WriteTimeout
	}
	if serverConfig.IdleTimeout > 0 {
		httpServer.IdleTimeout = serverConfig.IdleTimeout
	}
	if serverConfig.MaxHeaderSize != nil && serverConfig.MaxHeaderSize.ToByte() > 0{
		httpServer.MaxHeaderBytes = int(serverConfig.MaxHeaderSize.ToByte())
	}

	if serverConfig.Ssl != nil && serverConfig.Ssl.Enabled {
		return httpServer.ListenAndServeTLS(serverConfig.Ssl.CertFile, serverConfig.Ssl.KeyFile)
	}
	return httpServer.ListenAndServe()
}

func Shutdown (ctx context.Context) error {
	return httpServer.Shutdown(ctx)
}