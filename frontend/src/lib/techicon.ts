import { kebabCase } from 'lodash-es';
import type { ShipDesign } from './types/ShipDesign';

export function getHullIcon(design: ShipDesign | undefined): string {
	if (!design) {
		return '';
	}
	return `hull-${kebabCase(design.hull)}-${design.hullSetNumber ?? 0}`;
}
