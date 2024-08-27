import { fromHabType } from '$lib/services/Terraformer';
import type { CostFinder, DesignFinder } from '$lib/services/Universe';
import type { ProductionQueueItem } from '$lib/types/Production';
import type { BattleAttackWho, BattleRecord, BattleTactic, BattleTarget } from './Battle';
import { multiply, type Cost, minus, minZero } from './Cost';
import type { Fleet, WaypointTransportTasks } from './Fleet';
import { HabTypes, type Hab } from './Hab';
import type { Message } from './Message';
import type { MineField } from './MineField';
import type { MineralPacket } from './MineralPacket';
import type { MysteryTrader } from './MysteryTrader';
import type { Planet } from './Planet';
import { QueueItemTypes } from './QueueItemType';
import { humanoid, type Race } from './Race';
import type { Salvage } from './Salvage';
import type { ShipDesign } from './ShipDesign';
import {
	TechCategory,
	TerraformHabTypes,
	getBestTerraform,
	type Tech,
	type TechDefense,
	type TechPlanetaryScanner,
	type TechStore
} from './Tech';
import {
	TechField,
	emptyTechLevel,
	hasRequiredLevels,
	minTechLevel,
	type TechLevel
} from './TechLevel';
import type { Wormhole } from './Wormhole';

export type PlayerStatus = {
	updatedAt: string;

	userId?: number;
	num: number;
	color: string;
	name: string;
	race: Race;
	ready?: boolean;
	aiControlled?: boolean;
	guest?: boolean;
	submittedTurn?: boolean;
};

export type PlayerResponse = {
	id?: number;
	createdAt?: string;
	updatedAt?: string;

	gameId: number;
	userId?: number;
	num: number;
	color: string;
	name?: string;
	race: Race;
	ready?: boolean;
	aiControlled?: boolean;
	guest?: boolean;
	submittedTurn?: boolean;
	techLevels: TechLevel;
	techLevelsSpent: TechLevel;
	designs?: ShipDesign[];
	researchSpentLastYear?: number;
	achievedVictoryConditions?: number;
	relations: PlayerRelationship[];
	acquiredTechs?: Record<string, boolean>;
	spec: PlayerSpec;
} & PlayerOrders &
	PlayerMessages &
	PlayerPlans;

export type PlayerMessages = {
	messages: Message[];
};

export type PlayerPlans = {
	battlePlans?: BattlePlan[];
	productionPlans?: ProductionPlan[];
	transportPlans?: TransportPlan[];
};

export type PlayerIntels = {
	players: PlayerIntel[];
	scores: PlayerScore[][];
	planets: Planet[];
	fleets?: Fleet[];
	mineFields?: MineField[];
	mineralPackets?: MineralPacket[];
	salvages?: Salvage[];
	wormholes?: Wormhole[];
	mysteryTraders?: MysteryTrader[];
	battles?: BattleRecord[];
};

export type PlayerUniverse = {
	designs: ShipDesign[];
	planets: Planet[];
	fleets: Fleet[];
	starbases: Fleet[];
	mineFields: MineField[];
	mineralPackets: MineralPacket[];
	salvages: Salvage[];
};

export type PlayerOrders = {
	researching: TechField;
	nextResearchField: NextResearchField;
	researchAmount: number;
};

export type BattlePlan = {
	num: number;
	name: string;
	primaryTarget: BattleTarget;
	secondaryTarget: BattleTarget;
	tactic: BattleTactic;
	attackWho: BattleAttackWho;
	dumpCargo: boolean;
};

export type TransportPlan = {
	num: number;
	name: string;
	tasks: WaypointTransportTasks;
};

export type ProductionPlan = {
	num: number;
	name: string;
	items: ProductionQueueItem[];
	contributesOnlyLeftoverToResearch?: boolean;
};

export type PlayerSpec = {
	planetaryScanner?: TechPlanetaryScanner;
	defense?: TechDefense;
	resourcesPerYear?: number;
	resourcesPerYearResearch?: number;
	resourcesPerYearResearchEstimated?: number;
	currentResearchCost?: number;
};

