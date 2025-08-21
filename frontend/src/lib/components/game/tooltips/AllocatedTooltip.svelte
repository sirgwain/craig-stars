<script lang="ts" module>
	import { showTooltip } from '$lib/services/Stores';
	import type { Cost } from '$lib/types/cs-proto';
	import AllocatedTooltip from './AllocatedTooltip.svelte';

	export function onAllocatedTooltip(e: PointerEvent | MouseEvent, cost: Cost | undefined) {
		if (cost) {
			showTooltip<AllocatedTooltipProps>(e.x, e.y, AllocatedTooltip, { cost });
		}
	}

	export type AllocatedTooltipProps = {
		cost: Cost;
	};
</script>

<script lang="ts">
	import CostMini from '../CostMini.svelte';

	let { cost }: AllocatedTooltipProps = $props();
</script>

<div>
	Allocated <CostMini {cost} />
</div>
