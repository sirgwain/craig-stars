<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { type CommandedFleet } from '$lib/types/Fleet';
	import { equal, getMapObjectName } from '$lib/types/MapObject';
	import type { MapObjectTarget, Vector } from '$lib/types/cs';
	import {
		MapObjectTypeFleet,
		MapObjectTypeMineField,
		MapObjectTypeMineralPacket,
		MapObjectTypeNone,
		MapObjectTypePlanet,
		type MapObject
	} from '$lib/types/cs';
	import { flatten, keys } from 'lodash-es';
	import type { HTMLSelectAttributes } from 'svelte/elements';

	const { player, universe } = getGameContext();

	type Dictionary<T> = {
		[index: string]: T;
	};

	type Props = {
		fleet: CommandedFleet;
		otherMapObjectsHere: Dictionary<MapObject[]>;
		target: MapObjectTarget;
		position: Vector;
		onSelected: (selected: Partial<MapObject>) => void;
	} & HTMLSelectAttributes;

	let { fleet, otherMapObjectsHere, target, position, onSelected, ...rest }: Props = $props();

	// true if this mapObject is also our current target
	function isTarget(mo: MapObject) {
		if (
			target.targetType === MapObjectTypeFleet ||
			target.targetType === MapObjectTypeMineField ||
			target.targetType === MapObjectTypeMineralPacket
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
				k !== MapObjectTypePlanet && k !== MapObjectTypeFleet && k !== MapObjectTypeMineField
					? otherMapObjectsHere[k]
					: []
			)
		)
	);
	let allObjects = $derived([
		{ type: MapObjectTypeNone, position: position },
		...(otherMapObjectsHere[MapObjectTypePlanet] ?? []),
		...(otherMapObjectsHere[MapObjectTypeFleet] ?? []),
		...(otherMapObjectsHere[MapObjectTypeMineField] ?? []),
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
		<option selected={target.targetType === MapObjectTypeNone} value={0}
			>{`Space (${position.x ?? 0}, ${position.y ?? 0})`}</option
		>
	</optgroup>

	{#if otherMapObjectsHere[MapObjectTypePlanet]}
		<optgroup label="Planets">
			{#each otherMapObjectsHere[MapObjectTypePlanet] as mo, index}
				<option selected={isTarget(mo)} value={1 + index}>{mo.name}</option>
			{/each}
		</optgroup>
	{/if}

	{#if otherMapObjectsHere[MapObjectTypeFleet]}
		<optgroup label="Fleets">
			{#each otherMapObjectsHere[MapObjectTypeFleet] as mo, index}
				{#if !equal(fleet, mo)}
					<option
						style={mo.playerNum != $player.num
							? `color: ${$universe.getPlayerColor(mo.playerNum)};`
							: ''}
						selected={isTarget(mo)}
						value={1 + index + (otherMapObjectsHere[MapObjectTypePlanet]?.length ?? 0)}
						>{getMapObjectName(mo)}</option
					>
				{/if}
			{/each}
		</optgroup>
	{/if}

	{#if otherMapObjectsHere[MapObjectTypeMineField]}
		<optgroup label="Mine Fields">
			{#each otherMapObjectsHere[MapObjectTypeMineField] as mo, index}
				<option
					selected={isTarget(mo)}
					value={1 +
						index +
						(otherMapObjectsHere[MapObjectTypePlanet]?.length ?? 0) +
						(otherMapObjectsHere[MapObjectTypeFleet]?.length ?? 0)}>{mo.name}</option
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
						(otherMapObjectsHere[MapObjectTypePlanet]?.length ?? 0) +
						(otherMapObjectsHere[MapObjectTypeFleet]?.length ?? 0) +
						(otherMapObjectsHere[MapObjectTypeMineField]?.length ?? 0)}>{mo.name}</option
				>
			{/each}
		</optgroup>
	{/if}
</select>
