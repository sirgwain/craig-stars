<script lang="ts">
	import { add, emptyCargo, type Cargo } from '$lib/types/Cargo';
	import { createEventDispatcher } from 'svelte';
	import CargoBar from '$lib/components/game/CargoBar.svelte';
	import FuelBar from '$lib/components/game/FuelBar.svelte';
	import MineralBar from '$lib/components/game/MineralBar.svelte';

	export let cargo: Cargo = emptyCargo();
	export let fuel: number = 0;
	export let transferAmount: Cargo = emptyCargo();
	export let fuelTransferAmount: number = 0;
	export let cargoCapacity: number = 0;
	export let fuelCapacity: number = 0;

	const dispatch = createEventDispatcher();
</script>

<div class="grid grid-cols-2">
	<div class="text-right mr-1 h-8">Fuel</div>
	<div class="my-auto">
		<FuelBar value={fuel + fuelTransferAmount} capacity={fuelCapacity} />
	</div>

	<div class="text-right mr-1 h-8">Cargo Hold</div>
	<div class="my-auto">
		<CargoBar value={add(cargo, transferAmount)} capacity={cargoCapacity} />
	</div>

	<div class="col-span-2 mt-5" />

	<div class="text-right mr-1 h-8">Ironium</div>
	<div class="my-auto">
		<MineralBar value={(cargo.ironium ?? 0) + (transferAmount.ironium ?? 0)} capacity={cargoCapacity} color="ironium-bar" />
	</div>
	<div class="text-right mr-1 h-8">Boranium</div>
	<div class="my-auto">
		<MineralBar value={(cargo.boranium ?? 0) + (transferAmount.boranium ?? 0)} capacity={cargoCapacity} color="boranium-bar" />
	</div>
	<div class="text-right mr-1 h-8">Germanium</div>
	<div class="my-auto">
		<MineralBar value={(cargo.germanium ?? 0) + (transferAmount.germanium ?? 0)} capacity={cargoCapacity} color="germanium-bar" />
	</div>

	<div class="text-right mr-1 h-8">Colonists</div>
	<div class="my-auto">
		<MineralBar value={(cargo.colonists ?? 0) + (transferAmount.colonists ?? 0)} capacity={cargoCapacity} color="colonists-bar" />
	</div>
</div>
