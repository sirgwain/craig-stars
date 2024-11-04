<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import Design from '$lib/components/game/design/Design.svelte';
	import DesignCard from '$lib/components/game/DesignCard.svelte';
	import NotFound from '$lib/components/NotFound.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import Designs from '../../Designs.svelte';

	const { game, universe } = getGameContext();
	let playerNum = parseInt($page.params.playerNum);
	let num = parseInt($page.params.num);

	let design = $derived($universe.getDesign(playerNum, num));
</script>

<Breadcrumb>
	{#snippet crumbs()}
		<li><a href={`/games/${$game.id}/designs`}>Designs</a></li>
		<li>
			<a href={`/games/${$game.id}/designs/${playerNum}`}
				>{$universe.getPlayerPluralName(playerNum)}</a
			>
		</li>
		<li>{design?.name ?? 'not found'}</li>
	{/snippet}
</Breadcrumb>

{#if design}
	<div class="bg-base-200 rounded-md p-2">
		<Design {design} />
	</div>
{:else}
	<NotFound />
{/if}
