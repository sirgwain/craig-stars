package ai

import "github.com/sirgwain/craig-stars/cs"

// update player plans
func (ai *aiPlayer) plan() {
	// update the default planet production plan
	items := []cs.ProductionPlanItem{
		{Type: cs.QueueItemTypeAutoMinTerraform, Quantity: 5},
	}

	if !ai.Race.Spec.InnateResources {
		items = append(items, cs.ProductionPlanItem{Type: cs.QueueItemTypeAutoFactories, Quantity: 10})
	}

	if !ai.Race.Spec.InnateMining {
		items = append(items, cs.ProductionPlanItem{Type: cs.QueueItemTypeAutoMines, Quantity: 10})
	}
	items = append(items, cs.ProductionPlanItem{Type: cs.QueueItemTypeAutoMaxTerraform, Quantity: 5})

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
	ai.PlayerPlans.ProductionPlans[0].ContributesOnlyLeftoverToResearch = true

	// update the player orders
	ai.client.UpdatePlayerOrders(ai.Player, ai.Planets, ai.PlayerOrders, &ai.game.Rules)
}

