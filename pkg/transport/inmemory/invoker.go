package inmemory

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// InMemoryInvoker allows invoking RPC procedures in-process without opening a network port.
type InMemoryInvoker interface {
	// Invoke calls an RPC procedure using Protobuf messages in-memory.
	Invoke(ctx context.Context, procedure string, req proto.Message, resp proto.Message) error

	// InvokeJSON calls an RPC procedure using raw JSON in-memory.
	InvokeJSON(ctx context.Context, procedure string, reqJSON []byte) ([]byte, error)
}

type directInvoker struct {
	handler http.Handler
}

// NewInvoker creates a new in-memory invoker wrapping the given HTTP/ConnectRPC handler.
func NewInvoker(handler http.Handler) InMemoryInvoker {
	return &directInvoker{handler: handler}
}

func (i *directInvoker) formatPath(procedure string) string {
	if !strings.HasPrefix(procedure, "/") {
		return "/" + procedure
	}
	return procedure
}

func (i *directInvoker) Invoke(ctx context.Context, procedure string, req proto.Message, resp proto.Message) error {
	path := i.formatPath(procedure)

	reqBytes, err := proto.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal proto request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(reqBytes))
	if err != nil {
		return fmt.Errorf("failed to create in-memory request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/proto")
	httpReq.Header.Set("Connect-Protocol-Version", "1")

	rec := httptest.NewRecorder()
	i.handler.ServeHTTP(rec, httpReq)

	if rec.Code >= 400 {
		return fmt.Errorf("rpc procedure %q returned status %d: %s", procedure, rec.Code, rec.Body.String())
	}

	if resp != nil {
		if err := proto.Unmarshal(rec.Body.Bytes(), resp); err != nil {
			return fmt.Errorf("failed to unmarshal proto response: %w", err)
		}
	}

	return nil
}

func (i *directInvoker) InvokeJSON(ctx context.Context, procedure string, reqJSON []byte) ([]byte, error) {
	path := i.formatPath(procedure)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(reqJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create in-memory json request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Connect-Protocol-Version", "1")

	rec := httptest.NewRecorder()
	i.handler.ServeHTTP(rec, httpReq)

	if rec.Code >= 400 {
		return nil, fmt.Errorf("rpc procedure %q returned status %d: %s", procedure, rec.Code, rec.Body.String())
	}

	return rec.Body.Bytes(), nil
}

// ProtoToJSON is a utility to format a proto.Message into JSON.
func ProtoToJSON(msg proto.Message) ([]byte, error) {
	return protojson.Marshal(msg)
}

// JSONToProto is a utility to parse JSON into a proto.Message.
func JSONToProto(jsonData []byte, msg proto.Message) error {
	return protojson.Unmarshal(jsonData, msg)
}