export type PlayerIntel = {
	name: string;
	num: number;
	color: string;
	seen?: boolean;
	raceName?: string;
	racePluralName?: string;
};

export type PlayerScore = {
	planets: number;
	starbases: number;
	unarmedShips: number;
	escortShips: number;
	capitalShips: number;
	techLevels: number;
	resources: number;
	score: number;
	rank: number;
	achievedVictoryConditions?: number;
};

export enum NextResearchField {
	SameField = 'SameField',
	Energy = 'Energy',
	Weapons = 'Weapons',
	Propulsion = 'Propulsion',
	Construction = 'Construction',
	Electronics = 'Electronics',
	Biotechnology = 'Biotechnology',
	LowestField = 'LowestField'
}

export type PlayerRelationship = {
	relation?: PlayerRelation;
	shareMap?: boolean;
};

export enum PlayerRelation {
	Neutral = 'Neutral',
	Friend = 'Friend',
	Enemy = 'Enemy'
}

export class Player implements PlayerResponse, CostFinder {
	id = 0;
	createdAt?: string | undefined;
	updatedAt?: string | undefined;

	gameId = 0;
	num = 0;

	userId?: number | undefined;
	name = '';
	color = '#00FF00';
	race = { ...humanoid() };
	ready = false;
	aiControlled = false;
	submittedTurn = false;
	techLevels: TechLevel = { ...emptyTechLevel() };
	techLevelsSpent: TechLevel = { ...emptyTechLevel() };
	researchSpentLastYear = 0;
	researching: TechField = TechField.Energy;
	nextResearchField: NextResearchField = NextResearchField.Energy;
	researchAmount = 0;
	battlePlans: BattlePlan[] = [];
	productionPlans: ProductionPlan[] = [];
	transportPlans: TransportPlan[] = [];
	messages: Message[] = [];
	relations: PlayerRelationship[] = [];
	acquiredTechs: Record<string, boolean> = {};
	designs?: ShipDesign[] = [];
	spec: PlayerSpec = {};

	constructor(data?: PlayerResponse) {
		if (data) {
			Object.assign(this, data);
		}
	}

	isFriend(playerNum: number): boolean {
		return (
			playerNum > 0 &&
			playerNum <= this.relations.length &&
			this.relations[playerNum - 1].relation === PlayerRelation.Friend
		);
	}

	isSharingMap(playerNum: number): boolean {
		return (
			playerNum > 0 &&
			playerNum <= this.relations.length &&
			this.relations[playerNum - 1].relation === PlayerRelation.Friend &&
			!!this.relations[playerNum - 1].shareMap
		);
	}

	isNeutral(playerNum: number): boolean {
		return (
			playerNum > 0 &&
			playerNum <= this.relations.length &&
			this.relations[playerNum - 1].relation === PlayerRelation.Neutral
		);
	}

	isEnemy(playerNum: number): boolean {
		return (
			playerNum > 0 &&
			playerNum <= this.relations.length &&
			this.relations[playerNum - 1].relation === PlayerRelation.Enemy
		);
	}

	isFriendOrNeutral(playerNum: number): boolean {
		return (
			playerNum > 0 &&
			playerNum <= this.relations.length &&
			(this.relations[playerNum - 1].relation === PlayerRelation.Friend ||
				this.relations[playerNum - 1].relation === PlayerRelation.Neutral)
		);
	}

	getBattlePlan(num: number): BattlePlan | undefined {
		return this.battlePlans.find((p) => p.num === num);
	}

	getProductionPlan(num: number): ProductionPlan | undefined {
		return this.productionPlans.find((p) => p.num === num);
	}

	getTransportPlan(num: number): TransportPlan | undefined {
		return this.transportPlans.find((p) => p.num === num);
	}

	hasTech(tech: Tech): boolean {
		return (
			canLearnTech(this, tech) &&
			hasRequiredLevels(this.techLevels, tech.requirements) &&
			(!tech.requirements.acquirable || this.hasAcquiredTech(tech))
		);
	}

