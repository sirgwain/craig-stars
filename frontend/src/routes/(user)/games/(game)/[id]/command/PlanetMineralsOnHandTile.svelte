<script lang="ts">
	import { EventManager } from '$lib/EventManager';
	import type { CommandedPlanet, Planet } from '$lib/types/Planet';
	import { onMount } from 'svelte';
	import CommandTile from './CommandTile.svelte';

	export let planet: CommandedPlanet;

	onMount(() => {
		const unsubscribe = EventManager.subscribeCargoTransferredEvent((mo) => {
			if (planet == mo) {
				// trigger a reaction
				planet.cargo = (mo as CommandedPlanet).cargo;
			}
		});

		return () => unsubscribe();
	});
</script>

{#if planet}
	<CommandTile title="Minerals on Hand">
		<div class="flex justify-between">
			<div class="text-ironium">Ironium</div>
			<div>{planet.cargo.ironium}kT</div>
		</div>
		<div class="flex justify-between">
			<div class="text-boranium">Boranium</div>
			<div>{planet.cargo.boranium}kT</div>
		</div>
		<div class="flex justify-between">
			<div class="text-germanium">Germanium</div>
			<div>{planet.cargo.germanium}kT</div>
		</div>

		<div class="divider p-0 m-0" />

		<div class="flex justify-between">
			<div>Mines</div>
			<div>{planet.mines} of {planet.spec.maxMines}</div>
		</div>
		<div class="flex justify-between">
			<div>Factories</div>
			<div>{planet.factories} of {planet.spec.maxFactories}</div>
		</div>
	</CommandTile>
{/if}
