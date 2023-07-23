<script lang="ts">
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import DesignCard from '$lib/components/game/DesignCard.svelte';
	import { getGameContext } from '$lib/services/Contexts';

	const { game, player, universe, designs } = getGameContext();
</script>

<Breadcrumb>
	<svelte:fragment slot="crumbs">
		<li>Ship Designs</li>
	</svelte:fragment>
	<div slot="end">
		<a class="cs-link btn btn-sm" href={`/games/${$game.id}/designs/create`}>Create</a>
	</div>
</Breadcrumb>

<div class="flex flex-wrap justify-center gap-2">
	{#each $designs.filter((d) => d.playerNum === $player.num && !d.spec.starbase) as design (design.num)}
		<DesignCard
			{design}
			href={`/games/${$game.id}/designs/${design.num}`}
			copyhref={`/games/${$game.id}/designs/create/${design.hull}?copy=${design.num}`}
			on:delete={() => design.num && $game.deleteDesign(design.num)}
		/>
	{/each}

	{#each $designs.filter((d) => d.playerNum === $player.num && d.spec.starbase) as design (design.num)}
		<DesignCard
			{design}
			href={`/games/${$game.id}/designs/${design.num}`}
			copyhref={`/games/${$game.id}/designs/create/${design.hull}?copy=${design.num}`}
			on:delete={() => design.num && $game.deleteDesign(design.num)}
		/>
	{/each}
</div>
