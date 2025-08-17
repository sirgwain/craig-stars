package grpcwasm

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"
)

// Handler is a generic WASM-style interface for calling RPCs using raw bytes.
type Handler interface {
	Call(ctx context.Context, method string, req []byte) ([]byte, error)
}

// HandlerFunc is an adapter to allow the use of ordinary functions as Handlers.
type HandlerFunc func(ctx context.Context, method string, req []byte) ([]byte, error)

func (f HandlerFunc) Call(ctx context.Context, method string, req []byte) ([]byte, error) {
	return f(ctx, method, req)
}

func HandleUnary[Req proto.Message, Res proto.Message](
	ctx context.Context,
	reqBytes []byte,
	req Req,
	handler func(context.Context, Req) (Res, error),
) ([]byte, error) {
	if err := proto.Unmarshal(reqBytes, req); err != nil {
		return nil, fmt.Errorf("unmarshal request: %w", err)
	}

	resp, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}
	respBytes, err := proto.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response: %w", err)
	}
	return respBytes, nil
}
