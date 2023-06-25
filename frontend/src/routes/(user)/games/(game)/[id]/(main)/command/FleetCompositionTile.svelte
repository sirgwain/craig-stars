<script lang="ts">
	import { EventManager } from '$lib/EventManager';
	import { onShipDesignTooltip } from '$lib/components/game/tooltips/ShipDesignTooltip.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import { selectedWaypoint } from '$lib/services/Stores';
	import type { CommandedFleet } from '$lib/types/Fleet';
	import { createEventDispatcher } from 'svelte';
	import CommandTile from './CommandTile.svelte';

	const dispatch = createEventDispatcher();
	const { player, universe } = getGameContext();

	export let fleet: CommandedFleet;

	const split = () => {
		EventManager.publishSplitFleetDialogRequestedEvent(fleet);
	};

	const splitAll = async () => {
		dispatch('splitAll');
	};
	const merge = () => {
		EventManager.publishMergeFleetDialogRequestedEvent(fleet);
	};
</script>

{#if fleet.waypoints && $selectedWaypoint}
	<CommandTile title="Fleet Composition">
		<div class="bg-base-100 h-20 overflow-y-auto">
			<ul class="w-full h-full">
				{#each fleet.tokens as token, index}
					<li class="pl-1">
						<button
							type="button"
							class="w-full cursor-help"
							on:pointerdown|preventDefault={(e) =>
								onShipDesignTooltip(e, $universe.getDesign($player.num, token.designNum))}
						>
							<div class="flex flex-row justify-between">
								<div>
									{$universe.getDesign($player.num, token.designNum)?.name}
								</div>
								<div>
									{token.quantity}
								</div>
							</div>
						</button>
					</li>
				{/each}
			</ul>
		</div>
		<div class="flex justify-between my-1">
			<div class="my-auto">Battle Plan:</div>
			<div>
				<select
					class="select select-outline select-secondary select-sm text-sm"
					name="battlePlan"
					bind:value={fleet.battlePlanNum}
				>
					{#each $player.battlePlans as battlePlan}
						<option value={battlePlan.num}>{battlePlan.name}</option>
					{/each}
				</select>
			</div>
		</div>
		<div class="flex justify-between my-1">
			<div>Est Range:</div>
			<div>{fleet.spec.estimatedRange ? `${fleet.spec.estimatedRange} l.y.` : '--'}</div>
		</div>
		<div class="flex justify-between my-1">
			<div>Percent Cloaked</div>
			<div>{fleet.spec.cloakPercent ? fleet.spec.cloakPercent * 100 + '%' : 'none'}</div>
		</div>
		<div class="flex justify-between">
			<button on:click={split} class="btn btn-outline btn-sm normal-case btn-secondary"
				>Split</button
			>
			<button on:click={splitAll} class="btn btn-outline btn-sm normal-case btn-secondary"
				>Split All</button
			>
			<button on:click={merge} class="btn btn-outline btn-sm normal-case btn-secondary"
				>Merge</button
			>
		</div>
	</CommandTile>
{/if}
