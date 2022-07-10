import type { Fleet } from './Fleet';
import type { Planet } from './Planet';
import type { Race } from './Race';
import type { ShipDesign } from './ShipDesign';

export interface Player {
	id?: number;
	createdAt?: string;
	updatedat?: string;

	gameId: number;
	userId?: number;
	num: number;
	race: Race;
	ready?: boolean;
	aIControlled?: boolean;
	submittedTurn?: boolean;
	techLevels: TechLevel;
	techLevelsSpent: TechLevel;
	messages: Message[];
	designs: ShipDesign[];
	planets: Planet[];
	fleets: Fleet[];
	planetIntels: Planet[];
}

export interface TechLevel {
	energy?: number;
	weapons?: number;
	propulsion?: number;
	construction?: number;
	electronics?: number;
	biotechnology?: number;
}

export interface Message {
	id: number;
	createdAt: string;
	updatedat: string;

	playerId: number;
	type: string;
	text: string;
	targetType: MessageTargetType;
	targetId?: number;
}

export enum MessageTargetType {
	None = 'None',
	Planet = 'Planet',
	Fleet = 'Fleet',
	Wormhole = 'Wormhole',
	MineField = 'MineField',
	MysteryTrader = 'MysteryTrader',
	Battle = 'Battle'
}
