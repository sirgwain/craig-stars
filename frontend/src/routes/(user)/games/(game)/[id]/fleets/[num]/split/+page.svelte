<script lang="ts">
	import { run } from 'svelte/legacy';

	import { page } from '$app/stores';
	import { getGameContext } from '$lib/services/GameContext';
	import SplitFleet from '../../../dialogs/split/SplitFleet.svelte';
	import { onMount } from 'svelte';

	const { player, universe, commandMapObject, commandedFleet } = getGameContext();
	let num = parseInt($page.params.num);

	onMount(() => {
		if (!$commandedFleet || $commandedFleet.num !== num) {
			const fleet = $universe.getFleet($player.num, num);
			if (fleet) {
				commandMapObject(fleet);
			}
		}
	});
</script>

{#if $commandedFleet}
	<SplitFleet src={$commandedFleet} />
{/if}
