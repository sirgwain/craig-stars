import { showTooltip } from '$lib/services/Stores';
import type { TechLike } from '$lib/types/Tech';
import TechTooltip from './TechTooltip.svelte';

export function onTechTooltip(
	e: PointerEvent | MouseEvent,
	tech: TechLike | undefined,
	showResearchCost = false
) {
	e.preventDefault();
	if (tech) {
		showTooltip<TechTooltipProps>(e.x, e.y, TechTooltip, { tech, showResearchCost });
	}
}

export type TechTooltipProps = {
	tech: TechLike | undefined;
	showResearchCost?: boolean;
};
