import type { DesignFinder } from '$lib/services/Universe';
import {
	AiDifficulty,
	CostSchema,
	GameDBObjectSchema,
	HabSchema,
	PlayerOrdersSchema,
	PlayerPlansSchema,
	PlayerRelation,
	PlayerSchema,
	PlayerSpecSchema,
	PlayerStatsSchema,
	ProductionPlanItemSchema,
	Prt,
	QueueItemType,
	RaceSchema,
	RaceSpecSchema,
	ResearchCostSchema,
	SpendLeftoverPointsOn,
	TechField,
	type BattlePlan,
	type CargoJson as Cargo,
	type CostJson as Cost,
	type MapObjectTarget,
	type Player,
	type PlayerMessage,
	type PlayerRelationship,
	type PlayerScore,
	type ProductionPlan,
	type ProductionPlanItem,
	type ProductionQueueItem,
	type Race,
	type TechLevel,
	type TransportPlan
} from '$lib/types/cs-proto';
import type { CS } from '$lib/wasm';
import { create, merge, type UnknownField } from '@bufbuild/protobuf';
import { add, emptyCargo } from './Cargo';
import { multiply } from './Cost';

import { targetsEqual } from './MapObject';
import type { CommandedPlanet } from './Planet';
import { humanoid } from './Race';
import { type TechLike } from './Tech';
import { emptyTechLevel, hasRequiredLevels } from './TechLevel';
import { string } from './Vector';

export const TechFields: TechField[] = [
	TechField.ENERGY,
	TechField.WEAPONS,
	TechField.PROPULSION,
	TechField.CONSTRUCTION,
	TechField.ELECTRONICS,
	TechField.BIOTECHNOLOGY
];

export class CommandedPlayer implements Player {
	$typeName: 'craig_stars.v1.Player';
	$unknown = undefined;
	gameDbObject = create(GameDBObjectSchema);
	playerOrders = create(PlayerOrdersSchema);
	playerPlans = create(PlayerPlansSchema);

	userId = BigInt(0);
	name = '';
	num = 0;
	ready = false;
	aiControlled = false;
	aiDifficulty = AiDifficulty.NORMAL;
	guest: boolean = false;
	submittedTurn = false;
	color = '#00FF00';
	defaultHullSet: number = 0;
	race = new CommandedPlayerRace(humanoid());
	techLevels: TechLevel = emptyTechLevel();
	techLevelsSpent: TechLevel = emptyTechLevel();
	researchSpentLastYear = 0;
	relations: PlayerRelationship[] = [];
	messages: PlayerMessage[] = [];
	scoreHistory: PlayerScore[] = [];
	acquiredTechs: Record<string, boolean> = {};
	achievedVictoryConditions = 0;
	victor = false;
	archived = false;
	spec = create(PlayerSpecSchema);
	stats = create(PlayerStatsSchema, {});

	constructor(data?: Player) {
		this.$typeName = 'craig_stars.v1.Player';

		if (data) {
			merge(PlayerSchema, this, data);
		}
	}

	isFriend(playerNum: number | undefined): boolean {
		return (
			playerNum != undefined &&
			playerNum > 0 &&
			playerNum <= this.relations.length &&
			this.relations[playerNum - 1].relation === PlayerRelation.FRIEND
		);
	}

	isSharingMap(playerNum: number | undefined): boolean {
		return !!(
			playerNum &&
			playerNum > 0 &&
			playerNum <= this.relations.length &&
			this.relations[playerNum - 1].relation === PlayerRelation.FRIEND &&
			!!this.relations[playerNum - 1].shareMap
		);
	}

	isNeutral(playerNum: number): boolean {
		return (
			playerNum > 0 &&
			playerNum <= this.relations.length &&
			this.relations[playerNum - 1].relation === PlayerRelation.NEUTRAL
		);
	}

	isEnemy(playerNum: number): boolean {
		return (
			playerNum > 0 &&
			playerNum <= this.relations.length &&
			this.relations[playerNum - 1].relation === PlayerRelation.ENEMY
		);
	}

	isFriendOrNeutral(playerNum: number): boolean {
		return (
			playerNum > 0 &&
			playerNum <= this.relations.length &&
			(this.relations[playerNum - 1].relation === PlayerRelation.FRIEND ||
				this.relations[playerNum - 1].relation === PlayerRelation.NEUTRAL)
		);
	}

	getBattlePlan(num: number): BattlePlan | undefined {
		return this.playerPlans.battlePlans.find((p) => p.num === num);
	}

	getProductionPlan(num: number): ProductionPlan | undefined {
		return this.playerPlans.productionPlans.find((p) => p.num === num);
	}

	getTransportPlan(num: number): TransportPlan | undefined {
		return this.playerPlans.transportPlans.find((p) => p.num === num);
	}

