<script lang="ts">
	import { EventManager } from '$lib/EventManager';
	import { commandedFleet } from '$lib/services/Context';
	import CargoBar from '$lib/components/game/CargoBar.svelte';
	import FuelBar from '$lib/components/game/FuelBar.svelte';
	import CommandTile from './CommandTile.svelte';
	import type { Fleet } from '$lib/types/Fleet';
	import { onMount } from 'svelte';

	onMount(() => {
		const unsubscribe = EventManager.subscribeCargoTransferredEvent((mo) => {
			if ($commandedFleet == mo) {
				// trigger a reaction
				$commandedFleet.cargo = (mo as Fleet).cargo;
			}
		});

		return () => unsubscribe();
	});

	const transfer = () => {
		if ($commandedFleet) {
			EventManager.publishCargoTransferDialogRequestedEvent($commandedFleet);
		}
	};
</script>

{#if $commandedFleet?.spec}
	<CommandTile title="Fuel & Cargo">
		<div class="flex justify-between my-1">
			<div class="w-12">Fuel</div>
			<div class="ml-1 h-full w-full">
				<FuelBar value={$commandedFleet.fuel} capacity={$commandedFleet.spec.fuelCapacity} />
			</div>
		</div>

		<div class="flex justify-between my-1">
			<div class="w-12">Cargo</div>
			<div class="ml-1 h-full w-full">
				<CargoBar
					on:cargo-transfer={() => transfer()}
					value={$commandedFleet.cargo}
					capacity={$commandedFleet.spec.cargoCapacity}
				/>
			</div>
		</div>
		<div class="flex justify-between">
			<div class="text-ironium">Ironium</div>
			<div>{$commandedFleet.cargo?.ironium ?? 0}kT</div>
		</div>
		<div class="flex justify-between">
			<div class="text-boranium">Boranium</div>
			<div>{$commandedFleet.cargo?.boranium ?? 0}kT</div>
		</div>
		<div class="flex justify-between">
			<div class="text-germanium">Germanium</div>
			<div>{$commandedFleet.cargo?.germanium ?? 0}kT</div>
		</div>
		<div class="flex justify-between">
			<div class="text-colonists">Colonists</div>
			<div>{$commandedFleet.cargo?.colonists ?? 0}kT</div>
		</div>
	</CommandTile>
{/if}
