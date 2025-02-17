import { kebabCase } from 'lodash-es';
import type { ShipDesign } from './types/cs';
import { TechCategoryShipHull, TechCategoryStarbaseHull, type Tech } from './types/cs';

export function getHullIcon(design: ShipDesign | undefined): string {
	if (!design) {
		return '';
	}
	return `hull-${kebabCase(design.hull)}-${design.hullSetNumber ?? 0}`;
}

export function getTechIcon(tech: Tech | undefined, hullSetNumber: number): string {
	if (!tech) {
		return '';
	}
	const name = kebabCase(tech?.name.replace("'", '').replace(' ', '').replace('±', ''));
	if ([TechCategoryShipHull, TechCategoryStarbaseHull].includes(tech?.category)) {
		return `hull-${name}-${hullSetNumber ?? 0}`;
	} else {
		return `${name}`;
	}
}
