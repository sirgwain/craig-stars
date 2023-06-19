import type { MapObject } from './MapObject';

export type MineField = {
	mineFieldType: MineFieldType;
	numMines: number;
	spec: MineFieldSpec;
} & MapObject &
	MineFieldOrders;

export type MineFieldOrders = {
	detonate?: boolean;
};

export type MineFieldSpec = {
	radius: number;
	decayRate: number;
};

export enum MineFieldType {
	Standard = 'Standard',
	Heavy = 'Heavy',
	SpeedBump = 'SpeedBump'
}

export type MineFieldStats = {
	minDamagePerFleetRS: number;
	damagePerEngineRS: number;
	maxSpeed: number;
	chanceOfHit: number;
	minDamagePerFleet: number;
	damagePerEngine: number;
	sweepFactor: number;
	minDecay: number;
	canDetonate: boolean;
};
