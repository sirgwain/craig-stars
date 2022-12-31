<script lang="ts">
	import CargoBar from '$lib/components/game/CargoBar.svelte';
	import FuelBar from '$lib/components/game/FuelBar.svelte';
	import { EventManager } from '$lib/EventManager';
	import { findMyPlanetByNum } from '$lib/services/Context';
	import type { Fleet } from '$lib/types/Fleet';
	import { onMount } from 'svelte';
	import CommandTile from './CommandTile.svelte';

	export let fleet: Fleet;

	onMount(() => {
		const unsubscribe = EventManager.subscribeCargoTransferredEvent((mo) => {
			if (fleet == mo) {
				// trigger a reaction
				fleet.cargo = (mo as Fleet).cargo;
			}
		});

		return () => unsubscribe();
	});

	const transfer = () => {
		if (fleet.orbitingPlanetNum) {
			const planet = findMyPlanetByNum(fleet.orbitingPlanetNum);
			EventManager.publishCargoTransferDialogRequestedEvent(fleet, planet);
		} else {
			EventManager.publishCargoTransferDialogRequestedEvent(fleet);
		}
	};
</script>

{#if fleet?.spec}
	<CommandTile title="Fuel & Cargo">
		<div class="flex justify-between my-1">
			<div class="w-12">Fuel</div>
			<div class="ml-1 h-full w-full">
				<FuelBar value={fleet.fuel} capacity={fleet.spec.fuelCapacity} />
			</div>
		</div>

		<div class="flex justify-between my-1">
			<div class="w-12">Cargo</div>
			<div class="ml-1 h-full w-full">
				<CargoBar
					on:cargo-transfer={() => transfer()}
					value={fleet.cargo}
					capacity={fleet.spec.cargoCapacity}
				/>
			</div>
		</div>
		<div class="flex justify-between">
			<div class="text-ironium">Ironium</div>
			<div>{fleet.cargo?.ironium ?? 0}kT</div>
		</div>
		<div class="flex justify-between">
			<div class="text-boranium">Boranium</div>
			<div>{fleet.cargo?.boranium ?? 0}kT</div>
		</div>
		<div class="flex justify-between">
			<div class="text-germanium">Germanium</div>
			<div>{fleet.cargo?.germanium ?? 0}kT</div>
		</div>
		<div class="flex justify-between">
			<div class="text-colonists">Colonists</div>
			<div>{fleet.cargo?.colonists ?? 0}kT</div>
		</div>
	</CommandTile>
{/if}
