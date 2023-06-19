<script lang="ts" context="module">
	import { showTooltip } from '$lib/services/Context';
	import type { ShipDesign, ShipDesignIntel } from '$lib/types/ShipDesign';
	import ShipDesignTooltip from './ShipDesignTooltip.svelte';

	export function onShipDesignTooltip(
		e: PointerEvent,
		design: ShipDesign | ShipDesignIntel | undefined
	) {
		if (design) {
			showTooltip<ShipDesignTooltipProps>(e.x, e.y, ShipDesignTooltip, { design });
		}
	}

	export type ShipDesignTooltipProps = {
		design: ShipDesign | ShipDesignIntel;
	};
</script>

<script lang="ts">
	import Design from '$lib/components/game/design/Design.svelte';

	export let design: ShipDesign | ShipDesignIntel;
</script>

<div>
	<Design {design} />
</div>
