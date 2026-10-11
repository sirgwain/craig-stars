import { showTooltip } from '#lib/services/Stores.js';
import type { TokenDamage } from '#lib/types/Fleet.js';
import type { ShipDesign } from '#lib/types/cs-proto.js';
import ShipDesignTooltip from './ShipDesignTooltip.svelte';

export function onShipDesignTooltip(
	e: PointerEvent | MouseEvent,
	design: ShipDesign | undefined,
	damage?: TokenDamage
) {
	e.preventDefault();
	if (design) {
		showTooltip<ShipDesignTooltipProps>(e.x, e.y, ShipDesignTooltip, { design, damage });
	}
}

export type ShipDesignTooltipProps = {
	design: ShipDesign;
	damage?: TokenDamage;
};
