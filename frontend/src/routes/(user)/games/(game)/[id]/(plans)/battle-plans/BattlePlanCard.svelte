<script lang="ts">
	import {
		BattleAttackWho,
		BattleTactic,
		BattleTarget,
		type BattlePlan
	} from '$lib/types/cs-proto';
	import { enumToString } from '$lib/types/Enums';
	import { Trash } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	type Props = {
		plan: BattlePlan;
		href: string;
		showDelete?: boolean;
		onDelete?: (plan: BattlePlan) => void;
	};

	let { plan, href, showDelete = true, onDelete }: Props = $props();

	const deletePlan = async (plan: BattlePlan) => {
		if (confirm(`Are you sure you want to delete ${plan.name}?`)) {
			onDelete?.(plan);
		}
	};
</script>

<div
	class="card bg-base-200 shadow rounded-sm border-2 border-base-300 pt-2 m-1 w-full sm:w-[350px]"
>
	<div class="card-body">
		<h2 class="card-title">
			<a class="cs-link" {href}>{plan.name}</a>
		</h2>
		<div class="flex flex-col">
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-28">Name</div>
				<div>{plan.name}</div>
			</div>
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-28">Primary Target</div>
				<div>{enumToString(BattleTarget, plan.primaryTarget)}</div>
			</div>
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-28">Secondary Target</div>
				<div>{enumToString(BattleTarget, plan.secondaryTarget)}</div>
			</div>
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-28">Tactic</div>
				<div>{enumToString(BattleTactic, plan.tactic)}</div>
			</div>
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-28">Attack Who</div>
				<div>{enumToString(BattleAttackWho, plan.attackWho)}</div>
			</div>
		</div>
		{#if showDelete}
			<div class="card-actions justify-start">
				<div>
					<button
						class="btn"
						onclick={() => deletePlan(plan)}
						data-type="delete-button"
						data-id={`${plan.name}`}
					>
						<Icon src={Trash} size="24" class="hover:stroke-accent" />
					</button>
				</div>
			</div>
		{/if}
	</div>
</div>
