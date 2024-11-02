<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { CommandedFleet } from '$lib/types/Fleet';
	import { kebabCase } from 'lodash-es';
	import CommandTile from './CommandTile.svelte';

	const { player, universe, nextMapObject, previousMapObject, renameFleet } = getGameContext();

	export let fleet: CommandedFleet;

	let icon = '';

	async function onRename() {
		let name = prompt('Enter fleet name', fleet.baseName);
		if (!name || name === '') {
			name = $universe.getMyDesign(fleet.tokens[0].designNum)?.name ?? '';
		}
		if (name !== '') {
			await renameFleet(fleet, name);
		}
	}

	$: {
		icon = '';
		if (fleet.tokens.length > 0) {
			const designNum = fleet.tokens[0].designNum;
			const design = $universe.getDesign($player.num, designNum);
			if (design) {
				icon = `hull-${kebabCase(design.hull)}-${design.hullSetNumber ?? 0}`;
			}
		}
	}
</script>

<CommandTile title={fleet.name}>
	<div class="grid grid-cols-2">
		<div class="avatar border border-secondary p-2 bg-black m-auto relative">
			{#if fleet.tokens.reduce((count, t) => count + t.quantity, 0) > 1}
				<div class="absolute -right-2 -top-1 text-xl w-6 h-6">+</div>
			{/if}
			<div class="fleet-avatar {icon} bg-black" />
		</div>
		<div class="flex flex-col gap-y-1">
			<button
				on:click={() => previousMapObject()}
				type="button"
				class="btn btn-outline btn-sm normal-case btn-secondary">Prev</button
			>
			<button
				on:click={() => nextMapObject()}
				type="button"
				class="btn btn-outline btn-sm normal-case btn-secondary">Next</button
			>
			<button
				on:click={() => onRename()}
				type="button"
				class="btn btn-outline btn-sm normal-case btn-secondary">Rename</button
			>
		</div>
	</div>
</CommandTile>
