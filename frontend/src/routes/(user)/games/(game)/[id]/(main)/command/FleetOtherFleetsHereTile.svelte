<script lang="ts">
	import type {
		ShowCargoTransferDialogProps,
		ShowSplitFleetDialogProps
	} from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import type { CargoDest } from '$lib/types/CargoTransferRequest';
	import { MapObjectType, type Fleet } from '$lib/types/cs-proto';
	import { canLoadFuelOrCargo, type CommandedFleet } from '$lib/types/Fleet';
	import { commandable, getMapObjectName, key } from '$lib/types/MapObject';
	import { onDestroy } from 'svelte';
	import CommandTile from './CommandTile.svelte';
	import { getDisplayColor } from '$lib/utils/colorUtils';

	const {
		universe,
		game,
		player,
		settings,
		commandedFleet,
		commandedMapObjectKey,
		commandMapObject
	} = getGameContext();

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
		cargoDestsInOrbit.reduce<Record<string, CargoDest>>((acc, mo) => {
			acc[key(mo)] = mo;
			return acc;
		}, {})
	);

	let selectedMapObject: CargoDest | undefined = $derived(
		selectedMapObjectKey !== '' ? mapObjectsInOrbitByKey[selectedMapObjectKey] : undefined
	);

	let cargoDestsByPlayer = $derived(
		cargoDestsInOrbit.reduce<Record<number, CargoDest[]>>((acc, mo) => {
			const playerNum = mo?.mapObject?.playerNum ?? 0;
			acc[playerNum] ??= [];
			acc[playerNum].push(mo);
			return acc;
		}, {})
	);

	const onSelectedFleetChange = (key: string) => {
		selectedMapObjectKey = key;
	};

	const transfer = () => {
		if (
			!selectedMapObject ||
			!canLoadFuelOrCargo(fleet, selectedMapObject) ||
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
			selectedMapObject.mapObject?.type !== MapObjectType.FLEET ||
			!onShowSplitFleetDialog ||
			selectedMapObject.mapObject.playerNum !== $player.num
		) {
			return;
		}
		// Only player-owned real Fleets can be merge targets; intel objects are excluded above.
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
			value={selectedMapObjectKey}
			onchange={(e) => onSelectedFleetChange(e.currentTarget.value)}
			class="select select-outline select-secondary select-sm py-0 text-sm"
		>
			{#if cargoDestsByPlayer[0]}
				{#each cargoDestsByPlayer[0] as mo (key(mo))}
					<option
						style={mo?.mapObject?.playerNum !== $player.num
							? `color: ${getDisplayColor(mo?.mapObject?.playerNum, $player, $universe, $settings)};`
							: ''}
						value={key(mo)}
					>
						{getMapObjectName(mo)}
					</option>
				{/each}
			{/if}
			{#each cargoDestsByPlayer[$player.num].filter((m) => m && key(m) !== key(fleet)) as m (key(m))}
				<option
					style={m?.mapObject?.playerNum !== $player.num
						? `color: ${getDisplayColor(m?.mapObject?.playerNum, $player, $universe, $settings)};`
						: ''}
					value={key(m)}
				>
					{getMapObjectName(m)}
				</option>
			{/each}
			{#each $game.players as p (p.num)}
				{#if p.num !== $player.num && p.num in cargoDestsByPlayer}
					<optgroup
						label={$universe.getPlayerName(p.num)}
						style={`color: ${getDisplayColor(p.num, $player, $universe, $settings)};`}
					>
						{#each cargoDestsByPlayer[p.num] as m (key(m))}
							<option value={key(m)}>
								{getMapObjectName(m)}
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
						disabled={!commandable($player.num, selectedMapObject)}
						class="btn btn-outline btn-sm normal-case btn-secondary p-2"
						title="goto">Goto</button
					>
				</div>
				<div class="tooltip" data-tip="merge fleet">
					<button
						onclick={mergeTarget}
						disabled={selectedMapObject.mapObject?.type !== MapObjectType.FLEET ||
							selectedMapObject.mapObject.playerNum !== $player.num}
						class="btn btn-outline btn-sm normal-case btn-secondary p-2"
						title="goto"
						>Merge
					</button>
				</div>
				<div class="tooltip" data-tip="transfer cargo">
					<button
						onclick={transfer}
						disabled={!canLoadFuelOrCargo(fleet, selectedMapObject)}
						class="btn btn-outline btn-sm normal-case btn-secondary p-2"
						title="goto"
						>Transfer
					</button>
				</div>
			</div>
		{/if}
	</CommandTile>
{/if}
