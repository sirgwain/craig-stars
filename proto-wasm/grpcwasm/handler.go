package grpcwasm

import (
	"context"
)

type VTProtoMessage interface {
	MarshalVT() (dAtA []byte, err error)
	UnmarshalVT(dAtA []byte) error
}

// Handler is a generic WASM-style interface for calling RPCs using raw bytes.
type Handler interface {
	Call(ctx context.Context, method string, req []byte) ([]byte, error)
}

// HandlerFunc is an adapter to allow the use of ordinary functions as Handlers.
type HandlerFunc func(ctx context.Context, method string, req []byte) ([]byte, error)

func (f HandlerFunc) Call(ctx context.Context, method string, req []byte) ([]byte, error) {
	return f(ctx, method, req)
}

func HandleUnary[Req, Res VTProtoMessage](
	ctx context.Context,
	reqBytes []byte,
	req Req,
	handler func(context.Context, Req) (Res, error),
) ([]byte, error) {

	if err := req.UnmarshalVT(reqBytes); err != nil {
		return nil, err
	}

	resp, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}
	respBytes, err := resp.MarshalVT()
	if err != nil {
		return nil, err
	}
	return respBytes, nil
}
