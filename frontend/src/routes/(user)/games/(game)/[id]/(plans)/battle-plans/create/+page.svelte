<script lang="ts">
	import { goto } from '$app/navigation';
	import FormError from '$lib/components/FormError.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { addError } from '$lib/services/Errors';
	import { getGameContext } from '$lib/services/GameContext';
	import { notify } from '$lib/services/Notifications';
	import {
		BattleAttackWho,
		BattlePlanSchema,
		BattleTactic,
		BattleTarget,
		type BattlePlan
	} from '$lib/types/cs-proto';
	import { create } from '@bufbuild/protobuf';
	import type { ConnectError } from '@connectrpc/connect';
	import BattlePlanEditor from '../BattlePlanEditor.svelte';

	const { game, player, createBattlePlan } = getGameContext();

	let plan: BattlePlan = $state(
		create(BattlePlanSchema, {
			num: 0,
			name: '',
			primaryTarget: BattleTarget.ARMED_SHIPS,
			secondaryTarget: BattleTarget.ANY,
			tactic: BattleTactic.MAXIMIZE_DAMAGE_RATIO,
			attackWho: BattleAttackWho.ENEMIES_AND_NEUTRALS,
			dumpCargo: false
		})
	);

	let error = $state('');

	const onSubmit = async () => {
		error = '';

		try {
			if (plan && $game) {
				// save to server
				await createBattlePlan(plan);
				notify(`Saved ${plan.name}`);
				goto(
					`/games/${$game.id}/battle-plans/${$player.playerPlans.battlePlans[$player.playerPlans.battlePlans.length - 1].num}`
				);
			}
		} catch (e) {
			addError(e as ConnectError);
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
