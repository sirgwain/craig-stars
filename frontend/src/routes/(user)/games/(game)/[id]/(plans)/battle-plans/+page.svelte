<script lang="ts">
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import { addError, type CSError } from '$lib/services/Errors';
	import type { BattlePlan } from '$lib/types/Player';
	import BattlePlanCard from './BattlePlanCard.svelte';

	const { game, player, universe } = getGameContext();

	async function deletePlan(plan: BattlePlan) {
		if ($game) {
			try {
				await $game.deleteBattlePlan(plan.num);
				// trigger reactivity
				$player.battlePlans = $player.battlePlans;
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
		<a class="cs-link btn btn-sm" href={`/games/${$game.id}/battle-plans/create`}>Create</a>
	</div>
</Breadcrumb>

{#if $player.battlePlans.length}
	<div class="flex flex-wrap justify-center gap-2">
		{#each $player.battlePlans as plan (plan.num)}
			<BattlePlanCard
				{plan}
				href={`/games/${$game.id}/battle-plans/${plan.num}`}
				showDelete={plan.num !== 0}
				on:delete={() => deletePlan(plan)}
			/>
		{/each}
	</div>
{/if}
