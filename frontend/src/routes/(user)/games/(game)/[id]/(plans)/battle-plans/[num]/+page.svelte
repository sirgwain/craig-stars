<script lang="ts">
	import { page } from '$app/stores';
	import FormError from '$lib/components/FormError.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { game } from '$lib/services/Stores';
	import { CSError, addError } from '$lib/services/Errors';
	import BattlePlanEditor from '../BattlePlanEditor.svelte';

	let gameId = parseInt($page.params.id);
	let num = parseInt($page.params.num);

	$: plan = $game?.player.battlePlans.find((p) => p.num == num);

	let error = '';

	const onSubmit = async () => {
		error = '';

		try {
			if (plan && $game) {
				// save to server
				await $game.updateBattlePlan(plan);
			}
		} catch (e) {
			addError(e as CSError);
		}
	};
</script>

<form on:submit|preventDefault={onSubmit}>
	<Breadcrumb>
		<svelte:fragment slot="crumbs">
			<li><a href={`/games/${gameId}/battle-plans`}>Battle Plans</a></li>
			<li>{plan?.name ?? '<unknown>'}</li>
		</svelte:fragment>
		<div slot="end" class="flex justify-end mb-1">
			<button class="btn btn-success" type="submit" on:click={(e) => onSubmit()}>Save</button>
		</div>
	</Breadcrumb>

	<FormError {error} />

	{#if plan}
		<BattlePlanEditor bind:plan />
	{/if}
</form>
