<script lang="ts">
	import CargoBar from '$lib/components/game/CargoBar.svelte';
	import MineralBar from '$lib/components/game/MineralBar.svelte';
	import { add } from '$lib/types/Cargo';
	import { CargoTransferRequest } from '$lib/types/CargoTransferRequest.svelte';

	type Props = {
		transferAmount: CargoTransferRequest;
		cargo: CargoTransferRequest;
		cargoCapacity: number;
		fuelCapacity: number;
		allowFuelTransfers?: boolean;
		allowColonistTransfers?: boolean;
		onTransferFuel?: (amount: number) => number;
		onTransferIronium?: (amount: number) => number;
		onTransferBoranium?: (amount: number) => number;
		onTransferGermanium?: (amount: number) => number;
		onTransferColonists?: (amount: number) => number;
	};

	let {
		transferAmount = new CargoTransferRequest(),
		cargo = new CargoTransferRequest(),
		cargoCapacity = 0,
		fuelCapacity = 0,
		allowFuelTransfers = false,
		allowColonistTransfers = false,
		onTransferFuel: onTransferFuel,
		onTransferIronium: onFransferIronium,
		onTransferBoranium: onTransferBoranium,
		onTransferGermanium: onTransferGermanium,
		onTransferColonists: onTransferColonists
	}: Props = $props();
</script>

<div class="sm:grid sm:grid-cols-label-value">
	<div class="sm:text-right mr-1 h-8 select-none">Fuel</div>
	<div class="my-auto">
		<MineralBar
			value={cargo.fuel + transferAmount.fuel}
			capacity={fuelCapacity}
			color="fuel-bar"
			unit="mg"
			readonly={!allowFuelTransfers}
			onValueChanged={(value) => onTransferFuel?.(value - (cargo.fuel + transferAmount.fuel))}
		/>
	</div>

	<div class="sm:text-right mr-1 h-8 select-none">Cargo Hold</div>
	<div class="my-auto">
		<CargoBar value={add(cargo.cargo(), transferAmount.cargo())} capacity={cargoCapacity} />
	</div>

	<div class="col-span-2 mt-10 sm:mt-5"></div>

	<div class="sm:text-right mr-1 h-8 select-none">Ironium</div>
	<div class="my-auto">
		<MineralBar
			value={cargo.ironium + transferAmount.ironium}
			capacity={cargoCapacity}
			color="ironium-bar"
			onValueChanged={(value) =>
				onFransferIronium?.(value - (cargo.ironium + transferAmount.ironium))}
		/>
	</div>
	<div class="sm:text-right mr-1 h-8 select-none">Boranium</div>
	<div class="my-auto">
		<MineralBar
			value={cargo.boranium + transferAmount.boranium}
			capacity={cargoCapacity}
			color="boranium-bar"
			onValueChanged={(value) =>
				onTransferBoranium?.(value - (cargo.boranium + transferAmount.boranium))}
		/>
	</div>
	<div class="sm:text-right mr-1 h-8 select-none">Germanium</div>
	<div class="my-auto">
		<MineralBar
			value={cargo.germanium + transferAmount.germanium}
			capacity={cargoCapacity}
			color="germanium-bar"
			onValueChanged={(value) =>
				onTransferGermanium?.(value - (cargo.germanium + transferAmount.germanium))}
		/>
	</div>

	<div class="sm:text-right mr-1 h-8 select-none">Colonists</div>
	<div class="my-auto">
		<MineralBar
			value={cargo.colonists + transferAmount.colonists}
			capacity={cargoCapacity}
			color="colonists-bar"
			readonly={!allowColonistTransfers}
			onValueChanged={(value) =>
				onTransferColonists?.(value - (cargo.colonists + transferAmount.colonists))}
		/>
	</div>
	{#if cargo.colonists + transferAmount.colonists}
		<div class="my-auto col-span-2 ml-auto pr-1">
			<span class="italic text-sm"
				>{((cargo.colonists + transferAmount.colonists) * 100).toLocaleString()} colonists</span
			>
		</div>
	{/if}
</div>
