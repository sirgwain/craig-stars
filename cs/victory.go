package cs

import (
	"math"
	"sort"
)

type VictoryConditions struct {
	Conditions               Bitmask `json:"conditions"`
	NumCriteriaRequired      int     `json:"numCriteriaRequired"`
	YearsPassed              int     `json:"yearsPassed"`
	OwnPlanets               int     `json:"ownPlanets"`
	AttainTechLevel          int     `json:"attainTechLevel"`
	AttainTechLevelNumFields int     `json:"attainTechLevelNumFields"`
	ExceedsScore             int     `json:"exceedsScore"`
	ExceedsSecondPlaceScore  int     `json:"exceedsSecondPlaceScore"`
	ProductionCapacity       int     `json:"productionCapacity"`
	OwnCapitalShips          int     `json:"ownCapitalShips"`
	HighestScoreAfterYears   int     `json:"highestScoreAfterYears"`
}

type VictoryCondition Bitmask

const (
	VictoryConditionNone                        = 0
	VictoryConditionOwnPlanets VictoryCondition = 1 << (iota - 1)
	VictoryConditionAttainTechLevels
	VictoryConditionExceedsScore
	VictoryConditionExceedsSecondPlaceScore
	VictoryConditionProductionCapacity
	VictoryConditionOwnCapitalShips
	VictoryConditionHighestScoreAfterYears
)

// Evaluate achieved conditions and qualifying winners for the current scores.
type victoryChecker struct {
	game *FullGame
}

func newVictoryChecker(game *FullGame) victoryChecker {
	return victoryChecker{game}
}

func (v *victoryChecker) checkForVictor(player *Player) error {
	if len(player.ScoreHistory) == 0 {
		return nil
	}
	// reset player victory conditions to none
	player.AchievedVictoryConditions = VictoryConditionNone

	score := player.ScoreHistory[len(player.ScoreHistory)-1]
	if v.game.YearsPassed() < v.game.Rules.ShowPublicScoresAfterYears {
		// Clear premature saved flags without evaluating any victory conditions.
		score.AchievedVictoryConditions = VictoryConditionNone
		player.ScoreHistory[len(player.ScoreHistory)-1] = score
		return nil
	}

	v.checkOwnPlanets(player, score)
	v.checkAttainTechLevels(player)
	v.checkExceedScore(player, score)
	v.checkProductionCapacity(player, score)
	v.checkOwnCapitalShips(player, score)
	if len(v.game.Players) > 1 {
		v.checkExceedSecondPlaceScore(player, score)
		v.checkHighestScore(player, score)
	}

	// update the history with this player's AchievedVictoryConditions
	// this way we know over time when victories were achieved
	// and when a victor is declared, other players will know when
	// victory was achieved.
	score.AchievedVictoryConditions = player.AchievedVictoryConditions
	player.ScoreHistory[len(player.ScoreHistory)-1] = score

	// a game with no victory conditions (or no required criteria) never has a victor
	if v.game.VictoryConditions.Conditions == 0 || v.game.VictoryConditions.NumCriteriaRequired <= 0 {
		return nil
	}

	required := min(v.game.VictoryConditions.NumCriteriaRequired, v.game.VictoryConditions.Conditions.countBits())
	achieved := player.AchievedVictoryConditions & v.game.VictoryConditions.Conditions
	alive := score.Planets+score.UnarmedShips+score.EscortShips+score.CapitalShips > 0
	if !v.game.VictorDeclared && len(v.game.Players) > 1 && alive && required > 0 && achieved.countBits() >= required && v.game.YearsPassed() >= v.game.VictoryConditions.YearsPassed {
		// The caller declares all qualifying players together after checking everyone.
		player.Victor = true
	}

	return nil
}

func (v *victoryChecker) checkOwnPlanets(player *Player, score PlayerScore) {
	// i.e. if we own more than 60% of the planets, we have this victory condition
	if float64(score.Planets) >= math.Round(float64(len(v.game.Planets))*float64(v.game.VictoryConditions.OwnPlanets)/100) {
		player.AchievedVictoryConditions |= Bitmask(VictoryConditionOwnPlanets)
	}
}

func (v *victoryChecker) checkAttainTechLevels(player *Player) {
	numAttained := 0
	for _, field := range TechFields {
		if player.TechLevels.Get(field) >= v.game.VictoryConditions.AttainTechLevel {
			numAttained++
		}
	}
	if numAttained >= v.game.VictoryConditions.AttainTechLevelNumFields {
		player.AchievedVictoryConditions |= Bitmask(VictoryConditionAttainTechLevels)
	}
}

func (v *victoryChecker) checkExceedScore(player *Player, score PlayerScore) {
	if score.Score >= v.game.VictoryConditions.ExceedsScore {
		player.AchievedVictoryConditions |= Bitmask(VictoryConditionExceedsScore)
	}
}

func (v *victoryChecker) checkExceedSecondPlaceScore(player *Player, score PlayerScore) {
	if len(v.game.Players) > 1 {
		scores := make([]int, len(v.game.Players))
		for i, player := range v.game.Players {
			scores[i] = player.GetScore().Score
		}
		sort.Slice(scores, func(i, j int) bool {
			return scores[i] > scores[j]
		})

		threshold := int(float64(scores[1]) * (100 + float64(v.game.VictoryConditions.ExceedsSecondPlaceScore)) / 100)
		if score.Score == scores[0] && score.Score >= threshold {
			player.AchievedVictoryConditions |= Bitmask(VictoryConditionExceedsSecondPlaceScore)
		}
	}
}

func (v *victoryChecker) checkProductionCapacity(player *Player, score PlayerScore) {
	if score.Resources >= v.game.VictoryConditions.ProductionCapacity*1000 {
		player.AchievedVictoryConditions |= Bitmask(VictoryConditionProductionCapacity)
	}
}

func (v *victoryChecker) checkOwnCapitalShips(player *Player, score PlayerScore) {
	if score.CapitalShips >= v.game.VictoryConditions.OwnCapitalShips {
		player.AchievedVictoryConditions |= Bitmask(VictoryConditionOwnCapitalShips)
	}
}

func (v *victoryChecker) checkHighestScore(player *Player, score PlayerScore) {
	if v.game.YearsPassed() < v.game.VictoryConditions.HighestScoreAfterYears {
		return
	}
	for _, other := range v.game.Players {
		if other.Num != player.Num && other.GetScore().Score >= score.Score {
			return
		}
	}
	player.AchievedVictoryConditions |= Bitmask(VictoryConditionHighestScoreAfterYears)
}
