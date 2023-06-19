<!--
  @component
  Generates an SVG scatter plot. This component can also work if the x- or y-scale is ordinal, i.e. it has a `.bandwidth` method. See the [timeplot chart](https://layercake.graphics/example/Timeplot) for an example.
 -->
<script lang="ts">
	import { MapObjectType, type MapObject } from '$lib/types/MapObject';
	import type { Planet } from '$lib/types/Planet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import ScannerPlanet from './ScannerPlanet.svelte';

	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');

	let planets: Planet[] = [];

	$: planets = $data && $data.filter((mo: MapObject) => mo.type == MapObjectType.Planet);
</script>

<!-- Planets -->
{#each planets as planet}
	<ScannerPlanet {planet} />
{/each}
