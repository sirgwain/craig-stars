<script lang="ts">
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import TableSearchInput from '$lib/components/TableSearchInput.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import DesignCard from '$lib/components/game/DesignCard.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import type { ShipDesign } from '$lib/types/ShipDesign';

	const { game, player, universe, designs } = getGameContext();

	// filterable designs
	let filteredDesigns: ShipDesign[] = [];
	let search = '';

	$: filteredDesigns =
		$universe
			.getMyDesigns()
			.sort((a, b) => a.name.localeCompare(b.name))
			.filter(
				(i) =>
					i.name.toLowerCase().indexOf(search.toLowerCase()) != -1 ||
					i.hull.toLowerCase().indexOf(search.toLocaleLowerCase()) != -1
			) ?? [];
</script>

<Breadcrumb>
	<svelte:fragment slot="crumbs">
		<li>Ship Designs</li>
	</svelte:fragment>
	<div slot="end">
		<div class="flex flex-row justify-between gap-2 m-2">
			<TableSearchInput bind:value={search} />
			<div>
				<a class="cs-link btn btn-sm" href={`/games/${$game.id}/designer/create`}>Create</a>
			</div>
		</div>
	</div>
</Breadcrumb>

<div class="flex flex-wrap justify-evenly gap-2">
	{#each filteredDesigns.filter((d) => d.playerNum === $player.num && !d.spec.starbase) as design (design.num)}
		<DesignCard
			{design}
			href={`/games/${$game.id}/designer/${design.num}`}
			copyhref={`/games/${$game.id}/designer/create/${design.hull}?copy=${design.num}`}
			on:delete={() => design.num && $game.deleteDesign(design.num)}
		/>
	{/each}
</div>

<ItemTitle>Starbases</ItemTitle>
<div class="flex flex-wrap justify-evenly gap-2">
	{#each filteredDesigns.filter((d) => d.playerNum === $player.num && d.spec.starbase) as design (design.num)}
		<DesignCard
			{design}
			href={`/games/${$game.id}/designer/${design.num}`}
			copyhref={`/games/${$game.id}/designer/create/${design.hull}?copy=${design.num}`}
			on:delete={() => design.num && $game.deleteDesign(design.num)}
		/>
	{/each}
</div>
