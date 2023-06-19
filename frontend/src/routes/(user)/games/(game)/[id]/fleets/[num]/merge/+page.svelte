<script lang="ts">
	import { page } from '$app/stores';
	import { commandedFleet, game } from '$lib/services/Context';
	import type { CommandedFleet, Fleet } from '$lib/types/Fleet';
	import { onMount } from 'svelte';
	import MergeFleets from './MergeFleets.svelte';

	let gameId = parseInt($page.params.id);
	let num = parseInt($page.params.num);

	onMount(async () => {
		await $game?.load(gameId);

		// TODO: command the fleet by num
	});

	let fleetsInOrbit: Fleet[] = [];

	$: {
		if ($commandedFleet) {
			fleetsInOrbit = $game?.universe
				.getMyFleetsByPosition($commandedFleet)
				.filter((mo) => mo.num !== $commandedFleet?.num) as Fleet[];
		}
	}

	async function merge(fleet: CommandedFleet, fleetNums: number[]) {
		$game?.merge(fleet, fleetNums);
	}
</script>

{#if $commandedFleet}
	<MergeFleets
		fleet={$commandedFleet}
		otherFleetsHere={fleetsInOrbit}
		on:ok={(e) => merge(e.detail.fleet, e.detail.fleetNums)}
	/>
{/if}
