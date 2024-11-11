<script lang="ts">
	import { type MapObject } from '$lib/types/MapObject';
	import type { LayerCake } from 'layercake';
	import { getContext, type Snippet } from 'svelte';
	import { getScannerContext } from './Scanner';

	const { xGet, yGet } = getContext<LayerCake>('LayerCake');
	const { objectScale } = getScannerContext();

	type Props = {
		mapObject: MapObject;
		children?: Snippet;
	};

	let { mapObject, children }: Props = $props();
</script>

<!-- reverse the scale for MapObjects, we want to zoom in, but not scale the objects themselves. When you zoom 
    in you should be able to click between objects to select things that are close to each other. -->
<g transform={`translate(${$xGet(mapObject)}, ${$yGet(mapObject)}), scale(${1 / $objectScale})`}>
	{#if children}{@render children()}{:else}MapObject{/if}
</g>
