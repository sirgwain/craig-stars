<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import Designs from '../Designs.svelte';

	const { game, universe } = getGameContext();
	let playerNum = parseInt($page.params.playerNum);

	$: designs = $universe.getDesigns(playerNum);
</script>

<Breadcrumb>
	<svelte:fragment slot="crumbs">
		<li><a href={`/games/${$game.id}/designs`}>Designs</a></li>
		<li>{$universe.getPlayerName(playerNum)}</li>
	</svelte:fragment>
</Breadcrumb>

<Designs {designs} />
