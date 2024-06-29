<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import BattleView from '$lib/components/game/battle/BattleView.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	const { game, player, universe, gotoBattle } = getGameContext();
	let num = parseInt($page.params.num);

	$: battle = $universe.getBattle(num);

	function gotoTarget() {
		if (!battle) {
			return;
		}
		gotoBattle(battle.num);
		goto(`/games/${$game.id}`);
	}
</script>

{#if $game && battle}
	<Breadcrumb>
		<svelte:fragment slot="crumbs">
			<li><a class="cs-link" href={`/games/${$game.id}/battles`}>Battles</a></li>
			<li class="flex flex-row gap-1">
				{$universe.getBattleLocation(battle)}
				<button
					on:click={gotoTarget}
					class="btn btn-outline btn-sm normal-case btn-secondary p-2"
					title="goto">Goto</button
				>
			</li>
		</svelte:fragment>
	</Breadcrumb>

	<div class="grow px-1">
		<BattleView playerFinder={$universe} designFinder={$universe} battleRecord={battle} />
	</div>
{/if}
