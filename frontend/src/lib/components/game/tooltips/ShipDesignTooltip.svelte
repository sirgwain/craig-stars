<script lang="ts" context="module">
	import { showTooltip } from '$lib/services/Stores';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import ShipDesignTooltip from './ShipDesignTooltip.svelte';

	export function onShipDesignTooltip(e: PointerEvent, design: ShipDesign | undefined) {
		if (design) {
			showTooltip<ShipDesignTooltipProps>(e.x, e.y, ShipDesignTooltip, { design });
		}
	}

	export type ShipDesignTooltipProps = {
		design: ShipDesign;
	};
</script>

<script lang="ts">
	import Design from '$lib/components/game/design/Design.svelte';

	export let design: ShipDesign;
</script>

<div>
	<Design {design} />
</div>
