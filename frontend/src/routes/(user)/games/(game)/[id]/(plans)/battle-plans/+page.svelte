<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { game } from '$lib/services/Stores';
	import { addError, type CSError } from '$lib/services/Errors';
	import type { BattlePlan } from '$lib/types/Player';
	import BattlePlanCard from './BattlePlanCard.svelte';

	let gameId = parseInt($page.params.id);

	async function deletePlan(plan: BattlePlan) {
		if ($game) {
			try {
				await $game.deleteBattlePlan(plan.num);
				// trigger reactivity
				$game.player.battlePlans = $game.player.battlePlans;
			} catch (e) {
				addError(e as CSError);
			}
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
	<div class="flex flex-wrap justify-center gap-2">
		{#each $game.player.battlePlans as plan (plan.num)}
			<BattlePlanCard
				{plan}
				{gameId}
				showDelete={plan.num !== 0}
				on:delete={() => deletePlan(plan)}
			/>
		{/each}
	</div>
{/if}
