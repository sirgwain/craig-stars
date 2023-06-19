<script lang="ts">
	import type { BattlePlan } from '$lib/types/Player';
	import { Trash } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { startCase } from 'lodash-es';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	export let plan: BattlePlan;
	export let gameId: number;
	export let showDelete = true;

	const deletePlan = async (plan: BattlePlan) => {
		if (plan.name != undefined && confirm(`Are you sure you want to delete ${plan.name}?`)) {
			dispatch('delete', { plan });
		}
	};
</script>

<div class="card bg-base-200 shadow-xl rounded-sm border-2 border-base-300 pt-2 m-1 w-full sm:w-[350px]">
	<div class="card-body">
		<h2 class="card-title">
			<a class="cs-link" href={`/games/${gameId}/battle-plans/${plan.num}`}>{plan.name}</a>
		</h2>
		<div class="flex flex-col">
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-28">Name</div>
				<div>{plan.name}</div>
			</div>
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-28">Primary Target</div>
				<div>{startCase(plan.primaryTarget)}</div>
			</div>
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-28">Secondary Target</div>
				<div>{startCase(plan.secondaryTarget)}</div>
			</div>
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-28">Tactic</div>
				<div>{startCase(plan.tactic)}</div>
			</div>
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-28">Attack Who</div>
				<div>{startCase(plan.attackWho)}</div>
			</div>
		</div>
		{#if showDelete}
			<div class="card-actions justify-start">
				<div>
					<button class="btn" on:click={(e) => deletePlan(plan)}>
						<Icon src={Trash} size="24" class="hover:stroke-accent" />
					</button>
				</div>
			</div>
		{/if}
	</div>
</div>
