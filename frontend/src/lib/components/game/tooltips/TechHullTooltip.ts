import { showTooltip } from '$lib/services/Stores';
import type { TechHull } from '$lib/types/cs-proto';
import TechHullTooltip from './TechHullTooltip.svelte';

export function onTechHullTooltip(e: PointerEvent | MouseEvent, hull: TechHull | undefined) {
	if (hull) {
		showTooltip<TechHullTooltipProps>(e.x, e.y, TechHullTooltip, { hull });
	}
}

export type TechHullTooltipProps = {
	hull: TechHull;
};
