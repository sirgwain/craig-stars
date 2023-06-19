<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { game } from '$lib/services/Context';
	import type { BattlePlan } from '$lib/types/Player';
	import BattlePlanCard from './BattlePlanCard.svelte';

	let gameId = parseInt($page.params.id);

	function deletePlan(plan: BattlePlan) {
		if ($game) {
			$game.player.battlePlans = $game.player.battlePlans.filter((p) => p.num !== plan.num);
			$game.updatePlayerPlans();
		}
	}
</script>

<Breadcrumb>
	<svelte:fragment slot="crumbs">
		<li>Battle Plans</li>
	</svelte:fragment>
	<div slot="end">
		<a class="cs-link btn btn-sm" href={`/games/${gameId}/battle-plans/create`}>Create</a>
	</div>
</Breadcrumb>

{#if $game?.player.battlePlans.length}
	<div class="flex flex-wrap justify-center">
		{#each $game.player.battlePlans as plan (plan.num)}
			<div class="mb-2">
				<BattlePlanCard
					{plan}
					{gameId}
					showDelete={plan.num !== 0}
					on:delete={() => deletePlan(plan)}
				/>
			</div>
		{/each}
	</div>
{/if}
