<script lang="ts">
	import { goto } from '$app/navigation';
	import FormError from '$lib/components/FormError.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { addError } from '$lib/services/Errors';
	import { getGameContext } from '$lib/services/GameContext';
	import { notify } from '$lib/services/Notifications';
	import {
		TransportPlanSchema,
		WaypointTaskTransportAction,
		type TransportPlan
	} from '$lib/types/cs-proto';
	import { create } from '@bufbuild/protobuf';
	import type { ConnectError } from '@connectrpc/connect';
	import TransportPlanEditor from '../TransportPlanEditor.svelte';

	const { game, player, createTransportPlan } = getGameContext();

	let plan: TransportPlan = $state(
		create(TransportPlanSchema, {
			num: 0,
			name: '',
			tasks: {
				fuel: {
					action: WaypointTaskTransportAction.UNSPECIFIED
				},
				ironium: {
					action: WaypointTaskTransportAction.UNSPECIFIED
				},
				boranium: {
					action: WaypointTaskTransportAction.UNSPECIFIED
				},
				germanium: {
					action: WaypointTaskTransportAction.UNSPECIFIED
				},
				colonists: {
					action: WaypointTaskTransportAction.UNSPECIFIED
				}
			}
		})
	);

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
						$player.playerPlans.transportPlans[$player.playerPlans.transportPlans.length - 1].num
					}`
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
