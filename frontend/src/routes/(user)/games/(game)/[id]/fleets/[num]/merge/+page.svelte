<script lang="ts">
	import { page } from '$app/stores';
	import { getGameContext } from '$lib/services/Contexts';
	import { commandedFleet } from '$lib/services/Stores';
	import type { CommandedFleet, Fleet } from '$lib/types/Fleet';
	import MergeFleets from './MergeFleets.svelte';

	const { game, player, universe } = getGameContext();
	let num = parseInt($page.params.num);

	let fleetsInOrbit: Fleet[] = [];

	$: {
		if ($commandedFleet) {
			fleetsInOrbit = $universe
				.getMyFleetsByPosition($commandedFleet)
				.filter((mo) => mo.num !== $commandedFleet?.num) as Fleet[];
		}
	}

	async function merge(fleet: CommandedFleet, fleetNums: number[]) {
		$game.merge(fleet, fleetNums);
	}
</script>

{#if $commandedFleet}
	<MergeFleets
		fleet={$commandedFleet}
		otherFleetsHere={fleetsInOrbit}
		on:ok={(e) => merge(e.detail.fleet, e.detail.fleetNums)}
	/>
{/if}
