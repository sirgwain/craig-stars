import { showTooltip } from '$lib/services/Stores';
import type { ShipDesign } from '$lib/types/cs-proto';
import ShipDesignTooltip from './ShipDesignTooltip.svelte';

export function onShipDesignTooltip(e: PointerEvent | MouseEvent, design: ShipDesign | undefined) {
	e.preventDefault();
	if (design) {
		showTooltip<ShipDesignTooltipProps>(e.x, e.y, ShipDesignTooltip, { design });
	}
}

export type ShipDesignTooltipProps = {
	design: ShipDesign;
};
