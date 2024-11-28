<script lang="ts">
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { addError, type CSError } from '$lib/services/Errors';
	import { getGameContext } from '$lib/services/GameContext';
	import type { TransportPlan } from '$lib/types/Player';
	import TransportPlanCard from './TransportPlanCard.svelte';

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
	{#snippet crumbs()}
		<li>Transport Plans</li>
	{/snippet}

	{#snippet end()}
		<div class="flex justify-end mb-1">
			<a class="cs-link btn btn-sm" href={`/games/${$game.id}/transport-plans/create`}>Create</a>
		</div>
	{/snippet}
</Breadcrumb>

{#if $player.transportPlans.length}
	<div class="flex flex-wrap justify-center gap-2">
		{#each $player.transportPlans as plan (plan.num)}
			<TransportPlanCard
				{plan}
				href={`/games/${$game.id}/transport-plans/${plan.num}`}
				showDelete={plan.num !== 0}
				onDelete={() => deletePlan(plan)}
			/>
		{/each}
	</div>
{/if}
