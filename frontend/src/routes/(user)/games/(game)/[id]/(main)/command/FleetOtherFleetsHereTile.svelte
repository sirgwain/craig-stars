<script lang="ts">
	import type {
		ShowCargoTransferDialogProps,
		ShowSplitFleetDialogProps
	} from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import type { CargoDest } from '$lib/types/CargoTransferRequest.svelte';
	import { MapObjectTypeFleet, type Fleet } from '$lib/types/cs';
	import { canLoadCargo, type CommandedFleet } from '$lib/types/Fleet';
	import { commandable, getMapObjectName, key } from '$lib/types/MapObject';
	import { onDestroy } from 'svelte';
	import CommandTile from './CommandTile.svelte';

	const { universe, game, player, commandedFleet, commandedMapObjectKey, commandMapObject } =
		getGameContext();

	type Props = {
		fleet: CommandedFleet;
		cargoDestsInOrbit: CargoDest[];
	} & ShowCargoTransferDialogProps &
		ShowSplitFleetDialogProps;

	let { fleet, cargoDestsInOrbit, onShowCargoTransferDialog, onShowSplitFleetDialog }: Props =
		$props();

	let selectedMapObjectKey = $state(cargoDestsInOrbit.length > 0 ? key(cargoDestsInOrbit[0]) : '');
	$effect(() => {
		if (cargoDestsInOrbit.length > 1 && selectedMapObjectKey === '') {
			selectedMapObjectKey = key(cargoDestsInOrbit.find((f) => key(f) !== key(fleet)));
		}
	});

	let mapObjectsInOrbitByKey = $derived(
		cargoDestsInOrbit.reduce<Record<string, CargoDest>>((acc, fleet) => {
			acc[key(fleet)] = fleet;
			return acc;
		}, {})
	);

	let selectedMapObject: CargoDest | undefined = $derived(
		selectedMapObjectKey !== '' ? mapObjectsInOrbitByKey[selectedMapObjectKey] : undefined
	);

	let cargoDestsByPlayer = $derived(
		cargoDestsInOrbit.reduce<Record<number, CargoDest[]>>((acc, mo) => {
			if (!mo) {
				return acc;
			}
			if (!acc[mo.playerNum]) {
				acc[mo.playerNum] = [];
			}
			acc[mo.playerNum].push(mo);
			return acc;
		}, {})
	);

	const onSelectedFleetChange = (key: string) => {
		selectedMapObjectKey = key;
	};

	const transfer = () => {
		if (
			!selectedMapObject ||
			!canLoadCargo(fleet, selectedMapObject) ||
			!onShowCargoTransferDialog
		) {
			return;
		}
		onShowCargoTransferDialog({ src: fleet, dest: selectedMapObject });
	};

	const gotoTarget = () => {
		if (!selectedMapObject || !commandable($player.num, selectedMapObject)) {
			return;
		}
		commandMapObject(selectedMapObject);
	};

	const mergeTarget = () => {
		if (
			!$commandedFleet ||
			!selectedMapObject ||
			selectedMapObject.type !== MapObjectTypeFleet ||
			!onShowSplitFleetDialog ||
			selectedMapObject.playerNum !== $player.num
		) {
			return;
		}
		onShowSplitFleetDialog({ src: $commandedFleet, dest: selectedMapObject as Fleet });
	};

	// reset the waypoint index every time the commanded mapobject changes
	const unsubscribe = commandedMapObjectKey.subscribe(() => {
		selectedMapObjectKey = '';
	});
	onDestroy(unsubscribe);
</script>

{#if fleet}
	<CommandTile title="Other Entities Here">
		<select
			data-type="other-fleets-here-select"
			onchange={(e) => onSelectedFleetChange(e.currentTarget.value)}
			class="select select-outline select-secondary select-sm py-0 text-sm"
		>
			{#each cargoDestsByPlayer[$player.num]?.filter((f) => f && key(f) !== key(fleet)) as f}
				<option
					style={f?.playerNum !== $player.num
						? `color: ${$universe.getPlayerColor(f?.playerNum)};`
						: ''}
					value={key(f)}
				>
					{getMapObjectName(f)}
				</option>
			{/each}
			{#each $game.players as p}
				{#if p.num !== $player.num && p.num in cargoDestsByPlayer}
					<optgroup
						label={$universe.getPlayerName(p.num)}
						style={`color: ${$universe.getPlayerColor(p.num)};`}
					>
						{#each cargoDestsByPlayer[p.num] as f}
							<option value={key(f)}>
								{getMapObjectName(f)}
							</option>
						{/each}
					</optgroup>
				{/if}
			{/each}
		</select>

		{#if selectedMapObject}
			<div class="flex justify-between my-1 btn-group">
				<div class="tooltip" data-tip="goto fleet">
					<button
						onclick={gotoTarget}
						disabled={!selectedMapObject || !commandable($player.num, selectedMapObject)}
						class="btn btn-outline btn-sm normal-case btn-secondary p-2"
						title="goto">Goto</button
					>
				</div>
				<div class="tooltip" data-tip="merge fleet">
					<button
						onclick={mergeTarget}
						disabled={!selectedMapObject ||
							selectedMapObject.type !== MapObjectTypeFleet ||
							selectedMapObject.playerNum !== $player.num}
						class="btn btn-outline btn-sm normal-case btn-secondary p-2"
						title="goto"
						>Merge
					</button>
				</div>
				<div class="tooltip" data-tip="transfer cargo">
					<button
						onclick={transfer}
						disabled={!selectedMapObject || !canLoadCargo(fleet, selectedMapObject)}
						class="btn btn-outline btn-sm normal-case btn-secondary p-2"
						title="goto"
						>Transfer
					</button>
				</div>
			</div>
		{/if}
	</CommandTile>
{/if}
