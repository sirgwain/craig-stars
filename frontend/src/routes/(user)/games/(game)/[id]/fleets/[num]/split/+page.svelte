<script lang="ts">
	import { page } from '$app/state';
	import { getGameContext } from '$lib/services/GameContext';
	import { onMount } from 'svelte';
	import SplitFleet from '../../../dialogs/split/SplitFleet.svelte';

	const { universe, commandMapObject, commandedFleet } = getGameContext();
	let num = parseInt(page.params.num);

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
	<SplitFleet src={$commandedFleet} />
{/if}
