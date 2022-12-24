<script lang="ts">
	import { commandedPlanet } from '$lib/services/Context';
	import CommandTile from './CommandTile.svelte';
	import { onMount } from 'svelte';
	import { EventManager } from '$lib/EventManager';
	import type { Planet } from '$lib/types/Planet';

	onMount(() => {
		const unsubscribe = EventManager.subscribeCargoTransferredEvent((mo) => {
			if ($commandedPlanet == mo) {
				// trigger a reaction
				$commandedPlanet.cargo = (mo as Planet).cargo;
			}
		});

		return () => unsubscribe();
	});
</script>

{#if $commandedPlanet}
	<CommandTile title="Minerals on Hand">
		<div class="flex justify-between">
			<div class="text-ironium">Ironium</div>
			<div>{$commandedPlanet.cargo?.ironium ?? 0}kT</div>
		</div>
		<div class="flex justify-between">
			<div class="text-boranium">Boranium</div>
			<div>{$commandedPlanet.cargo?.boranium ?? 0}kT</div>
		</div>
		<div class="flex justify-between">
			<div class="text-germanium">Germanium</div>
			<div>{$commandedPlanet.cargo?.germanium ?? 0}kT</div>
		</div>

		<div class="divider p-0 m-0" />

		<div class="flex justify-between">
			<div>Mines</div>
			<div>{$commandedPlanet.mines} of {$commandedPlanet.spec?.maxMines}</div>
		</div>
		<div class="flex justify-between">
			<div>Factories</div>
			<div>{$commandedPlanet.factories} of {$commandedPlanet.spec?.maxFactories}</div>
		</div>
	</CommandTile>
{/if}
