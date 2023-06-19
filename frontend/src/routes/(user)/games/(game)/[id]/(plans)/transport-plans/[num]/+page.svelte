<script lang="ts">
	import { page } from '$app/stores';
	import FormError from '$lib/components/FormError.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { game } from '$lib/services/Context';
	import { CSError, addError } from '$lib/services/Errors';
	import TransportPlanEditor from '../TransportPlanEditor.svelte';

	let gameId = parseInt($page.params.id);
	let num = parseInt($page.params.num);

	$: plan = $game?.player.transportPlans.find((p) => p.num == num);

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
			<li><a href={`/games/${gameId}/transport-plans`}>Transport Plans</a></li>
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
