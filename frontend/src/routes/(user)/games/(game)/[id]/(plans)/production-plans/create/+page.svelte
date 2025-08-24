<script lang="ts">
	import { goto } from '$app/navigation';
	import FormError from '$lib/components/FormError.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { addError } from '$lib/services/Errors';
	import { getGameContext } from '$lib/services/GameContext';
	import { notify } from '$lib/services/Notifications';
	import { ProductionPlanSchema, type ProductionPlan } from '$lib/types/cs-proto';
	import { create } from '@bufbuild/protobuf';
	import type { ConnectError } from '@connectrpc/connect';
	import ProductionPlanEditor from '../ProductionPlanEditor.svelte';

	const { game, player, universe, createProductionPlan } = getGameContext();

	let plan: ProductionPlan = $state(
		create(ProductionPlanSchema, {
			num: 0,
			name: '',
			items: []
		})
	);

	let error = $state('');

	const onSubmit = async () => {
		error = '';

		try {
			// save to server
			await createProductionPlan(plan);
			notify(`Saved ${plan.name}`);
			goto(
				`/games/${$game.id}/production-plans/${
					$player.playerPlans.productionPlans[$player.playerPlans.productionPlans.length - 1].num
				}`
			);
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
			<li><a href={`/games/${$game.id}/production-plans`}>Production Plans</a></li>
			<li>{plan.name}</li>
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
