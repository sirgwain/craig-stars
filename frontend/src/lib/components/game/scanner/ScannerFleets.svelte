<!--
  @component
  Generates an SVG scatter plot. This component can also work if the x- or y-scale is ordinal, i.e. it has a `.bandwidth` method. See the [timeplot chart](https://layercake.graphics/example/Timeplot) for an example.
 -->
<script lang="ts">
	import { MapObjectType, type MapObject } from '$lib/types/MapObject';
	import type { Fleet } from '$lib/types/Fleet';
	import { getContext } from 'svelte';
	import ScannerFleet from './ScannerFleet.svelte';

	const { data } = getContext('LayerCake');

	let fleets: Fleet[] = [];

	$: fleets = $data && $data.filter((mo: MapObject) => mo.type == MapObjectType.Fleet);
</script>

{#each fleets as fleet}
	<ScannerFleet {fleet} />
{/each}
