package db

import "github.com/sirgwain/craig-stars/cs"

// Get the rules for a game
func (c *txClient) GetRulesForGame(gameID int64) (*cs.Rules, error) {

	// TODO: implement rules saving to DB
	rules := cs.NewRules()
	return &rules, nil
}
