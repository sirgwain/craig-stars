<script lang="ts">
	import { page } from '$app/stores';
	import { getGameContext } from '$lib/services/Contexts';
	import { commandMapObject, commandedFleet } from '$lib/services/Stores';
	import type { CommandedFleet, Fleet } from '$lib/types/Fleet';
	import MergeFleets from '../../../dialogs/merge/MergeFleets.svelte';

	const { game, player, universe } = getGameContext();
	let num = parseInt($page.params.num);

	let fleetsInOrbit: Fleet[] = [];

	$: {
		if ($commandedFleet && $commandedFleet.num === num) {
			fleetsInOrbit = $universe
				.getMyFleetsByPosition($commandedFleet)
				.filter((mo) => mo.num !== $commandedFleet?.num) as Fleet[];
		} else {
			const fleet = $universe.getFleet($player.num, num);
			if (fleet) {
				commandMapObject(fleet);
			}
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
