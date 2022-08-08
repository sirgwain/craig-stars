<script lang="ts">
	import { EventManager } from '$lib/EventManager';
	import { commandedFleet,commandMapObject,myMapObjectsByPosition } from '$lib/services/Context';
	import { PlanetService } from '$lib/services/PlanetService';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType,positionKey } from '$lib/types/MapObject';
	import CommandTile from './CommandTile.svelte';

	const planetService = new PlanetService();

	let fleetsInOrbit: Fleet[];
	let selectedFleet: Fleet | undefined;

	$: {
		if ($commandedFleet && $myMapObjectsByPosition) {
			const mapObjectsByPosition = $myMapObjectsByPosition[positionKey($commandedFleet)];
			fleetsInOrbit = mapObjectsByPosition.filter(
				(mo) => mo.type == MapObjectType.Fleet && mo != $commandedFleet
			) as Fleet[];
			if (fleetsInOrbit.length > 0) {
				selectedFleet = fleetsInOrbit[0];
			}
		}
	}

	const onSelectedFleetChange = (index: number) => {
		selectedFleet = fleetsInOrbit[index];
	};

	const transfer = () => {
		if ($commandedFleet && selectedFleet) {
			EventManager.publishCargoTransferDialogRequestedEvent($commandedFleet, selectedFleet);
		}
	};

	const gotoTarget = () => {
		if (selectedFleet) {
			commandMapObject(selectedFleet);
		}
	};

	const mergeTarget = () => {
		if (selectedFleet) {
			// EventManager.publishFleetMergeDialogRequestedEvent($commandedFleet, selectedFleet);
		}
	};

</script>

{#if $commandedFleet}
	<CommandTile title="Other Fleets Here">
		<select
			on:change={(e) => onSelectedFleetChange(parseInt(e.currentTarget.value))}
			class="select select-outline select-secondary select-sm py-0 text-sm"
		>
			{#each fleetsInOrbit as fleet, index}
				<option value={index}>{fleet.name}</option>>
			{/each}
		</select>

		{#if selectedFleet}
			<div class="flex justify-between my-1 btn-group">
				<div class="tooltip" data-tip="goto fleet">
					<button
						on:click={gotoTarget}
						disabled={!selectedFleet}
						class="btn btn-outline btn-sm normal-case btn-secondary p-2"
						title="goto"
						>Goto</button
					>
				</div>
				<div class="tooltip" data-tip="merge fleet">
					<button
						on:click={mergeTarget}
						disabled={!selectedFleet}
						class="btn btn-outline btn-sm normal-case btn-secondary p-2"
						title="goto"
						>Merge
					</button>
				</div>
				<div class="tooltip" data-tip="transfer cargo">
					<button
						on:click={transfer}
						disabled={!selectedFleet}
						class="btn btn-outline btn-sm normal-case btn-secondary p-2"
						title="goto"
						>Transfer
					</button>
				</div>
			</div>
		{/if}
	</CommandTile>
{/if}
