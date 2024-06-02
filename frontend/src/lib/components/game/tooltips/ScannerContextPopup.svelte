<script lang="ts" context="module">
	import type { Vector } from '$lib/types/Vector';
	import ScannerContextPopup from './ScannerContextPopup.svelte';

	export function onScannerContextPopup(
		e: PointerEvent | MouseEvent,
		position: Vector | undefined
	) {
		if (position) {
			showPopup<ScannerContextPopupProps>(e.x, e.y, ScannerContextPopup, { position });
		}
	}

	export type ScannerContextPopupProps = {
		position: Vector | undefined;
	};
</script>

<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { showPopup } from '$lib/services/Stores';
	import { MapObjectType, None, ownedBy, type MapObject } from '$lib/types/MapObject';
	import { flatten, keys } from 'lodash-es';
	import { createEventDispatcher } from 'svelte';
	import type { PopupEvent } from './Popup.svelte';
	import type { Fleet } from '$lib/types/Fleet';

	const { player, universe, commandMapObject, selectMapObject } = getGameContext();
	const dispatch = createEventDispatcher<PopupEvent>();

	export let position: Vector;

	$: otherMapObjectsHere = $universe.getOtherMapObjectsHereByType(position);
	$: everythingElse = flatten(
		keys(otherMapObjectsHere).map((k) =>
			k !== MapObjectType.Planet && k !== MapObjectType.Fleet && k !== MapObjectType.MineField
				? otherMapObjectsHere[k]
				: []
		)
	);

	function getTokenCount(mo: MapObject) {
		if (mo.type == MapObjectType.Fleet) {
			const fleet = mo as Fleet;
			return fleet.tokens ? fleet.tokens.reduce((count, t) => count + t.quantity, 0) : 0;
		}
		return 0;
	}

	function gotoTarget(mo: MapObject) {
		if (ownedBy(mo, $player.num)) {
			if (mo.type === MapObjectType.Planet || mo.type === MapObjectType.Fleet) {
				commandMapObject(mo);
			}
		}
		selectMapObject(mo);
		dispatch('close');
	}
</script>

<ul class="menu overflow-y-auto px-0.5">
	{#if otherMapObjectsHere[MapObjectType.Planet]}
		<li class="menu-title w-full">
			Planet
			<ul>
				{#each otherMapObjectsHere[MapObjectType.Planet] as mo}
					<li
						style={mo.playerNum != $player.num && mo.playerNum != None
							? `color: ${$universe.getPlayerColor(mo.playerNum)};`
							: ''}
					>
						<button
							class="py-1 pl-0.5 w-full text-left hover:text-accent"
							on:click={() => gotoTarget(mo)}>{mo.name}</button
						>
					</li>
				{/each}
			</ul>
		</li>
	{/if}
	{#if otherMapObjectsHere[MapObjectType.Fleet]}
		<li class="menu-title w-full">
			Fleets
			<ul>
				{#each otherMapObjectsHere[MapObjectType.Fleet] as mo}
					<li
						style={mo.playerNum != $player.num
							? `color: ${$universe.getPlayerColor(mo.playerNum)};`
							: ''}
					>
						<button
							class="py-1 pl-0.5 w-full text-left hover:text-accent"
							on:click={() => gotoTarget(mo)}>{mo.name} ({getTokenCount(mo)})</button
						>
					</li>
				{/each}
			</ul>
		</li>
	{/if}

	{#if otherMapObjectsHere[MapObjectType.MineField]}
		<li class="menu-title w-full">
			Mine Fields
			<ul>
				{#each otherMapObjectsHere[MapObjectType.MineField] as mo}
					<li
						style={mo.playerNum != $player.num
							? `color: ${$universe.getPlayerColor(mo.playerNum)};`
							: ''}
					>
						<button
							class="py-1 pl-0.5 w-full text-left hover:text-accent"
							on:click={() => gotoTarget(mo)}>{mo.name}</button
						>
					</li>
				{/each}
			</ul>
		</li>
	{/if}
	{#if everythingElse?.length > 0}
		<li class="menu-title w-full">
			Other
			<ul>
				{#each everythingElse as mo}
					<li
						style={mo.playerNum != $player.num
							? `color: ${$universe.getPlayerColor(mo.playerNum)};`
							: ''}
					>
						<button
							class="py-1 pl-0.5 w-full text-left hover:text-accent"
							on:click={() => gotoTarget(mo)}>{mo.name}</button
						>
					</li>
				{/each}
			</ul>
		</li>
	{/if}
</ul>
