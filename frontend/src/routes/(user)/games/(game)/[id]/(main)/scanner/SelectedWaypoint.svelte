<!--
  @component
  Renders the currently selected waypoint line pieces as a slightly larger line
 -->
<script lang="ts">
	import { commandedFleet, selectedWaypoint } from '$lib/services/Context';
	import { equal } from '$lib/types/Vector';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';

	const { xGet, yGet, xScale, yScale } = getContext<LayerCake>('LayerCake');
	const scale = getContext<Writable<number>>('scale');

	$: wpIndex = $commandedFleet?.waypoints.findIndex(
		(wp) => $selectedWaypoint && equal(wp.position, $selectedWaypoint.position)
	);

	$: wpNextIndex =
		wpIndex && $commandedFleet?.waypoints && $commandedFleet.waypoints.length > wpIndex + 1
			? wpIndex + 1
			: undefined;
	$: wpNextIndexLine =
		wpNextIndex && $commandedFleet
			? `L${$xGet($commandedFleet.waypoints[wpNextIndex])}, ${$yGet(
					$commandedFleet.waypoints[wpNextIndex]
			  )}`
			: '';

	$: strokeWidth = 6 / $scale;
</script>

{#if $commandedFleet && $selectedWaypoint && wpIndex && wpIndex != 0}
	<path
		d={`M${$xGet($commandedFleet.waypoints[wpIndex - 1])}, ${$yGet(
			$commandedFleet.waypoints[wpIndex - 1]
		)}L${$xGet($selectedWaypoint)}, ${$yGet($selectedWaypoint)}` + wpNextIndexLine ?? ''}
		class="waypoint-line-selected"
		stroke-width={strokeWidth}
	/>
{/if}
