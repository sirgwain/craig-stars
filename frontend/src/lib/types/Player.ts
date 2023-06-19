import type { BattleAttackWho, BattleRecord, BattleTactic, BattleTarget } from './Battle';
import type { Fleet } from './Fleet';
import type { Planet } from './Planet';
import { humanoid, type Race } from './Race';
import type { ShipDesign, ShipDesignIntel } from './ShipDesign';
import type { Tech, TechDefense, TechPlanetaryScanner } from './Tech';

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
	PlayerIntels &
	PlayerMessages;

export type PlayerMessages = {
	messages: Message[];
};

export type PlayerPlans = {
	battlePlans?: BattlePlan[];
};

export type PlayerIntels = {
	planetIntels: Planet[];
	fleetIntels?: Fleet[];
	shipDesignIntels?: ShipDesignIntel[];
	playerIntels: PlayerIntel[];
};
export type PlayerMapObjects = {
	planets: Planet[];
	fleets: Fleet[];
	starbases: Fleet[];
};

export type PlayerOrders = {
	researching: TechField;
	nextResearchField: NextResearchField;
	researchAmount: number;
};

export type BattlePlan = {
	name: string;
	primaryTarget: BattleTarget;
	secondaryTarget: BattleTarget;
	tactic: BattleTactic;
	attackWho: BattleAttackWho;
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

export interface Message {
	type: string;
	text: string;
	targetType?: MessageTargetType;
	targetNum?: number;
	battleNum?: number;
	targetPlayerNum?: number;
}

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
	planets: Planet[] = [];
	fleets: Fleet[] = [];
	starbases: Fleet[] = [];
	battlePlans: BattlePlan[] = [];
	designs: ShipDesign[] = [];
	planetIntels: Planet[] = [];
	fleetIntels: Fleet[] = [];
	shipDesignIntels: ShipDesignIntel[] = [];
	playerIntels: PlayerIntel[] = [];
	messages: Message[] = [];
	battles: BattleRecord[] = [];
	spec: PlayerSpec = {};

	constructor(data?: PlayerResponse) {
		if (data) {
			Object.assign(this, data);
		}
	}

	getPlayerIntel(num: number): PlayerIntel | undefined {
		if (num >= 1 && num <= this.playerIntels.length) {
			return this.playerIntels[num - 1];
		}
	}

	getDesign(playerNum: number, num: number): ShipDesign | ShipDesignIntel | undefined {
		if (playerNum == this.num) {
			return this.designs.find((d) => d.num === num);
		} else {
			return this.shipDesignIntels.find((d) => d.num === num);
		}
	}

	updateDesign(design: ShipDesign) {
		const filteredDesigns = this.designs.filter((d) => d.num != design.num) ?? [];
		this.designs = [...filteredDesigns, design];
	}

	getPlanetIntel(num: number): Planet | undefined {
		return this.planetIntels.find((p) => p.num === num);
	}

	getBattle(num: number): BattleRecord | undefined {
		return this.battles.find((b) => b.num === num);
	}

	getBattleLocation(battle: BattleRecord): string {
		if (battle.planetNum) {
			const planet = this.getPlanetIntel(battle.planetNum);
			return planet?.name ?? 'Unknown';
		}
		return `Space (${battle.position.x}, ${battle.position.y}`;
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
