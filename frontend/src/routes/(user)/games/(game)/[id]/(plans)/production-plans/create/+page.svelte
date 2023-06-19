<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import FormError from '$lib/components/FormError.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { game } from '$lib/services/Context';
	import { CSError, addError } from '$lib/services/Errors';
	import type { ProductionPlan } from '$lib/types/Player';
	import ProductionPlanEditor from '../ProductionPlanEditor.svelte';

	let gameId = parseInt($page.params.id);

	let plan: ProductionPlan = {
		num: 0,
		name: '',
		items: []
	};

	let error = '';

	const onSubmit = async () => {
		error = '';

		try {
			if (plan && $game) {
				// save to server
				await $game.createProductionPlan(plan);
				goto(
					`/games/${gameId}/production-plans/${
						$game.player.productionPlans[$game.player.productionPlans.length - 1].num
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
			<li><a href={`/games/${gameId}/production-plans`}>Production Plans</a></li>
			<li>{plan?.name ?? '<unknown>'}</li>
		</svelte:fragment>
		<div slot="end" class="flex justify-end mb-1">
			<button class="btn btn-success" type="submit">Save</button>
		</div>
	</Breadcrumb>

	<FormError {error} />
	{#if plan}
		<ProductionPlanEditor bind:plan />
	{/if}
</form>
