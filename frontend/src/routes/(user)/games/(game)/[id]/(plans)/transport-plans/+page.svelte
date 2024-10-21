<script lang="ts">
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { addError, type CSError } from '$lib/services/Errors';
	import type { TransportPlan } from '$lib/types/Player';
	import TransportPlanCard from './TransportPlanCard.svelte';
	import { notify } from '$lib/services/Notifications';

	const { game, player, deleteTransportPlan } = getGameContext();

	async function deletePlan(plan: TransportPlan) {
		if ($game) {
			try {
				await deleteTransportPlan(plan.num);
				// trigger reactivity
				$player.transportPlans = $player.transportPlans;
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

	<div slot="end" class="flex justify-end mb-1">
		<a class="cs-link btn btn-sm" href={`/games/${$game.id}/transport-plans/create`}>Create</a>
	</div>
</Breadcrumb>

{#if $player.transportPlans.length}
	<div class="flex flex-wrap justify-center gap-2">
		{#each $player.transportPlans as plan (plan.num)}
			<TransportPlanCard
				{plan}
				href={`/games/${$game.id}/transport-plans/${plan.num}`}
				showDelete={plan.num !== 0}
				on:delete={() => deletePlan(plan)}
			/>
		{/each}
	</div>
{/if}
