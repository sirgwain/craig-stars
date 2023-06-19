<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import DesignCard from '$lib/components/game/DesignCard.svelte';
	import { game } from '$lib/services/Context';

	let gameId = parseInt($page.params.id);
</script>

<Breadcrumb>
	<svelte:fragment slot="crumbs">
		<li>Designs</li>
	</svelte:fragment>
	<div slot="end">
		<a class="cs-link btn btn-sm" href={`/games/${gameId}/designs/create`}>Create</a>
	</div>
</Breadcrumb>

{#if $game?.player.designs.length}
	<div class="flex flex-wrap justify-center">
		{#each $game?.player.designs as design (design.num)}
			<div class="mb-2">
				<DesignCard
					{design}
					{gameId}
					on:delete={() => design.num && $game?.deleteDesign(design.num)}
				/>
			</div>
		{/each}
	</div>
{/if}
