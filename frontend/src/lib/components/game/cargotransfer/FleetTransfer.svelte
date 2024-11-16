<script lang="ts">
	import CargoBar from '$lib/components/game/CargoBar.svelte';
	import MineralBar from '$lib/components/game/MineralBar.svelte';
	import { CargoTransferRequest, add, type Cargo } from '$lib/types/Cargo';

	type Props = {
		transferAmount: CargoTransferRequest;
		cargo: CargoTransferRequest;
		cargoCapacity: number;
		fuelCapacity: number;
		allowFuelTransfers?: boolean;
		ontransferfuel?: (amount: number) => void;
		ontransferironium?: (amount: number) => void;
		ontransferboranium?: (amount: number) => void;
		ontransfergermanium?: (amount: number) => void;
		ontransfercolonists?: (amount: number) => void;
	};

	let {
		transferAmount = new CargoTransferRequest(),
		cargo = new CargoTransferRequest(),
		cargoCapacity = 0,
		fuelCapacity = 0,
		allowFuelTransfers = false,
		ontransferfuel,
		ontransferironium,
		ontransferboranium,
		ontransfergermanium,
		ontransfercolonists
	}: Props = $props();

	let cargoState = $derived(add(cargo, transferAmount));
	$inspect(transferAmount);
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
			onvaluechanged={(value) =>
				ontransferfuel && ontransferfuel(value - (cargo.fuel + transferAmount.fuel))}
		/>
	</div>

	<div class="sm:text-right mr-1 h-8">Cargo Hold</div>
	<div class="my-auto">
		<CargoBar value={cargoState} capacity={cargoCapacity} />
	</div>

	<div class="col-span-2 mt-10 sm:mt-5"></div>

	<div class="sm:text-right mr-1 h-8">Ironium</div>
	<div class="my-auto">
		<MineralBar
			value={cargo.ironium + transferAmount.ironium}
			capacity={cargoCapacity}
			color="ironium-bar"
			onvaluechanged={(value) =>
				ontransferironium && ontransferironium(value - (cargo.ironium + transferAmount.ironium))}
		/>
	</div>
	<div class="sm:text-right mr-1 h-8">Boranium</div>
	<div class="my-auto">
		<MineralBar
			value={cargo.boranium + transferAmount.boranium}
			capacity={cargoCapacity}
			color="boranium-bar"
			onvaluechanged={(value) =>
				ontransferboranium &&
				ontransferboranium(value - (cargo.boranium + transferAmount.boranium))}
		/>
	</div>
	<div class="sm:text-right mr-1 h-8">Germanium</div>
	<div class="my-auto">
		<MineralBar
			value={cargo.germanium + transferAmount.germanium}
			capacity={cargoCapacity}
			color="germanium-bar"
			onvaluechanged={(value) =>
				ontransfergermanium &&
				ontransfergermanium(value - (cargo.germanium + transferAmount.germanium))}
		/>
	</div>

	<div class="sm:text-right mr-1 h-8">Colonists</div>
	<div class="my-auto">
		<MineralBar
			value={cargo.colonists + transferAmount.colonists}
			capacity={cargoCapacity}
			color="colonists-bar"
			onvaluechanged={(value) =>
				ontransfercolonists &&
				ontransfercolonists(value - (cargo.colonists + transferAmount.colonists))}
		/>
	</div>
</div>
