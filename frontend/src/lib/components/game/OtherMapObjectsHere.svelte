<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import type { CommandedFleet, Target } from '$lib/types/Fleet';
	import { MapObjectType, equal, type MapObject } from '$lib/types/MapObject';
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
	{#if otherMapObjectsHere[MapObjectType.Planet]}
		<optgroup label="Planets">
			{#each otherMapObjectsHere[MapObjectType.Planet] as mo, index}
				<option selected={isTarget(mo)} value={index}>{mo.name}</option>
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
						value={index + (otherMapObjectsHere[MapObjectType.Planet]?.length ?? 0)}
						>{mo.name}</option
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
					value={index +
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
					value={index +
						(otherMapObjectsHere[MapObjectType.Planet]?.length ??
							0 + otherMapObjectsHere[MapObjectType.Fleet]?.length ??
							0 + otherMapObjectsHere[MapObjectType.MineField]?.length ??
							0)}>{mo.name}</option
				>
			{/each}
		</optgroup>
	{/if}
</select>
