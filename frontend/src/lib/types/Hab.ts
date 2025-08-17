import { HabSchema, type Hab, type HabJson } from '$lib/types/cs-proto';
import { create } from '@bufbuild/protobuf';

export type HabType = number /* int */;
export const Grav: HabType = 0;
export const Temp: HabType = 1;
export const Rad: HabType = 2;

export function emptyHab(): Hab {
	return create(HabSchema, {});
}

export const HabTypeShortString: string[] = ['grav', 'temp', 'rad'] as const;

export function habTypeString(type: HabType): string {
	switch (type) {
		case Grav:
			return 'Gravity';
		case Temp:
			return 'Temperature';
		case Rad:
			return 'Radiation';
		default:
			throw new Error(`Invalid habType: ${type}`);
	}
}

export function getHabValue(hab: HabJson | undefined, type: HabType): number {
	switch (type) {
		case Grav:
			return hab?.grav ?? 0;
		case Temp:
			return hab?.temp ?? 0;
		case Rad:
			return hab?.rad ?? 0;
		default:
			throw new Error(`Invalid habType: ${type}`);
	}
}

export function withHabValue(type: HabType, value: number): Hab {
	switch (type) {
		case Grav:
			return create(HabSchema, { grav: value });
		case Temp:
			return create(HabSchema, { temp: value });
		case Rad:
			return create(HabSchema, { rad: value });
		default:
			throw new Error(`Invalid habType: ${type}`);
	}
}

export function add(h1: Hab, h2: Hab): Hab {
	return create(HabSchema, {
		grav: (h1.grav ?? 0) + (h2.grav ?? 0),
		temp: (h1.temp ?? 0) + (h2.temp ?? 0),
		rad: (h1.rad ?? 0) + (h2.rad ?? 0)
	});
}

const gravFormatter = new Intl.NumberFormat(undefined, {
	style: 'decimal',
	minimumFractionDigits: 2,
	maximumFractionDigits: 2,
	roundingMode: 'trunc'
});

export function getGravString(grav: number): string {
	const tmp = Math.abs(grav - 50);

	let result = tmp <= 25 ? (tmp + 25) * 4 : tmp * 24 - 400;

	if (grav < 50) {
		result = Math.floor(10000 / result);
	}

	const formatted = gravFormatter.format(result / 100.0);

	return `${formatted}g`;
}

export function getTempString(temp: number): string {
	return `${(temp - 50) * 4}°C`;
}

export function getRadString(rad: number): string {
	return `${rad}mR`;
}

export function getHabValueString(habType: HabType, value: number): string {
	switch (habType) {
		case Grav:
			return getGravString(value);
		case Temp:
			return getTempString(value);
		case Rad:
			return getRadString(value);
	}
	return `${value}`;
}

export function absSum(hab: Hab | undefined): number {
	if (!hab) return 0;
	return Math.abs(hab.grav ?? 0) + Math.abs(hab.temp ?? 0) + Math.abs(hab.rad ?? 0);
}
