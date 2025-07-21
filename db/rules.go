package db

import (
	"context"

	"github.com/sirgwain/craig-stars/cs"
)

// Get the rules for a game
func (c *client) GetRulesForGame(ctx context.Context, gameID int64) (*cs.Rules, error) {

	// TODO: implement rules saving to DB
	rules := cs.NewRules()
	return &rules, nil
}
