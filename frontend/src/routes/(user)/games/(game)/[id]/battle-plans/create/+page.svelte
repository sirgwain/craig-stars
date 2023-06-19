<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { game } from '$lib/services/Context';
	import { BattleAttackWho, BattleTactic, BattleTarget } from '$lib/types/Battle';
	import type { BattlePlan } from '$lib/types/Player';
	import BattlePlanEditor from '../BattlePlanEditor.svelte';

	let gameId = parseInt($page.params.id);

	let plan: BattlePlan = {
		num: -1,
		name: '',
		primaryTarget: BattleTarget.ArmedShips,
		secondaryTarget: BattleTarget.Any,
		tactic: BattleTactic.MaximizeDamageRatio,
		attackWho: BattleAttackWho.EnemiesAndNeutrals,
		dumpCargo: false
	};

	let error = '';

	const onSave = async () => {
		error = '';

		try {
			if (plan && $game) {
				$game.player.battlePlans.push(plan);
				// update this design
				await $game.updatePlayerPlans();
				goto(
					`/games/${gameId}/battle-plans/${
						$game.player.battlePlans[$game.player.battlePlans.length - 1].num
					}`
				);
			}
		} catch (e) {
			error = (e as Error).message;
		}
	};
</script>

<Breadcrumb>
	<svelte:fragment slot="crumbs">
		<li><a href={`/games/${gameId}/battle-plans`}>Battle Plans</a></li>
		<li>{plan?.name ?? '<unknown>'}</li>
	</svelte:fragment>
	<div slot="end" class="flex justify-end mb-1">
		<button class="btn btn-success" type="submit" on:click={(e) => onSave()}>Save</button>
	</div>
</Breadcrumb>

{#if plan}
	<BattlePlanEditor bind:plan readonlyName={plan.num === 0} on:save={(e) => onSave()} bind:error />
{/if}
