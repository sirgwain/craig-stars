<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import ShipDesigner from '$lib/components/game/design/ShipDesigner.svelte';
	import { DesignService } from '$lib/services/DesignService';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import { onMount } from 'svelte';

	let gameId = $page.params.id;
	let num = parseInt($page.params.num);

	let design: ShipDesign | undefined;

	let error = '';

	onMount(async () => {
		try {
			design = await DesignService.get(gameId, num);
		} catch (e) {
			error = (e as Error).message;
		}
	});

	const onSave = async () => {
		error = '';

		try {
			if (design) {
				design = await DesignService.update(gameId, design);
				goto(`/games/${gameId}/designs/${design.num}`);
			}
		} catch (e) {
			error = (e as Error).message;
		}
	};
</script>

<div class="w-full mx-auto md:max-w-2xl">
	{#if design}
		<ShipDesigner {gameId} bind:design hullName={design.hull} on:save={(e) => onSave()} bind:error />
	{/if}
</div>
