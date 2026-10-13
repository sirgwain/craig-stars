// How a hit on a token is marked on the battle board: shield damage in blue with the bubble of
// an energy shield, armor damage in orange with a solid plate.

// 24x24 icon paths
export const SHIELD_ICON_PATH = 'M12 4a8 8 0 1 0 0 16a8 8 0 1 0 0-16Z';
export const ARMOR_ICON_PATH = 'M12 2 20 5v6c0 5.5-3.4 9.4-8 11-4.6-1.6-8-5.5-8-11V5Z';

export const SHIELD_COLOR = '#8fe3ff';
export const ARMOR_COLOR = '#ffb050';

export type DamageKind = 'shield' | 'armor';

export type DamageMark = {
	kind: DamageKind;
	amount: number;
	text: string;
	color: string;
	label: string;
};

type Hit = {
	damageDoneShields?: number;
	damageDoneArmor?: number;
};

export function formatDamage(amount: number): string {
	return `-${amount}`;
}

// the marks for one hit, armor first, then shields
export function damageMarks(hit: Hit | undefined): DamageMark[] {
	const marks: DamageMark[] = [];
	const armor = hit?.damageDoneArmor ?? 0;
	const shields = hit?.damageDoneShields ?? 0;
	if (armor > 0) {
		marks.push({
			kind: 'armor',
			amount: armor,
			text: formatDamage(armor),
			color: ARMOR_COLOR,
			label: `${armor} armor damage`
		});
	}
	if (shields > 0) {
		marks.push({
			kind: 'shield',
			amount: shields,
			text: formatDamage(shields),
			color: SHIELD_COLOR,
			label: `${shields} shield damage`
		});
	}
	return marks;
}

type Shot = Hit & { to?: { x?: number; y?: number } };

export type TargetDamage = { x: number; y: number; marks: DamageMark[] };

// the marks of an attack for each square it hit: the shots of one weapon slot can hit several
// targets, or the same one more than once
export function damageByTarget(shots: Shot[]): TargetDamage[] {
	const hits = new Map<string, Required<Hit> & { x: number; y: number }>();
	for (const shot of shots) {
		const x = shot.to?.x ?? 0;
		const y = shot.to?.y ?? 0;
		const key = `${x},${y}`;
		const hit = hits.get(key) ?? {
			x,
			y,
			damageDoneShields: 0,
			damageDoneArmor: 0
		};
		hit.damageDoneShields += shot.damageDoneShields ?? 0;
		hit.damageDoneArmor += shot.damageDoneArmor ?? 0;
		hits.set(key, hit);
	}
	return Array.from(hits.values(), (hit) => ({
		x: hit.x,
		y: hit.y,
		marks: damageMarks(hit)
	})).filter((target) => target.marks.length > 0);
}
