<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import DesignCard from '$lib/components/game/DesignCard.svelte';
	import { game } from '$lib/services/Stores';

	let gameId = parseInt($page.params.id);
</script>

<Breadcrumb>
	<svelte:fragment slot="crumbs">
		<li>Ship Designs</li>
	</svelte:fragment>
	<div slot="end">
		<a class="cs-link btn btn-sm" href={`/games/${gameId}/designs/create`}>Create</a>
	</div>
</Breadcrumb>

{#if $game?.universe.designs.length}
	<div class="flex flex-wrap justify-center gap-2">
		{#each $game?.universe.designs.filter((d) => d.playerNum === $game?.player.num) as design (design.num)}
			<DesignCard
				{design}
				{gameId}
				on:delete={() => design.num && $game?.deleteDesign(design.num)}
			/>
		{/each}
	</div>
{/if}
