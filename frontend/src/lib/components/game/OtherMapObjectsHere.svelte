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
props.props.props.props.
	interface Props {
		fleet: CommandedFleet;
		otherMapObjectsHere: Dictionary<MapObject[]>;
		target: Target;
		position: Vector;
		[key: string]: any
	}

	let {
		...props
	}: Props = $props();

	// true if this mapObject is also our current target
	function isTarget(mo: MapObject) {
		return (
			mo.type === props.target.targetType &&
			mo.num === props.target.targetNum &&
			(mo.playerNum ?? 0) === (props.target.targetPlayerNum ?? 0)
		);
	}

	function onSelectChange(index: number) {
		const selected = allObjects[index];
		dispatch('selected', selected);
	}

	let everythingElse = $derived(flatten(
		keys(props.otherMapObjectsHere).map((k) =>
			k !== MapObjectType.Planet && k !== MapObjectType.Fleet && k !== MapObjectType.MineField
				? props.otherMapObjectsHere[k]
				: []
		)
	));
	let allObjects = $derived([
		{ type: MapObjectType.None, props.position: props.position },
		...(props.otherMapObjectsHere[MapObjectType.Planet] ?? []),
		...(props.otherMapObjectsHere[MapObjectType.Fleet] ?? []),
		...(props.otherMapObjectsHere[MapObjectType.MineField] ?? []),
		...everythingElse
	]);
</script>

<select
	style={props.target.targetPlayerNum && props.target.targetPlayerNum != $player.num
		? `color: ${$universe.getPlayerColor(props.target.targetPlayerNum)};`
		: ''}
	onchange={(e) => onSelectChange(parseInt(e.currentTarget.value))}
	class={`select select-outline select-secondary select-sm text-sm ${props.class}`}
>
	<!-- allow for the non target -->
	<optgroup label="Space">
		<option selected={props.target.targetType === MapObjectType.None} value={0}
			>{`Space (${props.position.x ?? 0}, ${props.position.y ?? 0})`}</option
		>
	</optgroup>

	{#if props.otherMapObjectsHere[MapObjectType.Planet]}
		<optgroup label="Planets">
			{#each props.otherMapObjectsHere[MapObjectType.Planet] as mo, index}
				<option selected={isTarget(mo)} value={1 + index}>{mo.name}</option>
			{/each}
		</optgroup>
	{/if}

	{#if props.otherMapObjectsHere[MapObjectType.Fleet]}
		<optgroup label="Fleets">
			{#each props.otherMapObjectsHere[MapObjectType.Fleet] as mo, index}
				{#if !equal(props.fleet, mo)}
					<option
						style={mo.playerNum != $player.num
							? `color: ${$universe.getPlayerColor(mo.playerNum)};`
							: ''}
						selected={isTarget(mo)}
						value={1 + index + (props.otherMapObjectsHere[MapObjectType.Planet]?.length ?? 0)}
						>{getMapObjectName(mo)}</option
					>
				{/if}
			{/each}
		</optgroup>
	{/if}

	{#if props.otherMapObjectsHere[MapObjectType.MineField]}
		<optgroup label="Mine Fields">
			{#each props.otherMapObjectsHere[MapObjectType.MineField] as mo, index}
				<option
					selected={isTarget(mo)}
					value={1 +
						index +
						(props.otherMapObjectsHere[MapObjectType.Planet]?.length ??
							0 + props.otherMapObjectsHere[MapObjectType.Fleet]?.length ??
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
						(props.otherMapObjectsHere[MapObjectType.Planet]?.length ??
							0 + props.otherMapObjectsHere[MapObjectType.Fleet]?.length ??
							0 + props.otherMapObjectsHere[MapObjectType.MineField]?.length ??
							0)}>{mo.name}</option
				>
			{/each}
		</optgroup>
	{/if}
</select>
