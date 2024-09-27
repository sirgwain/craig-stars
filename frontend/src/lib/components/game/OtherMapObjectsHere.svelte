<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { type CommandedFleet, type Target } from '$lib/types/Fleet';
	import { MapObjectType, equal, getMapObjectName, type MapObject } from '$lib/types/MapObject';
	import type { Vector } from '$lib/types/Vector';
	import { flatten, keys } from 'lodash-es';
	import { createEventDispatcher } from 'svelte';

	const { game, player, universe, settings } = getGameContext();
	const dispatch = createEventDispatcher();

	interface Dictionary<T> {
		[index: string]: T;
	}

	export let fleet: CommandedFleet;
	export let otherMapObjectsHere: Dictionary<MapObject[]>;
	export let target: Target;
	export let position: Vector;

	// true if this mapObject is also our current target
	function isTarget(mo: MapObject) {
		return (
			mo.type == target.targetType &&
			mo.num == target.targetNum &&
			mo.playerNum == target.targetPlayerNum
		);
	}

	function onSelectChange(index: number) {
		const selected = allObjects[index];
		dispatch('selected', selected);
	}

	$: everythingElse = flatten(
		keys(otherMapObjectsHere).map((k) =>
			k !== MapObjectType.Planet && k !== MapObjectType.Fleet && k !== MapObjectType.MineField
				? otherMapObjectsHere[k]
				: []
		)
	);
	$: allObjects = [
		{ type: MapObjectType.None, position: position },
		...(otherMapObjectsHere[MapObjectType.Planet] ?? []),
		...(otherMapObjectsHere[MapObjectType.Fleet] ?? []),
		...(otherMapObjectsHere[MapObjectType.MineField] ?? []),
		...everythingElse
	];
</script>

<select
	style={target.targetPlayerNum && target.targetPlayerNum != $player.num
		? `color: ${$universe.getPlayerColor(target.targetPlayerNum)};`
		: ''}
	on:change={(e) => onSelectChange(parseInt(e.currentTarget.value))}
	class={`select select-outline select-secondary select-sm text-sm ${$$props.class}`}
>
	<!-- allow for the non target -->
	<optgroup label="Space">
		<option selected={target.targetType === MapObjectType.None} value={0}
			>{`Space (${position.x ?? 0}, ${position.y ?? 0})`}</option
		>
	</optgroup>

	{#if otherMapObjectsHere[MapObjectType.Planet]}
		<optgroup label="Planets">
			{#each otherMapObjectsHere[MapObjectType.Planet] as mo, index}
				<option selected={isTarget(mo)} value={1 + index}>{mo.name}</option>
			{/each}
		</optgroup>
	{/if}

	{#if otherMapObjectsHere[MapObjectType.Fleet]}
		<optgroup label="Fleets">
			{#each otherMapObjectsHere[MapObjectType.Fleet] as mo, index}
				{#if !equal(fleet, mo)}
					<option
						style={mo.playerNum != $player.num
							? `color: ${$universe.getPlayerColor(mo.playerNum)};`
							: ''}
						selected={isTarget(mo)}
						value={1 + index + (otherMapObjectsHere[MapObjectType.Planet]?.length ?? 0)}
						>{getMapObjectName(mo)}</option
					>
				{/if}
			{/each}
		</optgroup>
	{/if}

	{#if otherMapObjectsHere[MapObjectType.MineField]}
		<optgroup label="Mine Fields">
			{#each otherMapObjectsHere[MapObjectType.MineField] as mo, index}
				<option
					selected={isTarget(mo)}
					value={1 +
						index +
						(otherMapObjectsHere[MapObjectType.Planet]?.length ??
							0 + otherMapObjectsHere[MapObjectType.Fleet]?.length ??
							0)}>{mo.name}</option
				>
			{/each}
		</optgroup>
	{/if}

	{#if everythingElse?.length > 0}
		<optgroup label="Other">
			{#each everythingElse as mo, index}
				<option
					selected={isTarget(mo)}
					value={1 +
						index +
						(otherMapObjectsHere[MapObjectType.Planet]?.length ??
							0 + otherMapObjectsHere[MapObjectType.Fleet]?.length ??
							0 + otherMapObjectsHere[MapObjectType.MineField]?.length ??
							0)}>{mo.name}</option
				>
			{/each}
		</optgroup>
	{/if}
</select>
