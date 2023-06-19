<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
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

<Breadcrumb>
	<svelte:fragment slot="crumbs">
		<li>Designs</li>
	</svelte:fragment>
	<div slot="end">
		<a class="cs-link btn btn-sm" href={`/games/${gameId}/designs/create`}>Create</a>
	</div>
</Breadcrumb>

{#if designs?.length && gameId != undefined}
	<div class="flex flex-wrap justify-center">
		{#each designs as design (design.num)}
			<div class="mb-2">
				<DesignCard {design} {gameId} on:deleted={(e) => onDeleted(e.detail.design)} />
			</div>
		{/each}
	</div>
{/if}