	hasAcquiredTech(tech: Tech): boolean {
		if (!tech.requirements.acquirable) {
			return true;
		}
		return !!this.acquiredTechs[tech.name];
	}

	getAllies(): number[] {
		const allies: number[] = [];
		this.relations.forEach((r, index) => {
			if (r.relation === PlayerRelation.Friend) {
				allies.push(index + 1);
			}
		});
		return allies;
	}

	public getItemCost(
		item: ProductionQueueItem | undefined,
		designFinder: DesignFinder,
		techStore: TechStore,
		planet?: Planet,
		quantity = 1
	): Cost {
		if (item) {
			switch (item.type) {
				case QueueItemTypes.Starbase: // TODO: starbase upgrades...
					if (item.designNum) {
						const design = designFinder.getMyDesign(item.designNum);
						if (planet?.spec.hasStarbase) {
							const starbaseToUpgrade = designFinder.getMyDesign(planet.spec.starbaseDesignNum);
							if (starbaseToUpgrade && design) {
								return multiply(
									this.getStarbaseUpgradeCost(techStore, starbaseToUpgrade, design),
									quantity
								);
							}
						}
						return multiply(design?.spec.cost ?? {}, quantity);
					}
					break;
				case QueueItemTypes.ShipToken:
					if (item.designNum) {
						const design = designFinder.getMyDesign(item.designNum);
						return multiply(design?.spec.cost ?? {}, quantity);
					}
					break;
				default:
					return multiply(this.race?.spec?.costs[item.type] ?? {}, quantity);
			}
		}
		return {};
	}

	// get the cost of upgrading this starbase
	public getStarbaseUpgradeCost(
		techStore: TechStore,
		design: ShipDesign,
		updatedDesign: ShipDesign
	): Cost {
		// TODO: update this if we update the server side
		return minZero(minus(updatedDesign.spec?.cost ?? {}, design.spec?.cost ?? {}));
	}

	// get a player's ability to terraform
	public getTerraformAbility(techStore: TechStore): Hab {
		const terraformAbility: Hab = { grav: 0, temp: 0, rad: 0 };
		const bestTT = getBestTerraform(techStore, this, TerraformHabTypes.All);
		if (bestTT) {
			terraformAbility.grav = bestTT.ability;
			terraformAbility.temp = bestTT.ability;
			terraformAbility.rad = bestTT.ability;
		}

		Object.values(HabTypes).forEach((habType) => {
			const bestTerraform = getBestTerraform(techStore, this, fromHabType(habType));
			if (bestTerraform) {
				terraformAbility.grav = Math.max(bestTerraform.ability, terraformAbility.grav ?? 0);
				terraformAbility.temp = Math.max(bestTerraform.ability, terraformAbility.temp ?? 0);
				terraformAbility.rad = Math.max(bestTerraform.ability, terraformAbility.rad ?? 0);
			}
		});
		return terraformAbility;
	}

