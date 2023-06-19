<script lang="ts">
	import CargoBar from '$lib/components/game/CargoBar.svelte';
	import FuelBar from '$lib/components/game/FuelBar.svelte';
	import { EventManager } from '$lib/EventManager';
	import {
		commandedMapObjectName,
		commandMapObject,
		getMyMapObjectsByPosition
	} from '$lib/services/Context';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType } from '$lib/types/MapObject';
	import type { CommandedPlanet } from '$lib/types/Planet';
	import { ArrowTopRightOnSquare } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { onMount } from 'svelte';
	import CommandTile from './CommandTile.svelte';
	
	export let planet: CommandedPlanet;
	let fleetsInOrbit: Fleet[];
	let selectedFleet: Fleet | undefined;
	let selectedFleetIndex = 0;

	onMount(() => {
		const unsubscribe = EventManager.subscribeCargoTransferredEvent((mo) => {
			if (selectedFleet == mo) {
				// trigger a reaction
				selectedFleet.cargo = (mo as Fleet).cargo;
			}
		});

		return () => unsubscribe();
	});

	commandedMapObjectName.subscribe(() => (selectedFleetIndex = 0));

	$: {
		fleetsInOrbit = getMyMapObjectsByPosition(planet).filter(
			(mo) => mo.type == MapObjectType.Fleet
		) as Fleet[];
		if (fleetsInOrbit.length > 0) {
			selectedFleet = fleetsInOrbit[selectedFleetIndex];
		}
	}

	const onSelectedFleetChange = (index: number) => {
		selectedFleet = fleetsInOrbit[index];
		selectedFleetIndex = index;
	};

	const transfer = () => {
		if (selectedFleet) {
			EventManager.publishCargoTransferDialogRequestedEvent(selectedFleet, planet);
		}
	};

	const gotoTarget = () => {
		if (selectedFleet) {
			commandMapObject(selectedFleet);
		}
	};
</script>

<CommandTile title="Fleets In Orbit">
	<select
		on:change={(e) => onSelectedFleetChange(parseInt(e.currentTarget.value))}
		class="select select-outline select-secondary select-sm py-0 text-sm"
	>
		{#each fleetsInOrbit as fleet, index}
			<option value={index}>{fleet.name}</option>>
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
					on:cargo-transfer={transfer}
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
