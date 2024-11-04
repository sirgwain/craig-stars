<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import Designs from '../Designs.svelte';

	const { game, universe } = getGameContext();
	let playerNum = parseInt($page.params.playerNum);

	let designs = $derived($universe.getDesigns(playerNum));
</script>

<Breadcrumb>
	{#snippet crumbs()}
		<li><a href={`/games/${$game.id}/designs`}>Designs</a></li>
		<li>{$universe.getPlayerPluralName(playerNum)}</li>
	{/snippet}
</Breadcrumb>

<Designs {designs} />
