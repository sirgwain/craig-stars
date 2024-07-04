export enum TechField {
	Energy = 'Energy',
	Weapons = 'Weapons',
	Propulsion = 'Propulsion',
	Construction = 'Construction',
	Electronics = 'Electronics',
	Biotechnology = 'Biotechnology'
}

export interface TechLevel {
	energy?: number;
	weapons?: number;
	propulsion?: number;
	construction?: number;
	electronics?: number;
	biotechnology?: number;
}

export const emptyTechLevel = (): TechLevel => ({
	energy: 0,
	weapons: 0,
	propulsion: 0,
	construction: 0,
	electronics: 0,
	biotechnology: 0
});

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

export function minus(tl1: TechLevel, tl2: TechLevel): TechLevel {
	return {
		energy: (tl1.energy ?? 0) - (tl2.energy ?? 0),
		weapons: (tl1.weapons ?? 0) - (tl2.weapons ?? 0),
		propulsion: (tl1.propulsion ?? 0) - (tl2.propulsion ?? 0),
		construction: (tl1.construction ?? 0) - (tl2.construction ?? 0),
		electronics: (tl1.electronics ?? 0) - (tl2.electronics ?? 0),
		biotechnology: (tl1.biotechnology ?? 0) - (tl2.biotechnology ?? 0)
	};
}

export function sum(tl: TechLevel): number {
	return (
		(tl.energy ?? 0) +
		(tl.weapons ?? 0) +
		(tl.propulsion ?? 0) +
		(tl.construction ?? 0) +
		(tl.electronics ?? 0) +
		(tl.biotechnology ?? 0)
	);
}

export function levelsAbove(req: TechLevel, level: TechLevel): number | undefined {
	const diffs: number[] = [];
	if (req.energy) {
		diffs.push((level.energy ?? 0) - req.energy);
	}
	if (req.weapons) {
		diffs.push((level.weapons ?? 0) - req.weapons);
	}
	if (req.propulsion) {
		diffs.push((level.propulsion ?? 0) - req.propulsion);
	}
	if (req.construction) {
		diffs.push((level.construction ?? 0) - req.construction);
	}
	if (req.electronics) {
		diffs.push((level.electronics ?? 0) - req.electronics);
	}
	if (req.biotechnology) {
		diffs.push((level.biotechnology ?? 0) - req.biotechnology);
	}
	return Math.min(...diffs);
}

export function get(tl: TechLevel, field: TechField): number {
	switch (field) {
		case TechField.Energy:
			return tl.energy ?? 0;
		case TechField.Weapons:
			return tl.weapons ?? 0;
		case TechField.Propulsion:
			return tl.propulsion ?? 0;
		case TechField.Construction:
			return tl.construction ?? 0;
		case TechField.Electronics:
			return tl.electronics ?? 0;
		case TechField.Biotechnology:
			return tl.biotechnology ?? 0;
		default:
			throw new Error('invalid field: ' + field);
	}
}

export function minTechLevel(tl: TechLevel): number {
	return Math.min(
		tl.energy ?? 0,
		tl.weapons ?? 0,
		tl.propulsion ?? 0,
		tl.construction ?? 0,
		tl.electronics ?? 0,
		tl.biotechnology ?? 0
	);
}
