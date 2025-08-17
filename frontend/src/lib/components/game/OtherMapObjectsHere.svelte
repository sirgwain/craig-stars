<script lang="ts">
	import { MapObjectSchema, MapObjectType, type Vector } from '$lib/types/cs-proto';
	import { getGameContext } from '$lib/services/GameContext';
	import { type CommandedFleet } from '$lib/types/Fleet';
	import {
		equal,
		getMapObjectName,
		key,
		type MapObjectLike,
		type MapObjectTargetLike
	} from '$lib/types/MapObject';
	import { create } from '@bufbuild/protobuf';
	import { flatten } from 'lodash-es';
	import type { HTMLSelectAttributes } from 'svelte/elements';

	const { player, universe } = getGameContext();

	type Dictionary<T> = {
		[index: number]: T;
	};

	type Props = {
		fleet: CommandedFleet;
		otherMapObjectsHere: Dictionary<MapObjectLike[]>;
		target: MapObjectTargetLike;
		position: Vector | undefined;
		onSelected: (selected: Partial<MapObjectLike>) => void;
	} & HTMLSelectAttributes;

	let { fleet, otherMapObjectsHere, target, position, onSelected, ...rest }: Props = $props();

	// true if this mapObject is also our current target
	function isTarget(mo: MapObjectLike) {
		if (
			target.targetType === MapObjectType.FLEET ||
			target.targetType === MapObjectType.MINEFIELD ||
			target.targetType === MapObjectType.MINERAL_PACKET
		) {
			// fleets, minefields, and mineral packets are keyed off of player num as well as type/num
			return (
				mo.mapObject?.type === target.targetType &&
				mo.mapObject?.num === target.targetNum &&
				(mo.mapObject?.playerNum ?? 0) === (target.targetPlayerNum ?? 0)
			);
		} else {
			return mo.mapObject?.type === target.targetType && mo.mapObject?.num === target.targetNum;
		}
	}

	function onSelectChange(index: number) {
		const selected = allObjects[index];
		onSelected(selected);
	}

	let everythingElse = $derived(
		flatten(
			Object.entries(otherMapObjectsHere)
				.map(([key, value]) => [Number(key), value] as [MapObjectType, MapObjectLike[]])
				.filter(
					([k]) =>
						k !== MapObjectType.PLANET && k !== MapObjectType.FLEET && k !== MapObjectType.MINEFIELD
				)
				.map(([, value]) => value)
		)
	);
	let allObjects = $derived([
		{
			mapObject: create(MapObjectSchema, {
				type: MapObjectType.UNSPECIFIED,
				position: position
			})
		},
		...(otherMapObjectsHere[MapObjectType.PLANET] ?? []),
		...(otherMapObjectsHere[MapObjectType.FLEET] ?? []),
		...(otherMapObjectsHere[MapObjectType.MINEFIELD] ?? []),
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
		<option selected={target.targetType === MapObjectType.UNSPECIFIED} value={0}
			>{`Space (${position?.x ?? 0}, ${position?.y ?? 0})`}</option
		>
	</optgroup>

	{#if otherMapObjectsHere[MapObjectType.PLANET]}
		<optgroup label="Planets">
			{#each otherMapObjectsHere[MapObjectType.PLANET] as mo, index (key(mo))}
				<option selected={isTarget(mo)} value={1 + index}>{mo.mapObject?.name}</option>
			{/each}
		</optgroup>
	{/if}

	{#if otherMapObjectsHere[MapObjectType.FLEET]}
		<optgroup label="Fleets">
			{#each otherMapObjectsHere[MapObjectType.FLEET] as mo, index (key(mo))}
				{#if !equal(fleet, mo)}
					<option
						style={mo.mapObject?.playerNum != $player.num
							? `color: ${$universe.getPlayerColor(mo.mapObject?.playerNum)};`
							: ''}
						selected={isTarget(mo)}
						value={1 + index + (otherMapObjectsHere[MapObjectType.PLANET]?.length ?? 0)}
						>{getMapObjectName(mo)}</option
					>
				{/if}
			{/each}
		</optgroup>
	{/if}

	{#if otherMapObjectsHere[MapObjectType.MINEFIELD]}
		<optgroup label="Minefields">
			{#each otherMapObjectsHere[MapObjectType.MINEFIELD] as mo, index (key(mo))}
				<option
					selected={isTarget(mo)}
					value={1 +
						index +
						(otherMapObjectsHere[MapObjectType.PLANET]?.length ?? 0) +
						(otherMapObjectsHere[MapObjectType.FLEET]?.length ?? 0)}>{mo.mapObject?.name}</option
				>
			{/each}
		</optgroup>
	{/if}

	{#if everythingElse?.length > 0}
		<optgroup label="Other">
			{#each everythingElse as mo, index (key(mo))}
				<option
					selected={isTarget(mo)}
					value={1 +
						index +
						(otherMapObjectsHere[MapObjectType.PLANET]?.length ?? 0) +
						(otherMapObjectsHere[MapObjectType.FLEET]?.length ?? 0) +
						(otherMapObjectsHere[MapObjectType.MINEFIELD]?.length ?? 0)}
					>{mo.mapObject?.name}</option
				>
			{/each}
		</optgroup>
	{/if}
</select>
