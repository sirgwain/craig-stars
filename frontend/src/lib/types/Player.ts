import { fromHabType } from '$lib/services/Terraformer';
import type { CostFinder, DesignFinder } from '$lib/services/Universe';
import type { CS } from '$lib/wasm';
import { add, emptyCargo } from './Cargo';
import { multiply } from './Cost';
import type {
	BattlePlan,
	Cargo,
	CargoTransfers,
	MapObjectTarget,
	NextResearchField,
	Player,
	PlayerMessage,
	PlayerRelationship,
	PlayerScore,
	PlayerSpec,
	ProductionPlan,
	ProductionQueueItem,
	TechField,
	TransportPlan
} from './cs';
import {
	Biotechnology,
	Construction,
	Electronics,
	Energy,
	NextResearchFieldBiotechnology,
	NextResearchFieldConstruction,
	NextResearchFieldElectronics,
	NextResearchFieldEnergy,
	NextResearchFieldLowestField,
	NextResearchFieldPropulsion,
	NextResearchFieldSameField,
	NextResearchFieldWeapons,
	PlayerRelationEnemy,
	PlayerRelationFriend,
	PlayerRelationNeutral,
	Propulsion,
	QueueItemTypeShipToken,
	QueueItemTypeStarbase,
	TerraformHabTypeAll,
	Weapons,
	type Cost,
	type Hab,
	type Tech,
	type TechLevel,
	type TechStore
} from './cs';
import { HabTypes } from './Hab';
import { targetsEqual } from './MapObject';
import type { CommandedPlanet } from './Planet';
import { humanoid } from './Race';
import { getBestTerraform } from './Tech';
import { emptyTechLevel, hasRequiredLevels } from './TechLevel';
import { string } from './Vector';

export const TechFields: TechField[] = [
	Energy,
	Weapons,
	Propulsion,
	Construction,
	Electronics,
	Biotechnology
];

export const NextResearchFields: NextResearchField[] = [
	NextResearchFieldSameField,
	NextResearchFieldEnergy,
	NextResearchFieldWeapons,
	NextResearchFieldPropulsion,
	NextResearchFieldConstruction,
	NextResearchFieldElectronics,
	NextResearchFieldBiotechnology,
	NextResearchFieldLowestField
];

export class CommandedPlayer implements Player, CostFinder {
	id = 0;
	gameId = 0;
	createdAt = '';
	updatedAt = '';
	num = 0;

	userId?: number | undefined;
	name = '';
	color = '#00FF00';
	race = { ...humanoid() };
	ready = false;
	aiControlled = false;
	submittedTurn = false;
	victor = false;
	archived = false;
	techLevels: TechLevel = { ...emptyTechLevel() };
	techLevelsSpent: TechLevel = { ...emptyTechLevel() };
	researchSpentLastYear = 0;
	researching: TechField = Energy;
	nextResearchField: NextResearchField = NextResearchFieldEnergy;
	researchAmount = 0;
	cargoTransfers: CargoTransfers = {};
	battlePlans: BattlePlan[] = [];
	productionPlans: ProductionPlan[] = [];
	transportPlans: TransportPlan[] = [];
	messages: PlayerMessage[] = [];
	relations: PlayerRelationship[] = [];
	scoreHistory: PlayerScore[] = [];
	acquiredTechs: Record<string, boolean> = {};
	spec: PlayerSpec = {} as PlayerSpec;

	constructor(data?: Player) {
		if (data) {
			Object.assign(this, data);
		}
	}

	isFriend(playerNum: number | undefined): boolean {
		return (
			playerNum != undefined &&
			playerNum > 0 &&
			playerNum <= this.relations.length &&
			this.relations[playerNum - 1].relation === PlayerRelationFriend
		);
	}

	isSharingMap(playerNum: number): boolean {
		return (
			playerNum > 0 &&
			playerNum <= this.relations.length &&
			this.relations[playerNum - 1].relation === PlayerRelationFriend &&
			!!this.relations[playerNum - 1].shareMap
		);
	}

	isNeutral(playerNum: number): boolean {
		return (
			playerNum > 0 &&
			playerNum <= this.relations.length &&
			this.relations[playerNum - 1].relation === PlayerRelationNeutral
		);
	}

	isEnemy(playerNum: number): boolean {
		return (
			playerNum > 0 &&
			playerNum <= this.relations.length &&
			this.relations[playerNum - 1].relation === PlayerRelationEnemy
		);
	}

	isFriendOrNeutral(playerNum: number): boolean {
		return (
			playerNum > 0 &&
			playerNum <= this.relations.length &&
			(this.relations[playerNum - 1].relation === PlayerRelationFriend ||
				this.relations[playerNum - 1].relation === PlayerRelationNeutral)
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
			if (r.relation === PlayerRelationFriend) {
				allies.push(index + 1);
			}
		});
		return allies;
	}

	public getItemCost(
		cs: CS,
		item: ProductionQueueItem | undefined,
		designFinder: DesignFinder,
		planet?: CommandedPlanet,
		quantity = 1
	): Cost {
		if (item) {
			switch (item.type) {
				case QueueItemTypeStarbase: // TODO: starbase upgrades...
					if (item.designNum) {
						const design = designFinder.getMyDesign(item.designNum);
						if (planet?.spec.hasStarbase) {
							const starbaseToUpgrade = designFinder.getMyDesign(planet.spec.starbaseDesignNum);
							if (starbaseToUpgrade && design) {
								return cs.starbaseUpgradeCost(starbaseToUpgrade, design) ?? {};
							}
						}
						return multiply(design?.spec.cost ?? {}, quantity);
					}
					break;
				case QueueItemTypeShipToken:
					if (item.designNum) {
						const design = designFinder.getMyDesign(item.designNum);
						return multiply(design?.spec.cost ?? {}, quantity);
					}
					break;
				default:
					if (this.race?.spec?.costs) {
						return multiply(this.race.spec.costs[item.type] ?? {}, quantity);
					}
			}
		}
		return {};
	}

	// get a player's ability to terraform
	public getTerraformAbility(techStore: TechStore): Hab {
		const terraformAbility: Hab = { grav: 0, temp: 0, rad: 0 };
		const bestTT = getBestTerraform(techStore, this, TerraformHabTypeAll);
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

	public getByHandTransfer(target: MapObjectTarget): Cargo {
		const key = string(target.targetPosition);
		const transfers = this.cargoTransfers[key];
		let cargo = emptyCargo();
		if (!transfers) {
			return cargo;
		}

		// sum up all transfers for this target
		transfers.filter((t) => targetsEqual(target, t)).forEach((t) => (cargo = add(cargo, t.cargo)));

		return cargo;
	}
}

export function canLearnTech(player: CommandedPlayer, tech: Tech): boolean {
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
