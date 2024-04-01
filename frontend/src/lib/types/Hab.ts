export type HabType = (typeof HabTypes)[keyof typeof HabTypes];

export const HabTypes = {
	Gravity: 0,
	Temperature: 1,
	Radiation: 2
} as const;

export type Hab = {
	grav?: number;
	temp?: number;
	rad?: number;
};

export function habTypeString(type: HabType): string {
	switch (type) {
		case HabTypes.Gravity:
			return 'Gravity';
		case HabTypes.Temperature:
			return 'Temperature';
		case HabTypes.Radiation:
			return 'Radiation';
	}
}

export function getHabValue(hab: Hab | undefined, type: HabType): number {
	switch (type) {
		case HabTypes.Gravity:
			return hab?.grav ?? 0;
		case HabTypes.Temperature:
			return hab?.temp ?? 0;
		case HabTypes.Radiation:
			return hab?.rad ?? 0;
	}
}

export function withHabValue(type: HabType, value: number): Hab {
	switch (type) {
		case HabTypes.Gravity:
			return { grav: value };
		case HabTypes.Temperature:
			return { temp: value };
		case HabTypes.Radiation:
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
	const value = Math.pow(4, (grav - 50) / 50);
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
		case HabTypes.Gravity:
			return getGravString(value);
		case HabTypes.Temperature:
			return getTempString(value);
		case HabTypes.Radiation:
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
			return HabTypes.Gravity;
		} else {
			return HabTypes.Radiation;
		}
	} else {
		if (hab.temp >= hab.rad) {
			return HabTypes.Temperature;
		} else {
			return HabTypes.Radiation;
		}
	}
}

export function absSum(hab: Hab): number {
	return Math.abs(hab.grav ?? 0) + Math.abs(hab.temp ?? 0) + Math.abs(hab.rad ?? 0);
}
