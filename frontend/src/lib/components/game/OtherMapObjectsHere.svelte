<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { type CommandedFleet, type Target } from '$lib/types/Fleet';
	import { MapObjectType, equal, getMapObjectName, type MapObject } from '$lib/types/MapObject';
	import type { Vector } from '$lib/types/Vector';
	import { flatten, keys } from 'lodash-es';
	import type { HTMLAttributes, HTMLSelectAttributes } from 'svelte/elements';

	const { player, universe } = getGameContext();

	interface Dictionary<T> {
		[index: string]: T;
	}

	type Props = {
		fleet: CommandedFleet;
		otherMapObjectsHere: Dictionary<MapObject[]>;
		target: Target;
		position: Vector;
		onSelected: (selected: Partial<MapObject>) => void;
	} & HTMLSelectAttributes;

	let { fleet, otherMapObjectsHere, target, position, onSelected, ...rest }: Props = $props();

	// true if this mapObject is also our current target
	function isTarget(mo: MapObject) {
		if (
			target.targetType === MapObjectType.Fleet ||
			target.targetType === MapObjectType.MineField ||
			target.targetType === MapObjectType.MineralPacket
		) {
			// fleets, minefields, and mineral packets are keyed off of player num as well as type/num
			return (
				mo.type === target.targetType &&
				mo.num === target.targetNum &&
				(mo.playerNum ?? 0) === (target.targetPlayerNum ?? 0)
			);
		} else {
			return mo.type === target.targetType && mo.num === target.targetNum;
		}
	}

	function onSelectChange(index: number) {
		const selected = allObjects[index];
		onSelected(selected);
	}

	let everythingElse = $derived(
		flatten(
			keys(otherMapObjectsHere).map((k) =>
				k !== MapObjectType.Planet && k !== MapObjectType.Fleet && k !== MapObjectType.MineField
					? otherMapObjectsHere[k]
					: []
			)
		)
	);
	let allObjects = $derived([
		{ type: MapObjectType.None, position: position },
		...(otherMapObjectsHere[MapObjectType.Planet] ?? []),
		...(otherMapObjectsHere[MapObjectType.Fleet] ?? []),
		...(otherMapObjectsHere[MapObjectType.MineField] ?? []),
		...everythingElse
	]);
</script>

<select
	style={target.targetPlayerNum && target.targetPlayerNum != $player.num
		? `color: ${$universe.getPlayerColor(target.targetPlayerNum)};`
		: ''}
	onchange={(e) => onSelectChange(parseInt(e.currentTarget.value))}
	class={`select select-outline select-secondary select-sm text-sm ${rest.class ?? ''}`}
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
