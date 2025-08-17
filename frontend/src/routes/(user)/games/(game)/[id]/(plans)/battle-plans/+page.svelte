<script lang="ts">
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { addError } from '$lib/services/Errors';
	import { getGameContext } from '$lib/services/GameContext';
	import type { BattlePlan } from '$lib/types/cs-proto';
	import type { ConnectError } from '@connectrpc/connect';
	import BattlePlanCard from './BattlePlanCard.svelte';

	const { game, player, deleteBattlePlan } = getGameContext();

	async function deletePlan(plan: BattlePlan) {
		if ($game) {
			try {
				await deleteBattlePlan(plan.num);
				// trigger reactivity
				$player.playerPlans.battlePlans = $player.playerPlans.battlePlans;
			} catch (e) {
				addError(e as ConnectError);
			}
		}
	}
</script>

<Breadcrumb>
	{#snippet crumbs()}
		<li>Battle Plans</li>
	{/snippet}

	{#snippet end()}
		<div class="flex justify-end mb-1">
			<a class="cs-link btn btn-sm" href={`/games/${$game.id}/battle-plans/create`}>Create</a>
		</div>
	{/snippet}
</Breadcrumb>

{#if $player.playerPlans.battlePlans.length}
	<div class="flex flex-wrap justify-center gap-2">
		{#each $player.playerPlans.battlePlans as plan (plan.num)}
			<BattlePlanCard
				{plan}
				href={`/games/${$game.id}/battle-plans/${plan.num}`}
				showDelete={plan.num !== 0}
				onDelete={() => deletePlan(plan)}
			/>
		{/each}
	</div>
{/if}
