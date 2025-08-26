<script lang="ts">
	import { key, type MapObjectLike } from '$lib/types/MapObject';
	import type { LayerCake } from 'layercake';
	import { getContext, type Snippet } from 'svelte';
	import { getScannerContext } from './Scanner';

	const { xGet, yGet } = getContext<LayerCake>('LayerCake');
	const { objectScale } = getScannerContext();

	type Props = {
		mapObject: MapObjectLike;
		hideDataId?: boolean;
		children?: Snippet;
	};

	let { mapObject, hideDataId, children }: Props = $props();
</script>

<!-- reverse the scale for MapObjects, we want to zoom in, but not scale the objects themselves. When you zoom 
    in you should be able to click between objects to select things that are close to each other. -->
<g
	data-id={!hideDataId && key(mapObject)}
	transform={`translate(${$xGet(mapObject.mapObject)}, ${$yGet(mapObject.mapObject)}), scale(${1 / $objectScale})`}
>
	{#if children}{@render children()}{:else}MapObject{/if}
</g>
