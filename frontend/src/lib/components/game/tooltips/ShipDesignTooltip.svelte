<script lang="ts" module>
	import { showTooltip } from '$lib/services/Stores';
	import type { AnyShipDesign } from '$lib/services/Universe';
	import ShipDesignTooltip from './ShipDesignTooltip.svelte';

	export function onShipDesignTooltip(
		e: PointerEvent | MouseEvent,
		design: AnyShipDesign | undefined
	) {
		e.preventDefault();
		if (design) {
			showTooltip<ShipDesignTooltipProps>(e.x, e.y, ShipDesignTooltip, { design });
		}
	}

	export type ShipDesignTooltipProps = {
		design: AnyShipDesign;
	};
</script>

<script lang="ts">
	import Design from '$lib/components/game/design/Design.svelte';

	let { design }: ShipDesignTooltipProps = $props();
</script>

<div>
	<Design {design} />
</div>
