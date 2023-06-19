<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { game } from '$lib/services/Context';
	import { addError, type CSError } from '$lib/services/Errors';
	import type { TransportPlan } from '$lib/types/Player';
	import TransportPlanCard from './TransportPlanCard.svelte';

	let gameId = parseInt($page.params.id);

	async function deletePlan(plan: TransportPlan) {
		if ($game) {
			try {
				await $game.deleteTransportPlan(plan.num);
				// trigger reactivity
				$game.player.transportPlans = $game.player.transportPlans;
			} catch (e) {
				addError(e as CSError);
			}
		}
	}
</script>

<Breadcrumb>
	<svelte:fragment slot="crumbs">
		<li>Transport Plans</li>
	</svelte:fragment>

	<div slot="end">
		<a class="cs-link btn btn-sm" href={`/games/${gameId}/transport-plans/create`}>Create</a>
	</div>
</Breadcrumb>

{#if $game?.player.transportPlans.length}
	<div class="flex flex-wrap justify-center gap-2">
		{#each $game.player.transportPlans as plan (plan.num)}
			<TransportPlanCard
				{plan}
				{gameId}
				showDelete={plan.num !== 0}
				on:delete={() => deletePlan(plan)}
			/>
		{/each}
	</div>
{/if}
