<script lang="ts">
	import { page } from '$app/stores';
	import { getGameContext } from '$lib/services/GameContext';
	import type { CommandedFleet } from '$lib/types/Fleet';
	import SplitFleet from '../../../dialogs/split/SplitFleet.svelte';

	const { player, universe, commandMapObject, commandedFleet } = getGameContext();
	let num = parseInt($page.params.num);

	$: {
		if (!$commandedFleet || $commandedFleet.num !== num) {
			const fleet = $universe.getFleet($player.num, num);
			if (fleet) {
				commandMapObject(fleet);
			}
		}
	}

	function split(src: CommandedFleet, dest: CommandedFleet) {
		console.log('split, src: ', src, ' dest: ', dest);
	}
</script>

{#if $commandedFleet}
	<SplitFleet src={$commandedFleet} on:ok={(e) => split(e.detail.fleet, e.detail.fleetNums)} />
{/if}
