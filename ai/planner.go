package ai

import "github.com/sirgwain/craig-stars/cs"

// update player plans
func (ai *aiPlayer) plan() error {
	// update the default planet production plan
	items := []cs.ProductionPlanItem{}

	if !ai.Race.Spec.InnateResources {
		items = append(items, cs.ProductionPlanItem{Type: cs.QueueItemTypeAutoFactories, Quantity: 250})
	}

	if !ai.Race.Spec.InnateMining {
		items = append(items, cs.ProductionPlanItem{Type: cs.QueueItemTypeAutoMines, Quantity: 250})
	}

	if ai.doctrine.defensesPerYear > 0 && !ai.Race.Spec.LivesOnStarbases {
		items = append(items, cs.ProductionPlanItem{Type: cs.QueueItemTypeAutoDefenses, Quantity: ai.doctrine.defensesPerYear})
	}

	if !ai.Race.Spec.Instaforming && !(ai.Race.ImmuneGrav && ai.Race.ImmuneTemp && ai.Race.ImmuneRad) {
		items = append(items, cs.ProductionPlanItem{Type: cs.QueueItemTypeAutoMaxTerraform, Quantity: 10})
	}

	if !ai.Race.Spec.InnateResources {
		items = append(items, cs.ProductionPlanItem{Type: cs.QueueItemTypeAutoFactories, Quantity: 500})
	}

	if !ai.Race.Spec.InnateMining {
		items = append(items, cs.ProductionPlanItem{Type: cs.QueueItemTypeAutoMines, Quantity: 500})
	}

	if !ai.Race.Spec.LivesOnStarbases {
		items = append(items, cs.ProductionPlanItem{Type: cs.QueueItemTypeAutoDefenses, Quantity: 100})
	}

	ai.PlayerPlans.ProductionPlans[0].Items = items
	ai.PlayerPlans.ProductionPlans[0].ContributesOnlyLeftoverToResearch = ai.config.researchReserve == 0
	orders := ai.PlayerOrders
	orders.ResearchAmount = ai.config.researchReserve
	ai.client.UpdatePlayerOrders(ai.Player, ai.Planets, orders, &ai.game.Rules)
	for _, planet := range ai.Planets {
		leftoversOnly := ai.config.researchReserve == 0 || planet.Spec.PopulationDensity < .25
		if planet.ContributesOnlyLeftoverToResearch != leftoversOnly {
			planet.ContributesOnlyLeftoverToResearch = leftoversOnly
			if err := ai.client.UpdatePlanetOrders(&ai.game.Rules, ai.Player, planet, planet.PlanetOrders); err != nil {
				return err
			}
		}
	}
	return nil
}
