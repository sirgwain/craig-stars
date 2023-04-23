<script lang="ts">
	import { page } from '$app/stores';
	import DesignCard from '$lib/components/game/DesignCard.svelte';
	import Design from '$lib/components/game/design/Design.svelte';
	import { DesignService } from '$lib/services/DesignService';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import { onMount } from 'svelte';

	let gameId = parseInt($page.params.id);

	let designs: ShipDesign[] = [];

	onMount(async () => {
		try {
			designs = await DesignService.load(gameId);
		} catch (error) {
			// TODO: handle error
		}
	});

	function onDeleted(design: ShipDesign): void {
		designs = designs.filter((d) => d.num !== design.num);
	}
</script>

<div class="w-full mx-auto md:max-w-2xl">
	<div class="w-full flex justify-between gap-2 border-primary border-b-2 mb-2">
		<div class="breadcrumbs">
			<ul>
				<li>Designs</li>
			</ul>
		</div>

		<a class="cs-link btn btn-sm" href={`/games/${gameId}/designs/create`}>Create</a>
	</div>

	{#if designs?.length && gameId != undefined}
		<div class="flex flex-wrap justify-center">
			{#each designs as design (design.num)}
				<div class="mb-2">
					<DesignCard {design} {gameId} on:deleted={(e) => onDeleted(e.detail.design)} />
				</div>
			{/each}
		</div>
	{/if}
</div>
