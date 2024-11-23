<script lang="ts" module>
	import type { Vector } from '$lib/types/Vector';
	import ScannerContextPopup from './ScannerContextPopup.svelte';

	export function onScannerContextPopup(e: PointerEvent | MouseEvent, position?: Vector) {
		if (position) {
			showPopup<ScannerContextPopupProps>(e.x, e.y, ScannerContextPopup, { position });
		}
	}

	export type ScannerContextPopupProps = {
		position: Vector;
	};
</script>

<script lang="ts">
	import type { OnClose } from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import { showPopup } from '$lib/services/Stores';
	import { None } from '$lib/types/Constants';
	import { getMapObjectName, MapObjectType, ownedBy, type MapObject } from '$lib/types/MapObject';
	import { flatten, keys } from 'lodash-es';

	const { player, universe, commandMapObject, selectMapObject } = getGameContext();

	type Props = {
		onClose?: OnClose;
	} & ScannerContextPopupProps;

	let { position, onClose }: Props = $props();

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
