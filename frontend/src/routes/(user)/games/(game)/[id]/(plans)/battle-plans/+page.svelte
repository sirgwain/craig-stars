<script lang="ts">
	import Breadcrumb from '#lib/components/game/Breadcrumb.svelte';
	import { addError } from '#lib/services/Errors.js';
	import { getGameContext } from '#lib/services/GameContext.js';
	import type { BattlePlan } from '#lib/types/cs-proto.js';
	import type { ConnectError } from '@connectrpc/connect';
	import BattlePlanCard from './BattlePlanCard.svelte';

	const { game, player, deleteBattlePlan, readOnly } = getGameContext();

	async function deletePlan(plan: BattlePlan) {
		try {
			await deleteBattlePlan(plan.num);
			// trigger reactivity
			$player.playerPlans.battlePlans = $player.playerPlans.battlePlans;
		} catch (e) {
			addError(e as ConnectError);
		}
	}
</script>

<Breadcrumb>
	{#snippet crumbs()}
		<li>Battle Plans</li>
	{/snippet}

	{#snippet end()}
		<div class="flex justify-end mb-1">
			{#if !$readOnly}
				<a class="cs-link btn btn-sm" href={`/games/${$game.id}/battle-plans/create`}>Create</a>
			{/if}
		</div>
	{/snippet}
</Breadcrumb>

{#if $player.playerPlans.battlePlans.length}
	<div class="flex flex-wrap justify-center gap-2">
		{#each $player.playerPlans.battlePlans as plan (plan.num)}
			<BattlePlanCard
				{plan}
				href={`/games/${$game.id}/battle-plans/${plan.num}`}
				showDelete={plan.num !== 0 && !$readOnly}
				onDelete={() => deletePlan(plan)}
			/>
		{/each}
	</div>
{/if}
