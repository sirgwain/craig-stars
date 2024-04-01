<script lang="ts">
	import { page } from '$app/stores';
	import { getGameContext } from '$lib/services/GameContext';
	import { ownedBy } from '$lib/types/MapObject';
	import ProductionQueue from '../../../dialogs/production/ProductionQueue.svelte';

	const { player, universe, commandedPlanet, commandMapObject } = getGameContext();
	let num = parseInt($page.params.num);

	$: {
		const planet = $universe.getPlanet(num);
		if (planet && ownedBy(planet, $player.num)) {
			commandMapObject(planet);
		}
	}
</script>

{#if $commandedPlanet}
	<ProductionQueue planet={$commandedPlanet} />
{/if}
