<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import BattleView from '$lib/components/game/battle/BattleView.svelte';
	import { game } from '$lib/services/Stores';

	let id = $page.params.id;
	let num = parseInt($page.params.num);

	$: battle = $game?.universe.getBattle(num);
</script>

{#if $game && battle}
	<Breadcrumb>
		<svelte:fragment slot="crumbs">
			<li><a class="cs-link" href={`/games/${id}/battles`}>Battles</a></li>
			<li>{$game.universe.getBattleLocation(battle, $game.universe)}</li>
		</svelte:fragment>
	</Breadcrumb>

	<div class="grow px-1">
		<BattleView universe={$game.universe} battleRecord={battle} player={$game.player} />
	</div>
{/if}
