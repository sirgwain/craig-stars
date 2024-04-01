<script lang="ts">
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { addError, type CSError } from '$lib/services/Errors';
	import type { ProductionPlan } from '$lib/types/Player';
	import ProductionPlanCard from './ProductionPlanCard.svelte';

	const { game, player, universe, deleteProductionPlan } = getGameContext();

	async function deletePlan(plan: ProductionPlan) {
		if ($game) {
			try {
				await deleteProductionPlan(plan.num);
				// trigger reactivity
				$player.productionPlans = $player.productionPlans;
			} catch (e) {
				addError(e as CSError);
			}
		}
	}
</script>

<Breadcrumb>
	<svelte:fragment slot="crumbs">
		<li>Production Plans</li>
	</svelte:fragment>

	<div slot="end">
		<a class="cs-link btn btn-sm" href={`/games/${$game.id}/production-plans/create`}>Create</a>
	</div>
</Breadcrumb>

{#if $player.productionPlans.length}
	<div class="flex flex-wrap justify-center gap-2">
		{#each $player.productionPlans as plan (plan.num)}
			<ProductionPlanCard
				designFinder={$universe}
				{plan}
				href={`/games/${$game.id}/production-plans/${plan.num}`}
				showDelete={plan.num !== 0}
				on:delete={() => deletePlan(plan)}
			/>
		{/each}
	</div>
{/if}
