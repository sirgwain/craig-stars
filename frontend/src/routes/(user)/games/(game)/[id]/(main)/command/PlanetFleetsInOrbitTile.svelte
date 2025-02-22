<script lang="ts">
	import CargoBar from '$lib/components/game/CargoBar.svelte';
	import FuelBar from '$lib/components/game/FuelBar.svelte';
	import type { ShowCargoTransferDialogProps } from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import { canTransferCargo, CommandedFleet } from '$lib/types/Fleet';
	import { type Fleet } from '$lib/types/cs';
	import { getMapObjectName } from '$lib/types/MapObject';
	import type { CommandedPlanet } from '$lib/types/Planet';
	import { ArrowTopRightOnSquare } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { onMount } from 'svelte';
	import CommandTile from './CommandTile.svelte';

	const { universe, commandedMapObjectKey, commandMapObject } = getGameContext();

	type Props = {
		planet: CommandedPlanet;
		fleetsInOrbit: Fleet[];
	} & ShowCargoTransferDialogProps;

	let { planet, fleetsInOrbit, onShowCargoTransferDialog }: Props = $props();
	let selectedFleetIndex = $state(0);

	let selectedFleet: Fleet | undefined = $derived.by(() => {
		if (fleetsInOrbit.length > 0) {
			return fleetsInOrbit[selectedFleetIndex];
		} else {
			return undefined;
		}
	});

	const onSelectedFleetChange = (index: number) => {
		selectedFleetIndex = index;
	};

	const transfer = () => {
		if (!selectedFleet || !onShowCargoTransferDialog) {
			return;
		}
		const commandedFleet = new CommandedFleet(selectedFleet);
		onShowCargoTransferDialog({ src: commandedFleet, dest: planet });
	};

	const gotoTarget = () => {
		if (selectedFleet) {
			commandMapObject(selectedFleet);
		}
	};

	onMount(() => {
		return commandedMapObjectKey.subscribe(() => (selectedFleetIndex = 0));
	});
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
					onPointerDown={transfer}
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
