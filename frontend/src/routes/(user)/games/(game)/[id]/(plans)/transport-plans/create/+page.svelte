<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import FormError from '$lib/components/FormError.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { game } from '$lib/services/Stores';
	import { CSError, addError } from '$lib/services/Errors';
	import { WaypointTaskTransportAction } from '$lib/types/Fleet';
	import type { TransportPlan } from '$lib/types/Player';
	import TransportPlanEditor from '../TransportPlanEditor.svelte';

	let gameId = parseInt($page.params.id);

	let plan: TransportPlan = {
		num: 0,
		name: '',
		tasks: {
			fuel: {
				action: WaypointTaskTransportAction.None
			},
			ironium: {
				action: WaypointTaskTransportAction.None
			},
			boranium: {
				action: WaypointTaskTransportAction.None
			},
			germanium: {
				action: WaypointTaskTransportAction.None
			},
			colonists: {
				action: WaypointTaskTransportAction.None
			}
		}
	};

	let error = '';

	const onSubmit = async () => {
		error = '';

		try {
			if (plan && $game) {
				// save to server
				await $game.createTransportPlan(plan);
				goto(
					`/games/${gameId}/transport-plans/${
						$game.player.transportPlans[$game.player.transportPlans.length - 1].num
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
