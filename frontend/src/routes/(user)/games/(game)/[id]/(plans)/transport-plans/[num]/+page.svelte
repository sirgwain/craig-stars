<script lang="ts">
	import { page } from '$app/stores';
	import FormError from '$lib/components/FormError.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import { CSError, addError } from '$lib/services/Errors';
	import TransportPlanEditor from '../TransportPlanEditor.svelte';

	const { game, player, universe } = getGameContext();
	let num = parseInt($page.params.num);

	$: plan = $player.transportPlans.find((p) => p.num == num);

	let error = '';

	const onSubmit = async () => {
		error = '';

		try {
			if (plan && $game) {
				// save to server
				await $game.updateTransportPlan(plan);
			}
		} catch (e) {
			addError(e as CSError);
		}
	};
</script>

<form on:submit|preventDefault={onSubmit}>
	<Breadcrumb>
		<svelte:fragment slot="crumbs">
			<li><a href={`/games/${$game.id}/transport-plans`}>Transport Plans</a></li>
			<li>{plan?.name ?? '<unknown>'}</li>
		</svelte:fragment>
		<div slot="end" class="flex justify-end mb-1">
			<button class="btn btn-success" type="submit">Save</button>
		</div>
	</Breadcrumb>

	<FormError {error} />

	{#if plan}
		<TransportPlanEditor bind:plan />
	{/if}
</form>
