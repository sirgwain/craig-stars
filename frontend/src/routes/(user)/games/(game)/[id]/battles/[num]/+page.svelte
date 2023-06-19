<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import BattleView from '$lib/components/game/battle/BattleView.svelte';
	import { game } from '$lib/services/Context';

	let id = $page.params.id;
	let num = parseInt($page.params.num);

	$: battle = $game?.player.getBattle(num);
</script>

{#if $game && battle}
	<Breadcrumb>
		<svelte:fragment slot="crumbs">
			<li><a class="cs-link" href={`/games/${id}/battles`}>Battles</a></li>
			<li>{$game.player.getBattleLocation(battle)}</li>
		</svelte:fragment>
	</Breadcrumb>

	<div class="grow px-1">
		<BattleView battleRecord={battle} player={$game.player} />
	</div>
{/if}
