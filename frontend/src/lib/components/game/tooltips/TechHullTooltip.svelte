<script lang="ts" module>
	import { showTooltip } from '$lib/services/Stores';
	import type { TechHull } from '$lib/types/cs-proto';
	import TechHullTooltip from './TechHullTooltip.svelte';

	export function onTechHullTooltip(e: PointerEvent | MouseEvent, hull: TechHull | undefined) {
		if (hull) {
			showTooltip<TechHullTooltipProps>(e.x, e.y, TechHullTooltip, { hull: hull });
		}
	}

	export type TechHullTooltipProps = {
		hull: TechHull;
	};
</script>

<script lang="ts">
	import Hull from '../design/Hull.svelte';

	let { hull }: TechHullTooltipProps = $props();
</script>

<div class="flex flex-row justify-center">
	<Hull {hull} cargoCapacity={hull.cargoCapacity} />
</div>
