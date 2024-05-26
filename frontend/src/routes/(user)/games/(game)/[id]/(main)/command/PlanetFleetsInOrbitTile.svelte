<script lang="ts">
	import CargoBar from '$lib/components/game/CargoBar.svelte';
	import FuelBar from '$lib/components/game/FuelBar.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { canTransferCargo, CommandedFleet, type Fleet } from '$lib/types/Fleet';
	import type { CommandedPlanet } from '$lib/types/Planet';
	import { ArrowTopRightOnSquare } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher, onDestroy } from 'svelte';
	import type { CargoTransferDialogEvent } from '../../dialogs/cargo/CargoTranfserDialog.svelte';
	import CommandTile from './CommandTile.svelte';

	const dispatch = createEventDispatcher<CargoTransferDialogEvent>();

	const { universe, commandedMapObjectKey, commandMapObject } = getGameContext();

	export let planet: CommandedPlanet;
	export let fleetsInOrbit: Fleet[];
	let selectedFleet: Fleet | undefined;
	let selectedFleetIndex = 0;

	$: {
		if (fleetsInOrbit.length > 0) {
			selectedFleet = fleetsInOrbit[selectedFleetIndex];
		} else {
			selectedFleet = undefined;
		}
	}

	const onSelectedFleetChange = (index: number) => {
		selectedFleet = fleetsInOrbit[index];
		selectedFleetIndex = index;
	};

	const transfer = () => {
		if (selectedFleet) {
			const commandedFleet = new CommandedFleet(selectedFleet);
			dispatch('cargo-transfer-dialog', { src: commandedFleet, dest: planet });
		}
	};

	const gotoTarget = () => {
		if (selectedFleet) {
			commandMapObject(selectedFleet);
		}
	};

	const unsubscribe = commandedMapObjectKey.subscribe(() => (selectedFleetIndex = 0));
	onDestroy(unsubscribe);
</script>

<CommandTile title="Fleets In Orbit">
	<select
		on:change={(e) => onSelectedFleetChange(parseInt(e.currentTarget.value))}
		class="select select-outline select-secondary select-sm py-0 text-sm"
	>
		{#each fleetsInOrbit as fleet, index}
			<option value={index}>{fleet.name}</option>
		{/each}
	</select>

	{#if selectedFleet && selectedFleet.spec}
		<div class="flex justify-between my-1">
			<div class="w-12">Fuel</div>
			<div class="ml-1 h-full w-full">
				<FuelBar value={selectedFleet.fuel} capacity={selectedFleet.spec.fuelCapacity} />
			</div>
		</div>

		<div class="flex justify-between my-1">
			<div class="w-12">Cargo</div>
			<div class="ml-1 h-full w-full">
				<CargoBar
					on:cargo-transfer-dialog={transfer}
					canTransferCargo={canTransferCargo(selectedFleet, $universe)}
					value={selectedFleet.cargo}
					capacity={selectedFleet.spec.cargoCapacity}
				/>
			</div>
		</div>

		<div class="flex justify-between my-1">
			<div class="tooltip" data-tip="command fleet">
				<button
					on:click={gotoTarget}
					disabled={!selectedFleet}
					class="btn btn-outline btn-sm normal-case btn-secondary"
					title="goto"
					>Goto<Icon
						src={ArrowTopRightOnSquare}
						size="16"
						class="hover:stroke-accent inline"
					/></button
				>
			</div>
		</div>
	{/if}
</CommandTile>
