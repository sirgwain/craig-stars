<script lang="ts">
	import { goto } from '$app/navigation';
	import FormError from '$lib/components/FormError.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { CSError, addError } from '$lib/services/Errors';
	import { getGameContext } from '$lib/services/GameContext';
	import { notify } from '$lib/services/Notifications';
	import type { BattlePlan } from '$lib/types/cs';
	import {
		BattleAttackWhoEnemiesAndNeutrals,
		BattleTacticMaximizeDamageRatio,
		BattleTargetAny,
		BattleTargetArmedShips
	} from '$lib/types/cs';
	import BattlePlanEditor from '../BattlePlanEditor.svelte';

	const { game, player, createBattlePlan } = getGameContext();

	let plan: BattlePlan = $state({
		num: 0,
		name: '',
		primaryTarget: BattleTargetArmedShips,
		secondaryTarget: BattleTargetAny,
		tactic: BattleTacticMaximizeDamageRatio,
		attackWho: BattleAttackWhoEnemiesAndNeutrals,
		dumpCargo: false
	});

	let error = $state('');

	const onSubmit = async () => {
		error = '';

		try {
			if (plan && $game) {
				// save to server
				await createBattlePlan(plan);
				notify(`Saved ${plan.name}`);
				goto(
					`/games/${$game.id}/battle-plans/${$player.battlePlans[$player.battlePlans.length - 1].num}`
				);
			}
		} catch (e) {
			addError(e as CSError);
		}
	};
</script>

<form
	onsubmit={(e) => {
		e.preventDefault();
		onSubmit();
	}}
>
	<Breadcrumb>
		{#snippet crumbs()}
			<li><a href={`/games/${$game.id}/battle-plans`}>Battle Plans</a></li>
			<li>{plan?.name ?? '<unknown>'}</li>
		{/snippet}
		{#snippet end()}
			<div class="flex justify-end mb-1">
				<button class="btn btn-success mx-1" type="submit">Save</button>
			</div>
		{/snippet}
	</Breadcrumb>

	<FormError {error} />

	{#if plan}
		<BattlePlanEditor bind:plan />
	{/if}
</form>
