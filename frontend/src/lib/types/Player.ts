import type { Fleet } from './Fleet';
import { MapObjectType, type MapObject } from './MapObject';
import type { Planet } from './Planet';
import type { Race } from './Race';
import type { ShipDesign } from './ShipDesign';

export interface Player {
	id?: number;
	createdAt?: string;
	updatedAt?: string;

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
	fleetIntels: Fleet[];
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
	updatedAt: string;

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

export const findMyPlanet = (player: Player, planet: Planet): Planet | undefined =>
	player?.planets?.find((p) => p.num == planet.num);

export const findIntelMapObject = (player: Player, mo: MapObject): MapObject | undefined => {
	if (mo.type === MapObjectType.Planet) {
		return player?.planetIntels?.find((p) => p.num == mo.num) ?? mo;
	} else if (mo.type === MapObjectType.Fleet) {
		return player?.fleetIntels?.find((f) => f.num == mo.num) ?? mo;
	}
	return mo;
};
