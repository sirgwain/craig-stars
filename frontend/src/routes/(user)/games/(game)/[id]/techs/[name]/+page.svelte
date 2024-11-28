<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import TechHullSummary from '$lib/components/game/design/Hull.svelte';
	import TechSummary from '$lib/components/tech/TechSummary.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	import { TechCategory, type TechHull } from '$lib/types/Tech';

	const { game, player } = getGameContext();

	let nameSlug = $page.params.name;
	let tech = $derived($game.techs.getTech(nameSlug));

	let hull = $derived(tech as TechHull);
</script>

<Breadcrumb>
	{#snippet crumbs()}
		<li><a href={`/games/${$game.id}/techs`}>Techs</a></li>
		<li>{tech?.name ?? '<unknown>'}</li>
	{/snippet}
</Breadcrumb>

{#if tech}
	<TechSummary {tech} player={$player} />
	{#if (hull && tech.category == TechCategory.ShipHull) || tech.category == TechCategory.StarbaseHull}
		<h1 class="my-3 text-lg text-center font-semibold">Hull</h1>
		<div
			class="card bg-base-200 shadow w-full max-h-fit min-h-fit rounded-sm border-2 border-base-300"
		>
			<div class="w-full flex flex-row justify-center">
				<TechHullSummary {hull} />
			</div>
		</div>
	{/if}
{/if}
