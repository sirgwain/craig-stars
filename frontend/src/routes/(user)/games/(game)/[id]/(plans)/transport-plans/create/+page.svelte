<script lang="ts">
	import { goto } from '$app/navigation';
	import FormError from '$lib/components/FormError.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { CSError, addError } from '$lib/services/Errors';
	import { getGameContext } from '$lib/services/GameContext';
	import { notify } from '$lib/services/Notifications';
	import type { TransportPlan } from '$lib/types/cs';
	import { TransportActionNone } from '$lib/types/cs';
	import TransportPlanEditor from '../TransportPlanEditor.svelte';

	const { game, player, createTransportPlan } = getGameContext();

	let plan: TransportPlan = $state({
		num: 0,
		name: '',
		tasks: {
			fuel: {
				action: TransportActionNone
			},
			ironium: {
				action: TransportActionNone
			},
			boranium: {
				action: TransportActionNone
			},
			germanium: {
				action: TransportActionNone
			},
			colonists: {
				action: TransportActionNone
			}
		}
	});

	let error = $state('');

	const onSubmit = async () => {
		error = '';

		try {
			if (plan && $game) {
				// save to server
				await createTransportPlan(plan);
				notify(`Saved ${plan.name}`);
				goto(
					`/games/${$game.id}/transport-plans/${
						$player.transportPlans[$player.transportPlans.length - 1].num
					}`
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
			<li><a href={`/games/${$game.id}/transport-plans`}>Transport Plans</a></li>
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
		<TransportPlanEditor bind:plan />
	{/if}
</form>
