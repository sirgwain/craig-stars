<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import BattleView from '$lib/components/game/battle/BattleView.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	const { game, player, universe } = getGameContext();
	let num = parseInt($page.params.num);

	$: battle = $universe.getBattle(num);
</script>

{#if $game && battle}
	<Breadcrumb>
		<svelte:fragment slot="crumbs">
			<li><a class="cs-link" href={`/games/${$game.id}/battles`}>Battles</a></li>
			<li>{$universe.getBattleLocation(battle)}</li>
		</svelte:fragment>
	</Breadcrumb>

	<div class="grow px-1">
		<BattleView playerFinder={$universe} designFinder={$universe} battleRecord={battle} />
	</div>
{/if}
