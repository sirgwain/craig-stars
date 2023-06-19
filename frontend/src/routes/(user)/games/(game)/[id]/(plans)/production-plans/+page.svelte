<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { game } from '$lib/services/Context';
	import { addError, type CSError } from '$lib/services/Errors';
	import type { ProductionPlan } from '$lib/types/Player';
	import ProductionPlanCard from './ProductionPlanCard.svelte';

	let gameId = parseInt($page.params.id);

	async function deletePlan(plan: ProductionPlan) {
		if ($game) {
			try {
				await $game.deleteProductionPlan(plan.num);
				// trigger reactivity
				$game.player.productionPlans = $game.player.productionPlans;
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
		<a class="cs-link btn btn-sm" href={`/games/${gameId}/production-plans/create`}>Create</a>
	</div>
</Breadcrumb>

{#if $game?.player.productionPlans.length}
	<div class="flex flex-wrap justify-center gap-2">
		{#each $game.player.productionPlans as plan (plan.num)}
			<ProductionPlanCard
				{plan}
				{gameId}
				showDelete={plan.num !== 0}
				on:delete={() => deletePlan(plan)}
			/>
		{/each}
	</div>
{/if}
