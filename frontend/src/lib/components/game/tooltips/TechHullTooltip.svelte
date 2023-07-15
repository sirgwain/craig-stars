<script lang="ts" context="module">
	import { showTooltip } from '$lib/services/Stores';
	import TechHullTooltip from './TechHullTooltip.svelte';

	export function onTechHullTooltip(e: PointerEvent, hull: TechHull | undefined) {
		if (hull) {
			showTooltip<TechHullTooltipProps>(e.x, e.y, TechHullTooltip, { hull: hull });
		}
	}

	export type TechHullTooltipProps = {
		hull: TechHull;
	};
</script>

<script lang="ts">
	import type { TechHull } from '$lib/types/Tech';
	import Hull from '../design/Hull.svelte';

	export let hull: TechHull;
</script>

<div class="flex flex-row justify-center">
	<Hull {hull} cargoCapacity={hull.cargoCapacity} />
</div>
