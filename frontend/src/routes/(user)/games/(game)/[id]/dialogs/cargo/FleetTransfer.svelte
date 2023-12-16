<script lang="ts" context="module">
	export type FleetTransferEvent = {
		'transfer-fuel': number;
		'transfer-ironium': number;
		'transfer-boranium': number;
		'transfer-germanium': number;
		'transfer-colonists': number;
	};
</script>

<script lang="ts">
	import CargoBar from '$lib/components/game/CargoBar.svelte';
	import MineralBar from '$lib/components/game/MineralBar.svelte';
	import { CargoTransferRequest, add } from '$lib/types/Cargo';
	import { createEventDispatcher } from 'svelte';

	export let transferAmount = new CargoTransferRequest();
	export let cargo = new CargoTransferRequest();
	export let cargoCapacity: number = 0;
	export let fuelCapacity: number = 0;
	export let allowFuelTransfers = false;

	const dispatch = createEventDispatcher<FleetTransferEvent>();
</script>

<div class="grid grid-cols-2">
	<div class="text-right mr-1 h-8">Fuel</div>
	<div class="my-auto">
		<MineralBar
			value={cargo.fuel + transferAmount.fuel}
			capacity={fuelCapacity}
			color="fuel-bar"
			unit="mg"
			readonly={!allowFuelTransfers}
			on:valuechanged={(e) =>
				dispatch('transfer-fuel', e.detail - (cargo.fuel + transferAmount.fuel))}
		/>
	</div>

	<div class="text-right mr-1 h-8">Cargo Hold</div>
	<div class="my-auto">
		<CargoBar value={add(cargo, transferAmount)} capacity={cargoCapacity} />
	</div>

	<div class="col-span-2 mt-5" />

	<div class="text-right mr-1 h-8">Ironium</div>
	<div class="my-auto">
		<MineralBar
			value={cargo.ironium + transferAmount.ironium}
			capacity={cargoCapacity}
			color="ironium-bar"
			on:valuechanged={(e) =>
				dispatch('transfer-ironium', e.detail - (cargo.ironium + transferAmount.ironium))}
		/>
	</div>
	<div class="text-right mr-1 h-8">Boranium</div>
	<div class="my-auto">
		<MineralBar
			value={cargo.boranium + transferAmount.boranium}
			capacity={cargoCapacity}
			color="boranium-bar"
			on:valuechanged={(e) =>
				dispatch('transfer-boranium', e.detail - (cargo.boranium + transferAmount.boranium))}
		/>
	</div>
	<div class="text-right mr-1 h-8">Germanium</div>
	<div class="my-auto">
		<MineralBar
			value={cargo.germanium + transferAmount.germanium}
			capacity={cargoCapacity}
			color="germanium-bar"
			on:valuechanged={(e) =>
				dispatch('transfer-germanium', e.detail - (cargo.germanium + transferAmount.germanium))}
		/>
	</div>

	<div class="text-right mr-1 h-8">Colonists</div>
	<div class="my-auto">
		<MineralBar
			value={cargo.colonists + transferAmount.colonists}
			capacity={cargoCapacity}
			color="colonists-bar"
			on:valuechanged={(e) =>
				dispatch('transfer-colonists', e.detail - (cargo.colonists + transferAmount.colonists))}
		/>
	</div>
</div>
