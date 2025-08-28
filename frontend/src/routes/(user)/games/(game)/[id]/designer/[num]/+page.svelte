<script lang="ts">
	import { page } from '$app/state';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import Design from '$lib/components/game/design/Design.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	const { game, universe } = getGameContext();
	let num = parseInt(page.params.num || '0');

	let design = $derived($universe.getMyDesign(num));
</script>

{#if design}
	<Breadcrumb>
		{#snippet crumbs()}
			<li><a class="cs-link" href={`/games/${$game.id}/designer`}>Ship Designs</a></li>
			<li>{design?.name}</li>
			{#if !design.spec?.numInstances}
				<li><a class="cs-link" href={`/games/${$game.id}/designer/${design.num}/edit`}>Edit</a></li>
			{/if}
		{/snippet}
	</Breadcrumb>

	<div class="grow h-full px-1 md:p-0">
		<Design {design} />
	</div>
{/if}
