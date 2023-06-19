<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { game } from '$lib/services/Context';
	import { CSError, addError } from '$lib/services/Errors';
	import BattlePlanEditor from '../BattlePlanEditor.svelte';

	let gameId = parseInt($page.params.id);
	let num = parseInt($page.params.num);

	$: plan = $game?.player.battlePlans.find((p) => p.num == num);

	let error = '';

	const onSave = async () => {
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

<Breadcrumb>
	<svelte:fragment slot="crumbs">
		<li><a href={`/games/${gameId}/battle-plans`}>Battle Plans</a></li>
		<li>{plan?.name ?? '<unknown>'}</li>
	</svelte:fragment>
	<div slot="end" class="flex justify-end mb-1">
		<button class="btn btn-success" type="submit" on:click={(e) => onSave()}>Save</button>
	</div>
</Breadcrumb>

{#if plan}
	<BattlePlanEditor bind:plan on:save={(e) => onSave()} bind:error />
{/if}
