<script lang="ts">
	import { run } from 'svelte/legacy';

	import { getGameContext } from '$lib/services/GameContext';
	import type { CommandedFleet } from '$lib/types/Fleet';
	import { kebabCase } from 'lodash-es';
	import CommandTile from './CommandTile.svelte';

	const { player, universe, nextMapObject, previousMapObject, renameFleet } = getGameContext();

	type Props = {
		fleet: CommandedFleet;
	};

	let { fleet }: Props = $props();

	let icon = $state('');

	async function onRename() {
		let name = prompt('Enter fleet name', fleet.baseName);
		if (!name || name === '') {
			name = $universe.getMyDesign(fleet.tokens[0].designNum)?.name ?? '';
		}
		if (name !== '') {
			await renameFleet(fleet, name);
		}
	}

	run(() => {
		icon = '';
		if (fleet.tokens.length > 0) {
			const designNum = fleet.tokens[0].designNum;
			const design = $universe.getDesign($player.num, designNum);
			if (design) {
				icon = `hull-${kebabCase(design.hull)}-${design.hullSetNumber ?? 0}`;
			}
		}
	});
</script>

<CommandTile title={fleet.name}>
	<div class="grid grid-cols-2">
		<div class="avatar border border-secondary p-2 bg-black m-auto relative">
			{#if fleet.tokens.reduce((count, t) => count + t.quantity, 0) > 1}
				<div class="absolute -right-2 -top-1 text-xl w-6 h-6">+</div>
			{/if}
			<div class="fleet-avatar {icon} bg-black"></div>
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
