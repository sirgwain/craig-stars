import type { Planet } from './Planet';
import type { Race } from './Race';
import type { ShipDesign } from './ShipDesign';

export interface Player {
	id: number;
	createdAt: string;
	updatedat: string;
	deletedAt: null;
	gameId: number;
	userId: number;
	num: number;
	race: Race;
	ready?: boolean;
	aIControlled?: boolean;
	submittedTurn?: boolean;
	techLevels: TechLevels;
	techLevelsSpent: TechLevels;
	designs: ShipDesign[];
	planets: Planet[];
	planetIntels: Planet[];
}

export interface TechLevels {
	energy?: number;
	weapons?: number;
	propulsion?: number;
	construction?: number;
	electronics?: number;
	biotechnology?: number;
}
