<script lang="ts">
	import { page } from '$app/stores';
	import { commandMapObject, game } from '$lib/services/Context';
	import { PlanetService } from '$lib/services/PlanetService';
	import type { CommandedPlanet } from '$lib/types/Planet';
	import { onMount } from 'svelte';
	import ProductionQueue from '../../../dialogs/ProductionQueue.svelte';
	import { DesignService } from '$lib/services/DesignService';
	import type { ShipDesign } from '$lib/types/ShipDesign';

	let gameId = parseInt($page.params.id);
	let num = parseInt($page.params.num);

	let planet: CommandedPlanet | undefined;
	let designs: ShipDesign[] | undefined;
	let error = '';

	onMount(async () => {
		try {
			await Promise.all([
				PlanetService.get(gameId, num).then((p) => {
					planet = p;
					commandMapObject(planet);
				}),
				DesignService.load(gameId).then((items) => (designs = items))
			]);
		} catch (e) {
			error = (e as Error).message;
		}
	});
</script>

{#if planet && designs && $game}
	<ProductionQueue {planet} {designs} game={$game} player={$game.player} />
{/if}