	public getTechCost(t: Tech): Cost {
		// figure out miniaturization
		// this is 4% per level above the required tech we have.
		// We count the smallest diff, i.e. if you have
		// tech level 10 energy, 12 bio and the tech costs 9 energy, 4 bio
		// the smallest level difference you have is 1 energy level (not 8 bio levels)

		const levelDiff: TechLevel = {
			energy: -1,
			weapons: -1,
			propulsion: -1,
			construction: -1,
			electronics: -1,
			biotechnology: -1
		};

		const techLevels = {
			energy: this.techLevels.energy ?? 0,
			weapons: this.techLevels.weapons ?? 0,
			propulsion: this.techLevels.propulsion ?? 0,
			construction: this.techLevels.construction ?? 0,
			electronics: this.techLevels.electronics ?? 0,
			biotechnology: this.techLevels.biotechnology ?? 0
		};
		const requirements = {
			energy: t.requirements.energy ?? 0,
			weapons: t.requirements.weapons ?? 0,
			propulsion: t.requirements.propulsion ?? 0,
			construction: t.requirements.construction ?? 0,
			electronics: t.requirements.electronics ?? 0,
			biotechnology: t.requirements.biotechnology ?? 0
		};
		// From the diff between the player level and the requirements, find the lowest difference
		// i.e. 1 energey level in the example above
		let numTechLevelsAboveRequired = Number.MAX_SAFE_INTEGER;
		if (requirements.energy > 0) {
			levelDiff.energy = techLevels.energy - requirements.energy;
			numTechLevelsAboveRequired = Math.min(levelDiff.energy, numTechLevelsAboveRequired);
		}
		if (requirements.weapons > 0) {
			levelDiff.weapons = techLevels.weapons - requirements.weapons;
			numTechLevelsAboveRequired = Math.min(levelDiff.weapons, numTechLevelsAboveRequired);
		}
		if (requirements.propulsion > 0) {
			levelDiff.propulsion = techLevels.propulsion - requirements.propulsion;
			numTechLevelsAboveRequired = Math.min(levelDiff.propulsion, numTechLevelsAboveRequired);
		}
		if (requirements.construction > 0) {
			levelDiff.construction = techLevels.construction - requirements.construction;
			numTechLevelsAboveRequired = Math.min(levelDiff.construction, numTechLevelsAboveRequired);
		}
		if (requirements.electronics > 0) {
			levelDiff.electronics = techLevels.electronics - requirements.electronics;
			numTechLevelsAboveRequired = Math.min(levelDiff.electronics, numTechLevelsAboveRequired);
		}
		if (requirements.biotechnology > 0) {
			levelDiff.biotechnology = techLevels.biotechnology - requirements.biotechnology;
			numTechLevelsAboveRequired = Math.min(levelDiff.biotechnology, numTechLevelsAboveRequired);
		}

		// for starter techs, they are all 0 requirements, so just use our lowest field
		if (numTechLevelsAboveRequired == Number.MAX_SAFE_INTEGER) {
			numTechLevelsAboveRequired = minTechLevel(techLevels);
		}

		// As we learn techs, they get cheaper. We start off with full priced techs, but every additional level of research we learn makes
		// techs cost a little less, maxing out at some discount (i.e. 75% or 80% for races with BET)

		let miniaturization = Math.min(
			this.race.spec?.miniaturizationMax ?? 0,
			(this.race.spec?.miniaturizationPerLevel ?? 0) * numTechLevelsAboveRequired
		);
		// New techs cost BET races 2x
		// new techs will have 0 for miniaturization.
		let miniaturizationFactor = this.race.spec?.newTechCostFactor ?? 1;
		if (numTechLevelsAboveRequired > 0) {
			miniaturizationFactor = 1 - miniaturization;
		}

		let cost = Object.assign({}, t.cost);
		const costOffset = {
			engine: this.race.spec?.techCostOffset.engine ?? 0,
			beamWeapon: this.race.spec?.techCostOffset.beamWeapon ?? 0,
			torpedo: this.race.spec?.techCostOffset.torpedo ?? 0,
			bomb: this.race.spec?.techCostOffset.bomb ?? 0,
			planetaryDefense: this.race.spec?.techCostOffset.planetaryDefense ?? 0
		};
		switch (t.category) {
			case TechCategory.Engine:
				cost = multiply(cost, 1 + costOffset.engine);
				break;
			case TechCategory.BeamWeapon:
				cost = multiply(cost, 1 + costOffset.beamWeapon);
				break;
			case TechCategory.Bomb:
				cost = multiply(cost, 1 + costOffset.bomb);
				break;
			case TechCategory.Torpedo:
				cost = multiply(cost, 1 + costOffset.torpedo);
				break;
		}
		return {
			ironium: Math.ceil((cost.ironium ?? 0) * miniaturizationFactor),
			boranium: Math.ceil((cost.boranium ?? 0) * miniaturizationFactor),
			germanium: Math.ceil((cost.germanium ?? 0) * miniaturizationFactor),
			resources: Math.ceil((cost.resources ?? 0) * miniaturizationFactor)
		};
	}
}

export function canLearnTech(player: Player, tech: Tech): boolean {
	const requirements = tech.requirements;
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
