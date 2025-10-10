<script lang="ts" module>
	import ScannerContextPopup from './ScannerContextPopup.svelte';

	export type ScannerContextPopupProps = {
		position: Position;
	} & PopupPropsBase;

	export function onScannerContextPopup(e: PointerEvent | MouseEvent, position: Position) {
		showPopup<ScannerContextPopupProps>(e.x, e.y, ScannerContextPopup, { position });
	}
</script>

<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { None } from '$lib/types/Consts';
	import { MapObjectType } from '$lib/types/cs-proto';
	import {
		getMapObjectName,
		ownedBy,
		type MapObjectLike,
		type Position
	} from '$lib/types/MapObject';
	import { getDisplayColor } from '$lib/utils/colorUtils';
	import { flatten } from 'lodash-es';
	import { showPopup, type PopupPropsBase } from './Popup.svelte';

	const { player, universe, settings, commandMapObject, selectMapObject } = getGameContext();

	let { position, onClose }: ScannerContextPopupProps = $props();

	let otherMapObjectsHere = $derived($universe.getOtherMapObjectsHereByType(position));
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
	function gotoTarget(mo: MapObjectLike) {
		if (ownedBy(mo, $player.num)) {
			if (
				mo.mapObject?.type === MapObjectType.PLANET ||
				mo.mapObject?.type === MapObjectType.FLEET
			) {
				commandMapObject(mo);
			}
		}
		selectMapObject(mo);
		onClose?.();
	}
</script>

<ul data-type="popup" data-id="scanner-context-popup" class="menu overflow-y-auto px-0.5">
	{#if otherMapObjectsHere[MapObjectType.PLANET]}
		<li class="menu-title w-full">
			Planet
			<ul>
				{#each otherMapObjectsHere[MapObjectType.PLANET] as mo (mo)}
					<li
						style={mo.mapObject?.playerNum != $player.num && mo.mapObject?.playerNum != None
							? `color: ${getDisplayColor(mo.mapObject?.playerNum, $player, $universe, $settings)};`
							: ''}
					>
						<button
							class="py-1 pl-0.5 w-full text-left hover:text-accent"
							onclick={() => gotoTarget(mo)}>{mo.mapObject?.name}</button
						>
					</li>
				{/each}
			</ul>
		</li>
	{/if}
	{#if otherMapObjectsHere[MapObjectType.FLEET]}
		<li class="menu-title w-full">
			Fleets
			<ul>
				{#each otherMapObjectsHere[MapObjectType.FLEET] as mo (mo)}
					<li
						style={mo.mapObject?.playerNum != $player.num
							? `color: ${getDisplayColor(mo.mapObject?.playerNum, $player, $universe, $settings)};`
							: ''}
					>
						<button
							class="py-1 pl-0.5 w-full text-left hover:text-accent"
							onclick={() => gotoTarget(mo)}>{getMapObjectName(mo)}</button
						>
					</li>
				{/each}
			</ul>
		</li>
	{/if}

	{#if otherMapObjectsHere[MapObjectType.MINEFIELD]}
		<li class="menu-title w-full">
			Minefields
			<ul>
				{#each otherMapObjectsHere[MapObjectType.MINEFIELD] as mo (mo)}
					<li
						style={mo.mapObject?.playerNum != $player.num
							? `color: ${getDisplayColor(mo.mapObject?.playerNum, $player, $universe, $settings)};`
							: ''}
					>
						<button
							class="py-1 pl-0.5 w-full text-left hover:text-accent"
							onclick={() => gotoTarget(mo)}
						>
							{mo.mapObject?.name}
						</button>
					</li>
				{/each}
			</ul>
		</li>
	{/if}
	{#if everythingElse.length > 0}
		<li class="menu-title w-full">
			Other
			<ul>
				{#each everythingElse as mo (mo)}
					<li
						style={mo.mapObject?.playerNum != $player.num
							? `color: ${getDisplayColor(mo.mapObject?.playerNum, $player, $universe, $settings)};`
							: ''}
					>
						<button
							class="py-1 pl-0.5 w-full text-left hover:text-accent"
							onclick={() => gotoTarget(mo)}>{mo.mapObject?.name}</button
						>
					</li>
				{/each}
			</ul>
		</li>
	{/if}
</ul>
