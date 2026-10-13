package ai

// Experiment lets a test switch one behavior off for one player, so the AI can
// be compared with and without it in the same game. It returns true if the
// behavior is switched off for that player. It is nil outside those tests.
var Experiment func(playerNum int, behavior string) bool

// enabled returns true unless a test has switched this behavior off for this player
func (ai *aiPlayer) enabled(behavior string) bool {
	return Experiment == nil || !Experiment(ai.Num, behavior)
}
