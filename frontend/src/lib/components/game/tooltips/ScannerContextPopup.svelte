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
	import { None } from '$lib/types/cs';
	import { getMapObjectName, MapObjectType, ownedBy } from '$lib/types/MapObject';
	import { type MapObject } from '$lib/types/cs';
	import { flatten, keys } from 'lodash-es';
	import { showPopup, type PopupProps } from './Popup.svelte';

	const { player, universe, commandMapObject, selectMapObject } = getGameContext();

	let { position, onClose }: ScannerContextPopupProps = $props();

	let otherMapObjectsHere = $derived($universe.getOtherMapObjectsHereByType(position));
	let everythingElse = $derived(
		flatten(
			keys(otherMapObjectsHere).map((k) =>
				k !== MapObjectType.Planet && k !== MapObjectType.Fleet && k !== MapObjectType.MineField
					? otherMapObjectsHere[k]
					: []
			)
		)
	);

	function gotoTarget(mo: MapObject) {
		if (ownedBy(mo, $player.num)) {
			if (mo.type === MapObjectType.Planet || mo.type === MapObjectType.Fleet) {
				commandMapObject(mo);
			}
		}
		selectMapObject(mo);
		onClose?.();
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
							onclick={() => gotoTarget(mo)}>{mo.name}</button
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
							onclick={() => gotoTarget(mo)}>{getMapObjectName(mo)}</button
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
				{#each everythingElse as mo}
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
