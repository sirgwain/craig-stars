<script lang="ts">
	import { page } from '$app/state';
	import FormError from '$lib/components/FormError.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { addError } from '$lib/services/Errors';
	import { getGameContext } from '$lib/services/GameContext';
	import { notify } from '$lib/services/Notifications';
	import type { TransportPlan } from '$lib/types/cs-proto';
	import type { ConnectError } from '@connectrpc/connect';
	import TransportPlanEditor from '../TransportPlanEditor.svelte';

	const { game, player, updateTransportPlan } = getGameContext();
	let num = parseInt(page.params.num);

	let plan: TransportPlan | undefined = $derived(
		$player.playerPlans.transportPlans.find((p) => p.num == num)
	);

	let error = $state('');

	const onSubmit = async () => {
		error = '';

		try {
			if (plan) {
				// save to server
				await updateTransportPlan(plan);
				notify(`Saved ${plan.name}`);
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
