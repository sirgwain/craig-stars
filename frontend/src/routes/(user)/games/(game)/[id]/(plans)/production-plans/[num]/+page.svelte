<script lang="ts">
	import { page } from '$app/stores';
	import FormError from '$lib/components/FormError.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { CSError, addError } from '$lib/services/Errors';
	import ProductionPlanEditor from '../ProductionPlanEditor.svelte';
	import { notify } from '$lib/services/Notifications';
	import type { ProductionPlan } from '$lib/types/Player';

	const { game, player, universe, updateProductionPlan } = getGameContext();
	let num = parseInt($page.params.num);

	let plan: ProductionPlan | undefined = $state();
	$effect(() => {
		plan = $player.productionPlans.find((p) => p.num == num);
	});

	let error = $state('');

	const onSubmit = async () => {
		error = '';

		try {
			if (plan && $game) {
				// save to server
				await updateProductionPlan(plan);
				notify(`Saved ${plan.name}`);
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
			<li><a href={`/games/${$game.id}/production-plans`}>Production Plans</a></li>
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
		<ProductionPlanEditor designFinder={$universe} bind:plan />
	{/if}
</form>
