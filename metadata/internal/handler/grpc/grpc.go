package grpc

import (
	"movieexample.com/gen"
	"movieexample.com/metadata/internal/controller"
)

type Handler struct {
	gen.UnimplementedMetadataServiceServer
	svc *controller.MetadataService
}

func New(ctrl *metada)
