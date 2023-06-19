import type { BattleAttackWho, BattleRecord, BattleTactic, BattleTarget } from './Battle';
import type { Fleet, Target, WaypointTransportTasks } from './Fleet';
import type { MineField } from './MineField';
import type { MineralPacket } from './MineralPacket';
import type { MysteryTrader } from './MysteryTrader';
import type { Planet, ProductionQueueItem } from './Planet';
import { humanoid, type Race } from './Race';
import type { Salvage } from './Salvage';
import type { ShipDesign } from './ShipDesign';
import type { Tech, TechDefense, TechPlanetaryScanner } from './Tech';
import type { Wormhole } from './Wormhole';

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
	aIControlled?: boolean;
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
	seen: boolean;
	raceName?: string;
	racePluralName?: string;
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

export enum TechField {
	Energy = 'Energy',
	Weapons = 'Weapons',
	Propulsion = 'Propulsion',
	Construction = 'Construction',
	Electronics = 'Electronics',
	Biotechnology = 'Biotechnology'
}

export type TechLevel = {
	energy?: number;
	weapons?: number;
	propulsion?: number;
	construction?: number;
	electronics?: number;
	biotechnology?: number;
};

const emptyTechLevel: TechLevel = {
	energy: 0,
	weapons: 0,
	propulsion: 0,
	construction: 0,
	electronics: 0,
	biotechnology: 0
};

export type Message = {
	type: number;
	text: string;
	battleNum?: number;
} & Target;

export enum MessageTargetType {
	None = '',
	Planet = 'Planet',
	Fleet = 'Fleet',
	Wormhole = 'Wormhole',
	MineField = 'MineField',
	MysteryTrader = 'MysteryTrader',
	Battle = 'Battle'
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
	aIControlled = false;
	submittedTurn = false;
	techLevels: TechLevel = { ...emptyTechLevel };
	techLevelsSpent: TechLevel = { ...emptyTechLevel };
	researchSpentLastYear = 0;
	researching: TechField = TechField.Energy;
	nextResearchField: NextResearchField = NextResearchField.Energy;
	researchAmount = 15;
	battlePlans: BattlePlan[] = [];
	productionPlans: ProductionPlan[] = [];
	transportPlans: TransportPlan[] = [];
	messages: Message[] = [];
	spec: PlayerSpec = {};

	constructor(data?: PlayerResponse) {
		if (data) {
			Object.assign(this, data);
		}
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
}

export function hasRequiredLevels(tl: TechLevel, required: TechLevel): boolean {
	return (
		(tl.energy ?? 0) >= (required.energy ?? 0) &&
		(tl.weapons ?? 0) >= (required.weapons ?? 0) &&
		(tl.propulsion ?? 0) >= (required.propulsion ?? 0) &&
		(tl.construction ?? 0) >= (required.construction ?? 0) &&
		(tl.electronics ?? 0) >= (required.electronics ?? 0) &&
		(tl.biotechnology ?? 0) >= (required.biotechnology ?? 0)
	);
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
