<script lang="ts" module>
	import type { Vector } from '$lib/types/cs';
	import ScannerContextPopup from './ScannerContextPopup.svelte';

	export type ScannerContextPopupProps = {
		position: Vector;
	} & PopupProps;

	export function onScannerContextPopup(e: PointerEvent | MouseEvent, position?: Vector) {
		if (position) {
			showPopup<ScannerContextPopupProps>(e.x, e.y, ScannerContextPopup, { position });
		}
	}
</script>

<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import {
		MapObjectTypeFleet,
		MapObjectTypeMinefield,
		MapObjectTypePlanet,
		None,
		type MapObject
	} from '$lib/types/cs';
	import { getMapObjectName, ownedBy } from '$lib/types/MapObject';
	import { flatten, keys } from 'lodash-es';
	import { showPopup, type PopupProps } from './Popup.svelte';

	const { player, universe, commandMapObject, selectMapObject } = getGameContext();

	let { position, onClose }: ScannerContextPopupProps = $props();

	let otherMapObjectsHere = $derived($universe.getOtherMapObjectsHereByType(position));
	let everythingElse = $derived(
		flatten(
			keys(otherMapObjectsHere).map((k) =>
				k !== MapObjectTypePlanet && k !== MapObjectTypeFleet && k !== MapObjectTypeMinefield
					? otherMapObjectsHere[k]
					: []
			)
		)
	);

	function gotoTarget(mo: MapObject) {
		if (ownedBy(mo, $player.num)) {
			if (mo.type === MapObjectTypePlanet || mo.type === MapObjectTypeFleet) {
				commandMapObject(mo);
			}
		}
		selectMapObject(mo);
		onClose?.();
	}
</script>

<ul class="menu overflow-y-auto px-0.5">
	{#if otherMapObjectsHere[MapObjectTypePlanet]}
		<li class="menu-title w-full">
			Planet
			<ul>
				{#each otherMapObjectsHere[MapObjectTypePlanet] as mo (mo)}
					<li
						style={mo.playerNum != $player.num && mo.playerNum != None
							? `color: ${$universe.getPlayerColor(mo.playerNum)};`
							: ''}
					>
						<button
							class="py-1 pl-0.5 w-full text-left hover:text-accent"
							onclick={() => gotoTarget(mo)}>{mo.name}</button
						>
					</li>
				{/each}
			</ul>
		</li>
	{/if}
	{#if otherMapObjectsHere[MapObjectTypeFleet]}
		<li class="menu-title w-full">
			Fleets
			<ul>
				{#each otherMapObjectsHere[MapObjectTypeFleet] as mo (mo)}
					<li
						style={mo.playerNum != $player.num
							? `color: ${$universe.getPlayerColor(mo.playerNum)};`
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

	{#if otherMapObjectsHere[MapObjectTypeMinefield]}
		<li class="menu-title w-full">
			Minefields
			<ul>
				{#each otherMapObjectsHere[MapObjectTypeMinefield] as mo (mo)}
					<li
						style={mo.playerNum != $player.num
							? `color: ${$universe.getPlayerColor(mo.playerNum)};`
							: ''}
					>
						<button
							class="py-1 pl-0.5 w-full text-left hover:text-accent"
							onclick={() => gotoTarget(mo)}
						>
							{mo.name}
						</button>
					</li>
				{/each}
			</ul>
		</li>
	{/if}
	{#if everythingElse?.length > 0}
		<li class="menu-title w-full">
			Other
			<ul>
				{#each everythingElse as mo (mo)}
					<li
						style={mo.playerNum != $player.num
							? `color: ${$universe.getPlayerColor(mo.playerNum)};`
							: ''}
					>
						<button
							class="py-1 pl-0.5 w-full text-left hover:text-accent"
							onclick={() => gotoTarget(mo)}>{mo.name}</button
						>
					</li>
				{/each}
			</ul>
		</li>
	{/if}
</ul>
