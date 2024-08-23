<script lang="ts">
	import { type CommandedFleet, type Fleet } from '$lib/types/Fleet';
	import { createEventDispatcher, onDestroy } from 'svelte';
	import CommandTile from './CommandTile.svelte';
	import type { CargoTransferDialogEvent } from '../../dialogs/cargo/CargoTranfserDialog.svelte';
	import type { SplitFleetDialogEvent } from '../../dialogs/split/SplitFleetDialog.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { getMapObjectName } from '$lib/types/MapObject';

	const dispatch = createEventDispatcher<SplitFleetDialogEvent & CargoTransferDialogEvent>();

	const { commandedFleet, commandedMapObjectKey, commandMapObject } = getGameContext();

	export let fleet: CommandedFleet;
	export let fleetsInOrbit: Fleet[];

	let selectedFleet: Fleet | undefined;
	let selectedFleetIndex = 0;

	$: {
		if (fleetsInOrbit.length > 0) {
			selectedFleet = fleetsInOrbit[selectedFleetIndex];
		}
	}

	const onSelectedFleetChange = (index: number) => {
		selectedFleetIndex = index;
		selectedFleet = fleetsInOrbit[selectedFleetIndex];
	};

	const transfer = () => {
		if (selectedFleet) {
			dispatch('cargo-transfer-dialog', { src: fleet, dest: selectedFleet });
		}
	};

	const gotoTarget = () => {
		if (selectedFleet) {
			commandMapObject(selectedFleet);
		}
	};

	const mergeTarget = () => {
		if ($commandedFleet && selectedFleet) {
			dispatch('split-fleet-dialog', { src: $commandedFleet, dest: selectedFleet });
		}
	};

	// reset the waypoint index every time the commanded mapobject changes
	const unsubscribe = commandedMapObjectKey.subscribe(() => (selectedFleetIndex = 0));
	onDestroy(unsubscribe);
</script>

{#if fleet}
	<CommandTile title="Other Fleets Here">
		<select
			on:change={(e) => onSelectedFleetChange(parseInt(e.currentTarget.value))}
			class="select select-outline select-secondary select-sm py-0 text-sm"
		>
			{#each fleetsInOrbit as fleet, index}
				<option value={index}>{getMapObjectName(fleet)}</option>
			{/each}
		</select>

		{#if selectedFleet}
			<div class="flex justify-between my-1 btn-group">
				<div class="tooltip" data-tip="goto fleet">
					<button
						on:click={gotoTarget}
						disabled={!selectedFleet}
						class="btn btn-outline btn-sm normal-case btn-secondary p-2"
						title="goto">Goto</button
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
