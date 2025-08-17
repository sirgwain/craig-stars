<script lang="ts">
	import type { NextPrevMapObjectProps, RenameFleetProps } from '$lib/services/Events';

	import { getGameContext } from '$lib/services/GameContext';
	import type { AnyShipDesign } from '$lib/services/Universe';
	import { getHullIcon } from '$lib/techicon';
	import type { CommandedFleet } from '$lib/types/Fleet';
	import CommandTile from './CommandTile.svelte';

	const { universe } = getGameContext();

	type Props = {
		fleet: CommandedFleet;
		hideTitle?: boolean;
	} & NextPrevMapObjectProps &
		RenameFleetProps;

	let { fleet, hideTitle, onRenameFleet, onNextMapObject, onPreviousMapObject }: Props = $props();

	const design: AnyShipDesign | undefined = $derived.by(() => {
		if (fleet.tokens && fleet.tokens.length > 0) {
			const designNum = fleet.tokens[0].designNum;
			return $universe.getDesign(fleet.mapObject.playerNum, designNum);
		}
	});

	function rename() {
		let name = prompt('Enter fleet name', fleet.baseName);
		if (!name || name === '') {
			name = $universe.getMyDesign(fleet.tokens[0].designNum)?.name ?? '';
		}
		if (name !== '') {
			onRenameFleet?.({ fleet, name });
		}
	}
</script>

<CommandTile title={hideTitle ? '' : fleet.mapObject.name}>
	<div class="grid grid-cols-2">
		<div class="avatar border border-secondary p-2 bg-black m-auto relative">
			{#if fleet.tokens.reduce((count, t) => count + t.quantity, 0) > 1}
				<div class="absolute -right-2 -top-1 text-xl w-6 h-6">+</div>
			{/if}
			<div class="fleet-avatar {getHullIcon(design)} bg-black"></div>
		</div>
		<div class="flex flex-col gap-y-1">
			<button
				onclick={onPreviousMapObject}
				type="button"
				class="btn btn-outline btn-sm normal-case btn-secondary">Prev</button
			>
			<button
				onclick={onNextMapObject}
				type="button"
				class="btn btn-outline btn-sm normal-case btn-secondary">Next</button
			>
			<button
				onclick={rename}
				type="button"
				class="btn btn-outline btn-sm normal-case btn-secondary">Rename</button
			>
		</div>
	</div>
</CommandTile>
