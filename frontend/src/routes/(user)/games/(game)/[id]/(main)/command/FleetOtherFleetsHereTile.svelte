<script lang="ts">
	import type {
		ShowCargoTransferDialogProps,
		ShowSplitFleetDialogProps
	} from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import type { AnyFleet } from '$lib/services/Universe';
	import type { Fleet } from '$lib/types/cs';
	import { canLoadCargo, type CommandedFleet } from '$lib/types/Fleet';
	import { getMapObjectName, key } from '$lib/types/MapObject';
	import { onDestroy } from 'svelte';
	import CommandTile from './CommandTile.svelte';

	const { universe, game, player, commandedFleet, commandedMapObjectKey, commandMapObject } =
		getGameContext();

	type Props = {
		fleet: CommandedFleet;
		fleetsInOrbit: AnyFleet[];
	} & ShowCargoTransferDialogProps &
		ShowSplitFleetDialogProps;

	let { fleet, fleetsInOrbit, onShowCargoTransferDialog, onShowSplitFleetDialog }: Props = $props();

	let selectedFleetKey = $state(fleetsInOrbit.length > 0 ? key(fleetsInOrbit[0]) : '');
	$effect(() => {
		if (fleetsInOrbit.length > 0 && selectedFleetKey === '') {
			selectedFleetKey = key(fleetsInOrbit[0]);
		}
	});

	let fleetsInOrbitByKey = $derived(
		fleetsInOrbit.reduce<Record<string, AnyFleet>>((acc, fleet) => {
			acc[key(fleet)] = fleet;
			return acc;
		}, {})
	);

	let selectedFleet: AnyFleet | undefined = $derived(
		selectedFleetKey !== '' ? fleetsInOrbitByKey[selectedFleetKey] : undefined
	);

	let fleetsByPlayer = $derived(
		fleetsInOrbit.reduce<Record<number, AnyFleet[]>>((acc, fleet) => {
			if (!acc[fleet.playerNum]) {
				acc[fleet.playerNum] = [];
			}
			acc[fleet.playerNum].push(fleet);
			return acc;
		}, {})
	);

	const onSelectedFleetChange = (key: string) => {
		selectedFleetKey = key;
	};

	const transfer = () => {
		if (!selectedFleet || !canLoadCargo(fleet, selectedFleet) || !onShowCargoTransferDialog) {
			return;
		}
		onShowCargoTransferDialog({ src: fleet, dest: selectedFleet });
	};

	const gotoTarget = () => {
		if (!selectedFleet || selectedFleet.playerNum !== $player.num) {
			return;
		}
		commandMapObject(selectedFleet);
	};

	const mergeTarget = () => {
		if (
			!$commandedFleet ||
			!selectedFleet ||
			!onShowSplitFleetDialog ||
			selectedFleet.playerNum !== $player.num
		) {
			return;
		}
		onShowSplitFleetDialog({ src: $commandedFleet, dest: selectedFleet as Fleet });
	};

	// reset the waypoint index every time the commanded mapobject changes
	const unsubscribe = commandedMapObjectKey.subscribe(() => {		
		selectedFleetKey = '';
	});
	onDestroy(unsubscribe);
</script>

{#if fleet}
	<CommandTile title="Other Fleets Here">
		<select
			data-type="other-fleets-here-select"
			onchange={(e) => onSelectedFleetChange(e.currentTarget.value)}
			class="select select-outline select-secondary select-sm py-0 text-sm"
		>
			{#each fleetsByPlayer[$player.num].filter((f) => f.num !== fleet.num) as f}
				<option
					style={f.playerNum !== $player.num
						? `color: ${$universe.getPlayerColor(f.playerNum)};`
						: ''}
					value={key(f)}
				>
					{getMapObjectName(f)}
				</option>
			{/each}
			{#each $game.players as p}
				{#if p.num !== $player.num && p.num in fleetsByPlayer}
					<optgroup
						label={$universe.getPlayerName(p.num)}
						style={`color: ${$universe.getPlayerColor(p.num)};`}
					>
						{#each fleetsByPlayer[p.num] as f}
							<option value={key(f)}>
								{getMapObjectName(f)}
							</option>
						{/each}
					</optgroup>
				{/if}
			{/each}
		</select>

		{#if selectedFleet}
			<div class="flex justify-between my-1 btn-group">
				<div class="tooltip" data-tip="goto fleet">
					<button
						onclick={gotoTarget}
						disabled={!selectedFleet || selectedFleet.playerNum !== $player.num}
						class="btn btn-outline btn-sm normal-case btn-secondary p-2"
						title="goto">Goto</button
					>
				</div>
				<div class="tooltip" data-tip="merge fleet">
					<button
						onclick={mergeTarget}
						disabled={!selectedFleet || selectedFleet.playerNum !== $player.num}
						class="btn btn-outline btn-sm normal-case btn-secondary p-2"
						title="goto"
						>Merge
					</button>
				</div>
				<div class="tooltip" data-tip="transfer cargo">
					<button
						onclick={transfer}
						disabled={!selectedFleet || !canLoadCargo(fleet, selectedFleet)}
						class="btn btn-outline btn-sm normal-case btn-secondary p-2"
						title="goto"
						>Transfer
					</button>
				</div>
			</div>
		{/if}
	</CommandTile>
{/if}
