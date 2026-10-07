import { showTooltip } from '#lib/services/Stores.js';
import type { Cost } from '#lib/types/cs-proto.js';
import AllocatedTooltip from './AllocatedTooltip.svelte';

export function onAllocatedTooltip(e: PointerEvent | MouseEvent, cost: Cost | undefined) {
	if (cost) {
		showTooltip<AllocatedTooltipProps>(e.x, e.y, AllocatedTooltip, { cost });
	}
}

export type AllocatedTooltipProps = {
	cost: Cost;
};
