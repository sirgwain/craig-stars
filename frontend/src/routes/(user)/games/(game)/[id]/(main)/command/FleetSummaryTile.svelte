<script lang="ts">
	
	import { getGameContext } from '$lib/services/GameContext';
	import { getHullIcon } from '$lib/techicon';
	import type { CommandedFleet } from '$lib/types/Fleet';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import CommandTile from './CommandTile.svelte';

	const { player, universe, nextMapObject, previousMapObject, renameFleet } = getGameContext();

	type Props = {
		fleet: CommandedFleet;
	};

	let { fleet }: Props = $props();

	const design: ShipDesign | undefined = $derived.by(() => {
		if (fleet.tokens && fleet.tokens.length > 0) {
			const designNum = fleet.tokens[0].designNum;
			return $universe.getDesign(fleet.playerNum, designNum);
		}
	});

	async function onRename() {
		let name = prompt('Enter fleet name', fleet.baseName);
		if (!name || name === '') {
			name = $universe.getMyDesign(fleet.tokens[0].designNum)?.name ?? '';
		}
		if (name !== '') {
			await renameFleet(fleet, name);
		}
	}
</script>

<CommandTile title={fleet.name}>
	<div class="grid grid-cols-2">
		<div class="avatar border border-secondary p-2 bg-black m-auto relative">
			{#if fleet.tokens.reduce((count, t) => count + t.quantity, 0) > 1}
				<div class="absolute -right-2 -top-1 text-xl w-6 h-6">+</div>
			{/if}
			<div class="fleet-avatar {getHullIcon(design)} bg-black"></div>
		</div>
		<div class="flex flex-col gap-y-1">
			<button
				onclick={() => previousMapObject()}
				type="button"
				class="btn btn-outline btn-sm normal-case btn-secondary">Prev</button
			>
			<button
				onclick={() => nextMapObject()}
				type="button"
				class="btn btn-outline btn-sm normal-case btn-secondary">Next</button
			>
			<button
				onclick={() => onRename()}
				type="button"
				class="btn btn-outline btn-sm normal-case btn-secondary">Rename</button
			>
		</div>
	</div>
</CommandTile>
