<script lang="ts">
	import CargoBar from '$lib/components/game/CargoBar.svelte';
	import FuelBar from '$lib/components/game/FuelBar.svelte';
	import type { ShowCargoTransferDialogProps } from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import { canTransferCargo, type CommandedFleet } from '$lib/types/Fleet';
	import CommandTile from './CommandTile.svelte';

	const { universe } = getGameContext();

	type Props = {
		fleet: CommandedFleet;
	} & ShowCargoTransferDialogProps;

	let { fleet, onShowCargoTransferDialog }: Props = $props();

	function transfer() {
		if (!onShowCargoTransferDialog) {
			return;
		}

		onShowCargoTransferDialog({
			src: fleet,
			dest: fleet.getCargoTransferTarget($universe)
		});
	}
</script>

{#if fleet.spec}
	<CommandTile title="Fuel & Cargo">
		<div class="flex justify-between my-1">
			<div class="w-12 text-tile-item-title">Fuel</div>
			<div class="ml-1 h-full w-full">
				<!-- TODO: add fuel transfer -->
				<FuelBar value={fleet.fuel} capacity={fleet.spec.shipDesignSpec?.fuelCapacity ?? 0} />
			</div>
		</div>

		<div class="flex justify-between my-1">
			<div class="w-12 text-tile-item-title">Cargo</div>
			<div class="ml-1 h-full w-full">
				<CargoBar
					onPointerDown={transfer}
					canTransferCargo={canTransferCargo(fleet)}
					value={fleet.cargo}
					capacity={fleet.spec.shipDesignSpec?.cargoCapacity}
				/>
			</div>
		</div>
		<div class="flex justify-between">
			<div class="text-tile-item-title text-ironium">Ironium</div>
			<div>{fleet.cargo.ironium}kT</div>
		</div>
		<div class="flex justify-between">
			<div class="text-tile-item-title text-boranium">Boranium</div>
			<div>{fleet.cargo.boranium}kT</div>
		</div>
		<div class="flex justify-between">
			<div class="text-tile-item-title text-germanium">Germanium</div>
			<div>{fleet.cargo.germanium}kT</div>
		</div>
		<div class="flex justify-between">
			<div class="text-tile-item-title text-colonists">Colonists</div>
			<div>{fleet.cargo.colonists}kT</div>
		</div>
	</CommandTile>
{/if}
