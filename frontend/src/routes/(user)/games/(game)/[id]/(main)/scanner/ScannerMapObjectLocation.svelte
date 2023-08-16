<script lang="ts">
	import type { MapObject } from '$lib/types/MapObject';
	import { fade, fly } from 'svelte/transition';

	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';

	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');

	export let show: MapObject | undefined;
</script>

{#if show}
	<g out:fade>
		<path
			d={`M${$xGet(show) - 1000},${$yGet(show) - 1000} l2000,2000`}
			class="show-location-line"
		/>
		<path
			d={`M${$xGet(show) + 1000},${$yGet(show) - 1000} l-2000,2000`}
			class="show-location-line"
		/>
	</g>
{/if}
