package server

import (
	"context"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	mcpserver "github.com/mark3labs/mcp-go/server"
	mcpmark3labs "github.com/redpanda-data/protoc-gen-go-mcp/pkg/runtime/mark3labs"
	"github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1/craig_starsv1connect"
	"github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1/craig_starsv1mcp"
	"github.com/sirgwain/craig-stars/version"
)

// newMCPHandler creates the hosted Streamable HTTP MCP server and registers
// generated tools that forward to the existing Connect API.
func (s *server) newMCPHandler() http.Handler {
	raw, runtimeServer := mcpmark3labs.NewServer("craig-stars", version.Semver, mcpserver.WithToolCapabilities(true))
	opts := []connect.ClientOption{connect.WithInterceptors(newBearerForwardInterceptor())}
	baseURL := s.internalConnectBaseURL()
	client := http.DefaultClient

	craig_starsv1mcp.ForwardToConnectUserServiceClient(runtimeServer, craig_starsv1connect.NewUserServiceClient(client, baseURL, opts...))
	craig_starsv1mcp.ForwardToConnectRaceServiceClient(runtimeServer, craig_starsv1connect.NewRaceServiceClient(client, baseURL, opts...))
	craig_starsv1mcp.ForwardToConnectGameServiceClient(runtimeServer, craig_starsv1connect.NewGameServiceClient(client, baseURL, opts...))
	craig_starsv1mcp.ForwardToConnectPlayerServiceClient(runtimeServer, craig_starsv1connect.NewPlayerServiceClient(client, baseURL, opts...))
	craig_starsv1mcp.ForwardToConnectBattlePlanServiceClient(runtimeServer, craig_starsv1connect.NewBattlePlanServiceClient(client, baseURL, opts...))
	craig_starsv1mcp.ForwardToConnectProductionPlanServiceClient(runtimeServer, craig_starsv1connect.NewProductionPlanServiceClient(client, baseURL, opts...))
	craig_starsv1mcp.ForwardToConnectTransportPlanServiceClient(runtimeServer, craig_starsv1connect.NewTransportPlanServiceClient(client, baseURL, opts...))
	craig_starsv1mcp.ForwardToConnectShipDesignServiceClient(runtimeServer, craig_starsv1connect.NewShipDesignServiceClient(client, baseURL, opts...))
	craig_starsv1mcp.ForwardToConnectPlanetServiceClient(runtimeServer, craig_starsv1connect.NewPlanetServiceClient(client, baseURL, opts...))
	craig_starsv1mcp.ForwardToConnectFleetServiceClient(runtimeServer, craig_starsv1connect.NewFleetServiceClient(client, baseURL, opts...))
	craig_starsv1mcp.ForwardToConnectMinefieldServiceClient(runtimeServer, craig_starsv1connect.NewMinefieldServiceClient(client, baseURL, opts...))

	return mcpserver.NewStreamableHTTPServer(raw,
		mcpserver.WithEndpointPath("/api/mcp"),
		mcpserver.WithStateLess(true),
		mcpserver.WithHTTPContextFunc(func(ctx context.Context, r *http.Request) context.Context {
			if raw := bearerToken(r); raw != "" {
				return context.WithValue(ctx, apiBearerContextKey{}, raw)
			}
			return ctx
		}),
	)
}

// newBearerForwardInterceptor forwards the caller's API bearer token from the
// MCP request context to the internal Connect request.
func newBearerForwardInterceptor() connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if raw, ok := ctx.Value(apiBearerContextKey{}).(string); ok && raw != "" {
				req.Header().Set("Authorization", "Bearer "+raw)
			}
			return next(ctx, req)
		}
	})
}

// internalConnectBaseURL returns the base URL used by MCP tools to call this
// server's Connect API.
func (s *server) internalConnectBaseURL() string {
	addr := s.config.Address
	if strings.HasPrefix(addr, ":") {
		addr = "localhost" + addr
	}
	if !strings.Contains(addr, "://") {
		addr = "http://" + addr
	}
	return strings.TrimRight(addr, "/") + "/api/grpc"
}
