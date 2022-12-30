<script lang="ts">
	import { commandedPlanet } from '$lib/services/Context';
	import CommandTile from './CommandTile.svelte';
	import { onMount } from 'svelte';
	import { EventManager } from '$lib/EventManager';
	import type { Planet } from '$lib/types/Planet';
	export let planet: Planet;

	onMount(() => {
		const unsubscribe = EventManager.subscribeCargoTransferredEvent((mo) => {
			if (planet == mo) {
				// trigger a reaction
				planet.cargo = (mo as Planet).cargo;
			}
		});

		return () => unsubscribe();
	});
</script>

{#if planet}
	<CommandTile title="Minerals on Hand">
		<div class="flex justify-between">
			<div class="text-ironium">Ironium</div>
			<div>{planet.cargo?.ironium ?? 0}kT</div>
		</div>
		<div class="flex justify-between">
			<div class="text-boranium">Boranium</div>
			<div>{planet.cargo?.boranium ?? 0}kT</div>
		</div>
		<div class="flex justify-between">
			<div class="text-germanium">Germanium</div>
			<div>{planet.cargo?.germanium ?? 0}kT</div>
		</div>

		<div class="divider p-0 m-0" />

		<div class="flex justify-between">
			<div>Mines</div>
			<div>{planet.mines} of {planet.spec?.maxMines}</div>
		</div>
		<div class="flex justify-between">
			<div>Factories</div>
			<div>{planet.factories} of {planet.spec?.maxFactories}</div>
		</div>
	</CommandTile>
{/if}
