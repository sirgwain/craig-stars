import { kebabCase } from 'lodash-es';
import type { TechLike } from './types/Tech';
import { TechCategory, type ShipDesign } from './types/cs-proto';

export function getHullIcon(design: ShipDesign | undefined): string {
	if (!design) {
		return '';
	}
	return `hull-${kebabCase(design.hull)}-${design.hullSetNumber ?? 0}`;
}

export function getTechIcon(techLike: TechLike | undefined, hullSetNumber: number): string {
	if (!techLike?.tech) {
		return '';
	}
	const raw = (techLike.tech.name ?? '').replace("'", '').replace(' ', '').replace('±', '');
	const name = kebabCase(raw);
	const category = techLike.tech.category;
	if (category === TechCategory.SHIP_HULL || category === TechCategory.STARBASE_HULL) {
		return `hull-${name}-${hullSetNumber ?? 0}`;
	}
	return `${name}`;
}
