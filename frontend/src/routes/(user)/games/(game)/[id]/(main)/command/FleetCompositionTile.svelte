<script lang="ts">
	import { onShipDesignTooltip } from '$lib/components/game/tooltips/ShipDesignTooltip.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import { selectedWaypoint } from '$lib/services/Stores';
	import { getDamagePercentForToken, type CommandedFleet } from '$lib/types/Fleet';
	import { createEventDispatcher } from 'svelte';
	import CommandTile from './CommandTile.svelte';
	import type { SplitFleetEvent } from '../../dialogs/split/SplitFleetDialog.svelte';
	import type { MergeFleetsEvent } from '../../dialogs/merge/MergeFleetsDialog.svelte';
	import { Infinite } from '$lib/types/MapObject';

	const dispatchSplit = createEventDispatcher<SplitFleetEvent>();
	const dispatchMerge = createEventDispatcher<MergeFleetsEvent>();
	const { game, player, universe } = getGameContext();

	export let fleet: CommandedFleet;

	const split = () => {
		dispatchSplit('split-fleet-dialog', { fleet });
	};

	const splitAll = async () => {
		dispatchSplit('split-all');
	};
	const merge = () => {
		dispatchMerge('merge-fleets-dialog', {
			fleet,
			otherFleetsHere: $universe.getMyFleetsByPosition(fleet).filter((f) => f.num !== fleet.num)
		});
	};

	const updateBattlePlan = async (num: number) => {
		fleet.battlePlanNum = num;
		await $game.updateFleetOrders(fleet);
	};
</script>

{#if fleet.waypoints && $selectedWaypoint}
	<CommandTile title="Fleet Composition">
		<div class="bg-base-100 h-20 overflow-y-auto">
			<ul class="w-full h-full">
				{#each fleet.tokens as token}
					<li class="pl-1">
						<button
							type="button"
							class="w-full cursor-help"
							on:pointerdown|preventDefault={(e) =>
								onShipDesignTooltip(e, $universe.getDesign($player.num, token.designNum))}
						>
							<div class="flex flex-row justify-between relative">
								{#if (token.damage ?? 0) > 0 && (token.quantityDamaged ?? 0) > 0}
									<div
										style={`width: ${getDamagePercentForToken(
											token,
											$universe.getMyDesign(token.designNum)
										).toFixed()}%`}
										class="damage-bar h-full absolute opacity-50"
									/>
								{/if}
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
					on:change={(e) => updateBattlePlan(parseInt(e.currentTarget.value))}
				>
					{#each $player.battlePlans as battlePlan}
						<option value={battlePlan.num}>{battlePlan.name}</option>
					{/each}
				</select>
			</div>
		</div>
		<div class="flex justify-between my-1">
			<div>Est Range:</div>
			<div>
				{fleet.spec.estimatedRange
					? fleet.spec.estimatedRange === Infinite
						? 'Infinite'
						: `${fleet.spec.estimatedRange} l.y.`
					: '--'}
			</div>
		</div>
		<div class="flex justify-between my-1">
			<div>Percent Cloaked</div>
			<div>{fleet.spec.cloakPercent ? fleet.spec.cloakPercent + '%' : 'none'}</div>
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
