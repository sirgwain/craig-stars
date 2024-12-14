<script lang="ts">
	import CargoBar from '$lib/components/game/CargoBar.svelte';
	import MineralBar from '$lib/components/game/MineralBar.svelte';
	import { add } from '$lib/types/Cargo';
	import {
		type CargoTransferRequest,
		newCargoTransferRequest
	} from '$lib/types/CargoTransferRequest';

	type Props = {
		transferAmount: CargoTransferRequest;
		cargo: CargoTransferRequest;
		cargoCapacity: number;
		fuelCapacity: number;
		allowFuelTransfers?: boolean;
		onTransferFuel?: (amount: number) => void;
		onTransferIronium?: (amount: number) => void;
		onTransferBoranium?: (amount: number) => void;
		onTransferGermanium?: (amount: number) => void;
		onTransferColonists?: (amount: number) => void;
	};

	let {
		transferAmount = newCargoTransferRequest(),
		cargo = newCargoTransferRequest(),
		cargoCapacity = 0,
		fuelCapacity = 0,
		allowFuelTransfers = false,
		onTransferFuel: onTransferVuel,
		onTransferIronium: onFransferIronium,
		onTransferBoranium: onTransferBoranium,
		onTransferGermanium: onTransferGermanium,
		onTransferColonists: onTransferColonists
	}: Props = $props();
</script>

<div class="sm:grid sm:grid-cols-label-value">
	<div class="sm:text-right mr-1 h-8">Fuel</div>
	<div class="my-auto">
		<MineralBar
			value={cargo.fuel + transferAmount.fuel}
			capacity={fuelCapacity}
			color="fuel-bar"
			unit="mg"
			readonly={!allowFuelTransfers}
			onValueChanged={(value) => onTransferVuel?.(value - (cargo.fuel + transferAmount.fuel))}
		/>
	</div>

	<div class="sm:text-right mr-1 h-8">Cargo Hold</div>
	<div class="my-auto">
		<CargoBar value={add(cargo, transferAmount)} capacity={cargoCapacity} />
	</div>

	<div class="col-span-2 mt-10 sm:mt-5"></div>

	<div class="sm:text-right mr-1 h-8">Ironium</div>
	<div class="my-auto">
		<MineralBar
			value={cargo.ironium + transferAmount.ironium}
			capacity={cargoCapacity}
			color="ironium-bar"
			onValueChanged={(value) =>
				onFransferIronium?.(value - (cargo.ironium + transferAmount.ironium))}
		/>
	</div>
	<div class="sm:text-right mr-1 h-8">Boranium</div>
	<div class="my-auto">
		<MineralBar
			value={cargo.boranium + transferAmount.boranium}
			capacity={cargoCapacity}
			color="boranium-bar"
			onValueChanged={(value) =>
				onTransferBoranium?.(value - (cargo.boranium + transferAmount.boranium))}
		/>
	</div>
	<div class="sm:text-right mr-1 h-8">Germanium</div>
	<div class="my-auto">
		<MineralBar
			value={cargo.germanium + transferAmount.germanium}
			capacity={cargoCapacity}
			color="germanium-bar"
			onValueChanged={(value) =>
				onTransferGermanium?.(value - (cargo.germanium + transferAmount.germanium))}
		/>
	</div>

	<div class="sm:text-right mr-1 h-8">Colonists</div>
	<div class="my-auto">
		<MineralBar
			value={cargo.colonists + transferAmount.colonists}
			capacity={cargoCapacity}
			color="colonists-bar"
			onValueChanged={(value) =>
				onTransferColonists?.(value - (cargo.colonists + transferAmount.colonists))}
		/>
	</div>
</div>
