<script lang="ts">
	import { type MapObject } from '$lib/types/MapObject';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Readable } from 'svelte/store';

	const { xGet, yGet } = getContext<LayerCake>('LayerCake');
	const objectScale = getContext<Readable<number>>('objectScale');

	interface Props {
		mapObject: MapObject;
		children?: import('svelte').Snippet;
	}

	let { mapObject, children }: Props = $props();
</script>

<!-- reverse the scale for MapObjects, we want to zoom in, but not scale the objects themselves. When you zoom 
    in you should be able to click between objects to select things that are close to each other. -->
<g transform={`translate(${$xGet(mapObject)}, ${$yGet(mapObject)}), scale(${1 / $objectScale})`}>
	{#if children}{@render children()}{:else}MapObject{/if}
</g>
