package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"PackageServer/db"
	"PackageServer/fs"
	"PackageServer/logger"
	fileServer "PackageServer/server/file"
	httpServer "PackageServer/server/http"
)

func Start() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	logger.InitLog()
	fs.InitFs()
	db.InitCache()
	db.InitDb()
	srv := httpServer.InitHttpServer()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Errorf("%v\n", err)
			quit <- syscall.SIGTERM
		}
	}()

	fileCtx, fileCancel := context.WithCancel(context.Background())
	defer fileCancel()
	fileSrv := fileServer.NewFileServer(fileCtx)
	go func() {
		if err := fileSrv.Serve(); err != nil {
			logger.Log.Errorf("%v\n", err)
			quit <- syscall.SIGTERM
		}
	}()

	<-quit
	logger.Log.Panicln("shutdown server gracefully... ")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error("servcer shutdown...")
	}

	select {
	case <-ctx.Done():
		logger.Log.Infoln("server shuwdown timeout of 5 sec")
	}

}
