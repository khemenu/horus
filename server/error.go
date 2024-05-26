package server

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrPermissionDenied = status.Error(codes.PermissionDenied, codes.PermissionDenied.String())
)

func newErrMissingRequiredField(path string) error {
	return status.Errorf(codes.InvalidArgument, "required: %s", path)
}
