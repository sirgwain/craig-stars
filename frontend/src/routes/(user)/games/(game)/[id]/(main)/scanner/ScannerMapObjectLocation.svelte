<script lang="ts">
	import type { MapObject } from '$lib/types/MapObject';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import { fade } from 'svelte/transition';

	const { xGet, yGet } = getContext<LayerCake>('LayerCake');

	type Props = {
		show: MapObject | undefined;
	};

	let { show }: Props = $props();
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
