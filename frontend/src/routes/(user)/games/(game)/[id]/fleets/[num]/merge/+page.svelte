<script lang="ts">
	import { page } from '$app/state';
	import { getGameContext } from '$lib/services/GameContext';
	import type { Fleet } from '$lib/types/cs';
	import { onMount } from 'svelte';
	import MergeFleets from '../../../dialogs/merge/MergeFleets.svelte';

	const { universe, commandedFleet, commandMapObject, merge } = getGameContext();
	let num = parseInt(page.params.num);

	let fleetsInOrbit: Fleet[] = $derived.by(() => {
		if ($commandedFleet && $commandedFleet.num === num) {
			return $universe
				.getMyFleetsByPosition($commandedFleet)
				.filter((mo) => mo.num !== $commandedFleet?.num) as Fleet[];
		}
		return [];
	});

	onMount(() => {
		if (!$commandedFleet || $commandedFleet.num !== num) {
			const fleet = $universe.getMyFleet(num);
			if (fleet) {
				commandMapObject(fleet);
			}
		}
	});
</script>

{#if $commandedFleet}
	<MergeFleets
		fleet={$commandedFleet}
		otherFleetsHere={fleetsInOrbit}
		onOk={(e) => merge(e.fleet, e.fleetNums)}
	/>
{/if}
