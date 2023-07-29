package cs

import (
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

type victory struct {
	game *FullGame
}

type victoryChecker interface {
	checkForVictor(player *Player) error
}

func newVictoryChecker(game *FullGame) victoryChecker {
	return &victory{game}
}

func (v *victory) checkForVictor(player *Player) error {
	if len(player.ScoreHistory) == 0 {
		return nil
	}
	score := player.ScoreHistory[len(player.ScoreHistory)-1]

	if v.game.VictoryConditions.Conditions&Bitmask(VictoryConditionOwnPlanets) > 0 {
		v.checkOwnPlanets(player, score)
	}
	if v.game.VictoryConditions.Conditions&Bitmask(VictoryConditionAttainTechLevels) > 0 {
		v.checkAttainTechLevels(player, score)
	}
	if v.game.VictoryConditions.Conditions&Bitmask(VictoryConditionExceedsScore) > 0 {
		v.checkExceedScore(player, score)
	}
	if v.game.VictoryConditions.Conditions&Bitmask(VictoryConditionExceedsSecondPlaceScore) > 0 {
		v.checkExceedSecondPlaceScore(player, score)
	}
	if v.game.VictoryConditions.Conditions&Bitmask(VictoryConditionProductionCapacity) > 0 {
		v.checkProductionCapacity(player, score)
	}
	if v.game.VictoryConditions.Conditions&Bitmask(VictoryConditionOwnCapitalShips) > 0 {
		v.checkOwnCapitalShips(player, score)
	}
	if v.game.VictoryConditions.Conditions&Bitmask(VictoryConditionHighestScoreAfterYears) > 0 {
		v.checkHighestScore(player, score)
	}

	if player.AchievedVictoryConditions.countBits() >= v.game.VictoryConditions.NumCriteriaRequired && v.game.YearsPassed() >= v.game.VictoryConditions.YearsPassed {
		// we have a victor!
		player.Victor = true
		v.game.VictorDeclared = true
	}

	return nil
}

func (v *victory) checkOwnPlanets(player *Player, score PlayerScore) {
	// i.e. if we own more than 60% of the planets, we have this victory condition
	if float64(score.Planets) >= float64(len(v.game.Planets))*(float64(v.game.VictoryConditions.OwnPlanets)/100) {
		player.AchievedVictoryConditions |= Bitmask(VictoryConditionOwnPlanets)
	}
}

func (v *victory) checkAttainTechLevels(player *Player, score PlayerScore) {
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

func (v *victory) checkExceedScore(player *Player, score PlayerScore) {
	if score.Score > v.game.VictoryConditions.ExceedsScore {
		player.AchievedVictoryConditions |= Bitmask(VictoryConditionExceedsScore)
	}
}

func (v *victory) checkExceedSecondPlaceScore(player *Player, score PlayerScore) {
	if len(v.game.Players) > 1 {
		scores := make([]int, len(v.game.Players))
		for i := range v.game.Players {
			scores[i] = score.Score
		}
		sort.Slice(scores, func(i, j int) bool {
			return scores[i] > scores[j]
		})
		if float64(score.Score)/float64(scores[1]-1)*100 > float64(v.game.VictoryConditions.ExceedsSecondPlaceScore) {
			player.AchievedVictoryConditions |= Bitmask(VictoryConditionExceedsSecondPlaceScore)
		}
	}
}

func (v *victory) checkProductionCapacity(player *Player, score PlayerScore) {
	productionCapacity := 0
	for _, planet := range v.game.Planets {
		if planet.PlayerNum == player.Num {
			productionCapacity += planet.Spec.ResourcesPerYear
		}
	}
	if productionCapacity >= v.game.VictoryConditions.ProductionCapacity * 1000 {
		player.AchievedVictoryConditions |= Bitmask(VictoryConditionProductionCapacity)
	}
}

func (v *victory) checkOwnCapitalShips(player *Player, score PlayerScore) {
	if score.CapitalShips >= v.game.VictoryConditions.OwnCapitalShips {
		player.AchievedVictoryConditions |= Bitmask(VictoryConditionOwnCapitalShips)
	}
}

func (v *victory) checkHighestScore(player *Player, score PlayerScore) {
	if v.game.YearsPassed() >= v.game.VictoryConditions.HighestScoreAfterYears {
		sortedScores := make([]int, len(v.game.Players))
		for i := range v.game.Players {
			sortedScores[i] = score.Score
		}
		sort.Slice(sortedScores, func(i, j int) bool {
			return sortedScores[i] > sortedScores[j]
		})
		if score.Score == sortedScores[0] {
			player.AchievedVictoryConditions |= Bitmask(VictoryConditionHighestScoreAfterYears)
		}
	}
}
