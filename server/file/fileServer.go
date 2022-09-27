package fileServer

import (
	"PackageServer/logger"
	"PackageServer/service"
	"PackageServer/variable"
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
)

type FileServer struct {
	ctx             context.Context
	FileTransferSvc service.FileTransferSvc
}

func NewFileServer(ctx context.Context) *FileServer {
	return &FileServer{
		ctx:             ctx,
		FileTransferSvc: service.NewFileTransferSvc(),
	}
}

func (f *FileServer) Serve() error {
	for {
		select {
		case <-f.ctx.Done():
			return errors.New("FileTransferServer quit by ctx")
		case job := <-variable.FilePushJobChan:
			logger.Log.Infof("receive a file job: %v", job.FileName)
			v := validator.New()
			if err := v.Struct(&job); err != nil {
				logger.Log.Infof("invalid jobChan Err: %v", err)
			} else {
				f.FileTransferSvc.StartPackagePushJob(job)
			}
		}
	}
}
