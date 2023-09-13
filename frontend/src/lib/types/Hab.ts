export enum HabType {
	Gravity = 'Gravity',
	Temperature = 'Temperature',
	Radiation = 'Radiation'
}

export interface Hab {
	grav?: number;
	temp?: number;
	rad?: number;
}

export function getHabValue(hab: Hab | undefined, type: HabType): number {
	switch (type) {
		case HabType.Gravity:
			return hab?.grav ?? 0;
		case HabType.Temperature:
			return hab?.temp ?? 0;
		case HabType.Radiation:
			return hab?.rad ?? 0;
	}
}

export function withHabValue(type: HabType, value: number): Hab {
	switch (type) {
		case HabType.Gravity:
			return { grav: value };
		case HabType.Temperature:
			return { temp: value };
		case HabType.Radiation:
			return { rad: value };
	}
}

export function add(h1: Hab, h2: Hab) {
	return {
		grav: (h1.grav ?? 0) + (h2.grav ?? 0),
		temp: (h1.temp ?? 0) + (h2.temp ?? 0),
		rad: (h1.rad ?? 0) + (h2.rad ?? 0)
	};
}

// get gravity as a string. This goes from 0.25g to 4.00g
// math credit @ekolis
export function getGravString(grav: number): string {
	const value = Math.pow(4, (grav - 50) / 50)
	return `${value.toFixed(2)}g`;
}

export function getTempString(temp: number): string {
	return `${(temp - 50) * 4}°C`;
}

export function getRadString(rad: number): string {
	return `${rad}mR`;
}

export function getHabValueString(habType: HabType, value: number): string {
	switch (habType) {
		case HabType.Gravity:
			return getGravString(value);
		case HabType.Temperature:
			return getTempString(value);
		case HabType.Radiation:
			return getRadString(value);
	}
	return `${value}`;
}
