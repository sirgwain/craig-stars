<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import Design from '$lib/components/game/design/Design.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	const { game, universe, player } = getGameContext();
	let num = parseInt($page.params.num);

	$: design = $universe.designs.find((d) => d.playerNum == $player.num && d.num === num);
</script>

{#if design}
	<Breadcrumb>
		<svelte:fragment slot="crumbs">
			<li><a class="cs-link" href={`/games/${$game.id}/designer`}>Ship Designs</a></li>
			<li>{design?.name}</li>
			{#if !design.spec?.numInstances}
				<li><a class="cs-link" href={`/games/${$game.id}/designer/${design.num}/edit`}>Edit</a></li>
			{/if}
		</svelte:fragment>
	</Breadcrumb>

	<div class="grow h-full px-1 md:p-0">
		<Design {design} />
	</div>
{/if}
