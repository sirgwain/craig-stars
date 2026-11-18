import { showTooltip } from '#lib/services/Stores.js';
import type { TechLike } from '#lib/types/Tech.js';
import TechTooltip from './TechTooltip.svelte';
import type { CommandedPlayer } from '#lib/types/Player.js';
import type { CS } from '#lib/wasm.js';

export function onTechTooltip(
	e: PointerEvent | MouseEvent,
	tech: TechLike | undefined,
	showResearchCost = false,
	environment: Pick<TechTooltipProps, 'playerOverride' | 'client'> = {}
) {
	e.preventDefault();
	if (tech) {
		showTooltip<TechTooltipProps>(e.x, e.y, TechTooltip, {
			tech,
			showResearchCost,
			...environment
		});
	}
}

export type TechTooltipProps = {
	tech: TechLike | undefined;
	showResearchCost?: boolean;
	playerOverride?: CommandedPlayer;
	client?: CS;
};
