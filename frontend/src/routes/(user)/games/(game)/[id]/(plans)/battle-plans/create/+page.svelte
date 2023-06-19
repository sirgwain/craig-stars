<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import FormError from '$lib/components/FormError.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { game } from '$lib/services/Stores';
	import { CSError, addError } from '$lib/services/Errors';
	import { BattleAttackWho, BattleTactic, BattleTarget } from '$lib/types/Battle';
	import type { BattlePlan } from '$lib/types/Player';
	import BattlePlanEditor from '../BattlePlanEditor.svelte';

	let gameId = parseInt($page.params.id);

	let plan: BattlePlan = {
		num: 0,
		name: '',
		primaryTarget: BattleTarget.ArmedShips,
		secondaryTarget: BattleTarget.Any,
		tactic: BattleTactic.MaximizeDamageRatio,
		attackWho: BattleAttackWho.EnemiesAndNeutrals,
		dumpCargo: false
	};

	let error = '';

	const onSubmit = async () => {
		error = '';

		try {
			if (plan && $game) {
				// save to server
				await $game.createBattlePlan(plan);
				goto(
					`/games/${gameId}/battle-plans/${
						$game.player.battlePlans[$game.player.battlePlans.length - 1].num
					}`
				);
			}
		} catch (e) {
			addError(e as CSError);
		}
	};
</script>

<form on:submit|preventDefault={onSubmit}>
	<Breadcrumb>
		<svelte:fragment slot="crumbs">
			<li><a href={`/games/${gameId}/battle-plans`}>Battle Plans</a></li>
			<li>{plan?.name ?? '<unknown>'}</li>
		</svelte:fragment>
		<div slot="end" class="flex justify-end mb-1">
			<button class="btn btn-success" type="submit" on:click={(e) => onSubmit()}>Save</button>
		</div>
	</Breadcrumb>

	<FormError {error} />

	{#if plan}
		<BattlePlanEditor bind:plan />
	{/if}
</form>
