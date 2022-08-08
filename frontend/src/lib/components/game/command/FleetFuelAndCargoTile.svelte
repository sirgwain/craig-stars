<script lang="ts">
	import { EventManager } from '$lib/EventManager';
	import { commandedFleet } from '$lib/services/Context';
	import CargoBar from '../CargoBar.svelte';
	import FuelBar from '../FuelBar.svelte';
	import CommandTile from './CommandTile.svelte';

	const transfer = () => {
		if ($commandedFleet) {
			EventManager.publishCargoTransferDialogRequestedEvent($commandedFleet);
		}
	};
</script>

{#if $commandedFleet}
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
			<div class="text-germanium">Colonists</div>
			<div>{$commandedFleet.cargo?.colonists ?? 0}kT</div>
		</div>
	</CommandTile>
{/if}
