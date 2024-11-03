<script lang="ts">
	import { preventDefault } from 'svelte/legacy';

	import { goto } from '$app/navigation';
	import FormError from '$lib/components/FormError.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { CSError, addError } from '$lib/services/Errors';
	import type { ProductionPlan } from '$lib/types/Player';
	import ProductionPlanEditor from '../ProductionPlanEditor.svelte';
	import { notify } from '$lib/services/Notifications';

	const { game, player, universe, createProductionPlan } = getGameContext();

	let plan: ProductionPlan = $state({
		num: 0,
		name: '',
		items: []
	});

	let error = $state('');

	const onSubmit = async () => {
		error = '';

		try {
			if (plan && $game) {
				// save to server
				await createProductionPlan(plan);
				notify(`Saved ${plan.name}`);
				goto(
					`/games/${$game.id}/production-plans/${
						$player.productionPlans[$player.productionPlans.length - 1].num
					}`
				);
			}
		} catch (e) {
			addError(e as CSError);
		}
	};
</script>

<form onsubmit={preventDefault(onSubmit)}>
	<Breadcrumb>
		{#snippet crumbs()}
			
				<li><a href={`/games/${$game.id}/production-plans`}>Production Plans</a></li>
				<li>{plan?.name ?? '<unknown>'}</li>
			
			{/snippet}
		{#snippet end()}
				<div  class="flex justify-end mb-1">
				<button class="btn btn-success mx-1" type="submit">Save</button>
			</div>
			{/snippet}
	</Breadcrumb>

	<FormError {error} />
	{#if plan}
		<ProductionPlanEditor designFinder={$universe} bind:plan />
	{/if}
</form>
