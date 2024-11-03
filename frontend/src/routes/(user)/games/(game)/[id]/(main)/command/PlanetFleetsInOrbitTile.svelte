<script lang="ts">
	import { run } from 'svelte/legacy';

	import CargoBar from '$lib/components/game/CargoBar.svelte';
	import FuelBar from '$lib/components/game/FuelBar.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { canTransferCargo, CommandedFleet, type Fleet } from '$lib/types/Fleet';
	import type { CommandedPlanet } from '$lib/types/Planet';
	import { ArrowTopRightOnSquare } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher, onDestroy } from 'svelte';
	import type { CargoTransferDialogEvent } from '../../dialogs/cargo/CargoTransferDialog.svelte';
	import CommandTile from './CommandTile.svelte';
	import { getMapObjectName } from '$lib/types/MapObject';

	const dispatch = createEventDispatcher<CargoTransferDialogEvent>();

	const { universe, commandedMapObjectKey, commandMapObject } = getGameContext();

	interface Props {
		planet: CommandedPlanet;
		fleetsInOrbit: Fleet[];
	}

	let { planet, fleetsInOrbit }: Props = $props();
	let selectedFleet: Fleet | undefined = $state();
	let selectedFleetIndex = $state(0);

	run(() => {
		if (fleetsInOrbit.length > 0) {
			selectedFleet = fleetsInOrbit[selectedFleetIndex];
		} else {
			selectedFleet = undefined;
		}
	});

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
		onchange={(e) => onSelectedFleetChange(parseInt(e.currentTarget.value))}
		class="select select-outline select-secondary select-sm py-0 text-sm"
	>
		{#each fleetsInOrbit as fleet, index}
			<option value={index}>{getMapObjectName(fleet)}</option>
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
					onclick={gotoTarget}
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
