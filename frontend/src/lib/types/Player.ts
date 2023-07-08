import type { BattleAttackWho, BattleRecord, BattleTactic, BattleTarget } from './Battle';
import type { Fleet, WaypointTransportTasks } from './Fleet';
import type { Message } from './Message';
import type { MineField } from './MineField';
import type { MineralPacket } from './MineralPacket';
import type { MysteryTrader } from './MysteryTrader';
import type { Planet, ProductionQueueItem } from './Planet';
import { humanoid, type Race } from './Race';
import type { Salvage } from './Salvage';
import type { ShipDesign } from './ShipDesign';
import type { Tech, TechDefense, TechPlanetaryScanner } from './Tech';
import { emptyTechLevel, hasRequiredLevels, TechField, type TechLevel } from './TechLevel';
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
	submittedTurn?: boolean;
	techLevels: TechLevel;
	techLevelsSpent: TechLevel;
	designs?: ShipDesign[];
	researchSpentLastYear?: number;
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

export class Player implements PlayerResponse {
	id = 0;
	createdAt?: string | undefined;
	updatedAt?: string | undefined;

	gameId = 0;
	num = 0;

	userId?: number | undefined;
	name = '';
	color = '#00FF00';
	race = { ...humanoid };
	ready = false;
	aiControlled = false;
	submittedTurn = false;
	techLevels: TechLevel = { ...emptyTechLevel() };
	techLevelsSpent: TechLevel = { ...emptyTechLevel() };
	researchSpentLastYear = 0;
	researching: TechField = TechField.Energy;
	nextResearchField: NextResearchField = NextResearchField.Energy;
	researchAmount = 15;
	battlePlans: BattlePlan[] = [];
	productionPlans: ProductionPlan[] = [];
	transportPlans: TransportPlan[] = [];
	messages: Message[] = [];
	relations: PlayerRelationship[] = [];
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
		return canLearnTech(this, tech) && hasRequiredLevels(this.techLevels, tech.requirements)
	}
}

export function canLearnTech(player: PlayerResponse, tech: Tech): boolean {
	const requirements = tech.requirements;
	if (requirements.prtRequired && requirements.prtRequired !== player.race.prt) {
		return false;
	}
	if (requirements.prtDenied && player.race.prt === requirements.prtDenied) {
		return false;
	}
	if (
		requirements.lrtsRequired &&
		(player.race.lrts & (1 << requirements.lrtsRequired)) !== 1 << requirements.lrtsRequired
	) {
		return false;
	}
	if (requirements.lrtsDenied && (player.race.lrts & (1 << requirements.lrtsDenied)) !== 0) {
		return false;
	}
	return true;
}
