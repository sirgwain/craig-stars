import { Grav, Rad, Temp, type Hab, type HabType } from './cs';

export const HabTypes: HabType[] = [Grav, Temp, Rad] as const;

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

export function getHabValue(hab: Hab | undefined, type: HabType): number {
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
			return { grav: value };
		case Temp:
			return { temp: value };
		case Rad:
			return { rad: value };
		default:
			throw new Error(`Invalid habType: ${type}`);
	}
}

export function add(h1: Hab, h2: Hab) {
	return {
		grav: (h1.grav ?? 0) + (h2.grav ?? 0),
		temp: (h1.temp ?? 0) + (h2.temp ?? 0),
		rad: (h1.rad ?? 0) + (h2.rad ?? 0)
	};
}

export function getGravString(grav: number): string {
	let result = 0;
	const tmp = Math.abs(grav - 50);
	if (tmp <= 25) result = (tmp + 25) * 4;
	else result = tmp * 24 - 400;
	if (grav < 50) result = Math.floor(10000 / result);

	const value = result + (result % 100) / 100.0;

	return `${(value / 100).toFixed(2)}g`;
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

export function getLargest(hab: Hab): HabType {
	hab.grav = hab.grav ?? 0;
	hab.temp = hab.temp ?? 0;
	hab.rad = hab.rad ?? 0;
	if (hab.grav >= hab.temp) {
		if (hab.grav >= hab.rad) {
			return Grav;
		} else {
			return Rad;
		}
	} else {
		if (hab.temp >= hab.rad) {
			return Temp;
		} else {
			return Rad;
		}
	}
}

export function absSum(hab: Hab): number {
	return Math.abs(hab.grav ?? 0) + Math.abs(hab.temp ?? 0) + Math.abs(hab.rad ?? 0);
}
