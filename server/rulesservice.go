package server

import (
	"context"

	"connectrpc.com/connect"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/proto/converter"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1/craig_starsv1connect"
)

// NewRulesServiceHandler constructs the RulesService handler.
func NewRulesServiceHandler() craig_starsv1connect.RulesServiceHandler {
	return &rulesService{}
}

type rulesService struct{}

// GetRules returns a single Rules by name, or NOT_FOUND if missing.
func (s *rulesService) GetRules(ctx context.Context, req *connect.Request[craig_starsv1.GetRulesRequest]) (*connect.Response[craig_starsv1.GetRulesResponse], error) {
	rules := cs.StandardRules
	return connect.NewResponse(&craig_starsv1.GetRulesResponse{
		Rules: converter.C.ConvertCSRules(&rules),
	}), nil
}