	hasTech(tech: TechLike | undefined): boolean {
		if (!tech) {
			return false;
		}
		return (
			canLearnTech(this, tech) &&
			hasRequiredLevels(this.techLevels, tech.tech?.requirements?.techLevel) &&
			(!tech.tech?.requirements?.acquirable || this.hasAcquiredTech(tech))
		);
	}

	hasAcquiredTech(tech: TechLike): boolean {
		if (!tech.tech?.requirements?.acquirable) {
			return true;
		}
		return !!this.acquiredTechs[tech.tech?.name ?? ''];
	}

	getAllies(): number[] {
		const allies: number[] = [];
		this.relations.forEach((r, index) => {
			if (r.relation === PlayerRelation.FRIEND) {
				allies.push(index + 1);
			}
		});
		return allies;
	}

	public async getItemCost(
		cs: CS,
		item: ProductionQueueItem | undefined,
		designFinder: DesignFinder,
		planet?: CommandedPlanet,
		quantity = 1
	): Promise<Cost> {
		if (item) {
			switch (item.type) {
				case QueueItemType.STARBASE:
					if (item.designNum) {
						const design = designFinder.getMyDesign(item.designNum);
						if (planet?.spec.planetStarbaseSpec?.hasStarbase) {
							const starbaseToUpgrade = designFinder.getMyDesign(
								planet.spec.planetStarbaseSpec?.starbaseDesignNum
							);
							if (starbaseToUpgrade && design) {
								const { cost } = await cs.wasmService.getStarbaseUpgradeCost({
									design: starbaseToUpgrade,
									newDesign: design
								});
								return cost ?? create(CostSchema, {});
							}
						}
						return multiply(design?.spec?.cost ?? {}, quantity);
					}
					break;
				case QueueItemType.SHIP_TOKEN:
					if (item.designNum) {
						const design = designFinder.getMyDesign(item.designNum);
						return multiply(design?.spec?.cost ?? {}, quantity);
					}
					break;
				default:
					if (this.race?.spec?.costs) {
						return multiply(this.race.spec.costs[item.type] ?? {}, quantity);
					}
			}
		}
		return create(CostSchema, {});
	}

	public getByHandTransfer(target: MapObjectTarget): Cargo {
		const key = string(target.targetPosition);
		const transfers = this.playerOrders.cargoTransfers[key]?.transfers ?? [];
		let cargo = emptyCargo();
		if (!transfers) {
			return cargo;
		}

		// sum up all transfers for this target
		transfers
			.filter((t) => targetsEqual(target, t.mapObjectTarget))
			.forEach((t) => (cargo = add(cargo, t.cargo)));

		return cargo;
	}
}

export class CommandedPlayerRace implements Race {
	$typeName: 'craig_stars.v1.Race';
	$unknown?: UnknownField[] | undefined;

	id = BigInt(0);
	createdAt = undefined;
	updatedAt = undefined;
	userId = BigInt(0);
	name = '';
	pluralName = '';
	spendLeftoverPointsOn = SpendLeftoverPointsOn.UNSPECIFIED;
	prt = Prt.UNSPECIFIED;
	lrts = 0;
	habLow = create(HabSchema);
	habHigh = create(HabSchema);
	growthRate = 0;
	popEfficiency = 0;
	factoryOutput = 0;
	factoryCost = 0;
	numFactories = 0;
	factoriesCostLess = false;
	immuneGrav = false;
	immuneTemp = false;
	immuneRad = false;
	mineOutput = 0;
	mineCost = 0;
	numMines = 0;
	researchCost = create(ResearchCostSchema);
	techsStartHigh = false;
	spec = create(RaceSpecSchema);

	constructor(data?: Race) {
		this.$typeName = 'craig_stars.v1.Race';

		if (data) {
			merge(RaceSchema, this, data);
		}
	}
}

export function canLearnTech(player: CommandedPlayer, tech: TechLike): boolean {
	if (!tech.tech?.requirements) {
		return true;
	}
	const requirements = tech.tech.requirements;
	if (
		requirements.prtsRequired?.length &&
		requirements.prtsRequired.indexOf(player.race.prt) == -1
	) {
		return false;
	}
	if (requirements.prtsDenied?.length && requirements.prtsDenied.indexOf(player.race.prt) != -1) {
		return false;
	}

	if (requirements.lrtsRequired && (player.race.lrts & requirements.lrtsRequired) == 0) {
		return false;
	}
	if (requirements.lrtsDenied && (player.race.lrts & requirements.lrtsDenied) > 0) {
		return false;
	}

	if (!player.hasAcquiredTech(tech)) {
		return false;
	}
	return true;
}

export function planItemFromQueueItemType(type: QueueItemType): ProductionPlanItem {
	return create(ProductionPlanItemSchema, {
		type,
		quantity: 1
	});
}
