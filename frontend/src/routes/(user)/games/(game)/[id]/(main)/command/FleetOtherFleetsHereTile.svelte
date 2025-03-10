<script lang="ts">
	import type {
		ShowCargoTransferDialogProps,
		ShowSplitFleetDialogProps
	} from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import { type CommandedFleet } from '$lib/types/Fleet';
	import { type Fleet } from '$lib/types/cs';
	import { getMapObjectName } from '$lib/types/MapObject';
	import { onDestroy } from 'svelte';
	import CommandTile from './CommandTile.svelte';

	const { commandedFleet, commandedMapObjectKey, commandMapObject } = getGameContext();

	type Props = {
		fleet: CommandedFleet;
		fleetsInOrbit: Fleet[];
	} & ShowCargoTransferDialogProps &
		ShowSplitFleetDialogProps;

	let { fleet, fleetsInOrbit, onShowCargoTransferDialog, onShowSplitFleetDialog }: Props = $props();

	let selectedFleetIndex = $state(0);
	let selectedFleet: Fleet | undefined = $derived(
		fleetsInOrbit.length > 0 ? fleetsInOrbit[selectedFleetIndex] : undefined
	);

	const onSelectedFleetChange = (index: number) => {
		selectedFleetIndex = index;
	};

	const transfer = () => {
		if (!selectedFleet || !onShowCargoTransferDialog) {
			return;
		}
		onShowCargoTransferDialog({ src: fleet, dest: selectedFleet });
	};

	const gotoTarget = () => {
		if (!selectedFleet) {
			return;
		}
		commandMapObject(selectedFleet);
	};

	const mergeTarget = () => {
		if (!$commandedFleet || !selectedFleet || !onShowSplitFleetDialog) {
			return;
		}
		onShowSplitFleetDialog({ src: $commandedFleet, dest: selectedFleet });
	};

	// reset the waypoint index every time the commanded mapobject changes
	const unsubscribe = commandedMapObjectKey.subscribe(() => (selectedFleetIndex = 0));
	onDestroy(unsubscribe);
</script>

{#if fleet}
	<CommandTile title="Other Fleets Here">
		<select
			data-type="other-fleets-here-select"
			onchange={(e) => onSelectedFleetChange(parseInt(e.currentTarget.value))}
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
						onclick={gotoTarget}
						disabled={!selectedFleet}
						class="btn btn-outline btn-sm normal-case btn-secondary p-2"
						title="goto">Goto</button
					>
				</div>
				<div class="tooltip" data-tip="merge fleet">
					<button
						onclick={mergeTarget}
						disabled={!selectedFleet}
						class="btn btn-outline btn-sm normal-case btn-secondary p-2"
						title="goto"
						>Merge
					</button>
				</div>
				<div class="tooltip" data-tip="transfer cargo">
					<button
						onclick={transfer}
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
