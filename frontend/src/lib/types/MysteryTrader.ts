import type { MovingMapObject } from './MapObject';
import type { TechLevel } from './TechLevel';

export type MysteryTrader = {
	requestedBoon: number;
} & MovingMapObject;

export type MysteryTraderRewardType =
	(typeof MysteryTraderRewardTypes)[keyof typeof MysteryTraderRewardTypes];

export const MysteryTraderRewardTypes = {
	None: '',
	Research: 'Research',
	Engine: 'Engine',
	Bomb: 'Bomb',
	Armor: 'Armor',
	Shield: 'Shield',
	Electrical: 'Electrical',
	Mechanical: 'Mechanical',
	Torpedo: 'Torpedo',
	MineRobot: 'MineRobot',
	ShipHull: 'ShipHull',
	BeamWeapon: 'BeamWeapon',
	Genesis: 'Genesis',
	JumpGate: 'JumpGate',
	Lifeboat: 'Lifeboat'
} as const;

export type MysteryTraderReward = {
	type: MysteryTraderRewardType;
	techLevels: TechLevel;
	tech?: string;
	ship?: string;
	shipCount?: number;
};

export function isHullComponent(type: MysteryTraderRewardType): boolean {
	switch (type) {
		case MysteryTraderRewardTypes.Engine:
		case MysteryTraderRewardTypes.Bomb:
		case MysteryTraderRewardTypes.Armor:
		case MysteryTraderRewardTypes.Shield:
		case MysteryTraderRewardTypes.Electrical:
		case MysteryTraderRewardTypes.Mechanical:
		case MysteryTraderRewardTypes.Torpedo:
		case MysteryTraderRewardTypes.MineRobot:
		case MysteryTraderRewardTypes.BeamWeapon:
			return true;
	}
	return false;
}
