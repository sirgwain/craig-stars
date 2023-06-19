package dbsqlx

import "github.com/sirgwain/craig-stars/game"

// Get the rules for a game
func (c *client) GetRulesForGame(gameID int64) (*game.Rules, error) {

	// TODO: implement rules saving to DB
	rules := game.NewRules()
	return &rules, nil
}
