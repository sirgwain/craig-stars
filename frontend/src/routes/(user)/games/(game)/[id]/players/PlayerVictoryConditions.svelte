<script lang="ts">
	import { getGameContext } from '#lib/services/GameContext.js';
	import {
		VictoryConditionAttainTechLevels,
		VictoryConditionExceedsScore,
		VictoryConditionExceedsSecondPlaceScore,
		VictoryConditionHighestScoreAfterYears,
		VictoryConditionOwnCapitalShips,
		VictoryConditionOwnPlanets,
		VictoryConditionProductionCapacity
	} from '#lib/types/Consts.js';
	import { CheckBadge } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	const { game, universe } = getGameContext();

	let players = $derived($universe.playerIntels);
</script>

<!-- Names run up the page on phones so they can't overlap. The condition column keeps a readable
     width and stays put while the players scroll sideways. -->
<div class="overflow-x-auto">
	<table class="table table-zebra table-pin-cols mx-auto w-auto">
		<thead>
			<tr>
				<th></th>
				{#each players as player (player.num)}
					<td class="px-2 align-bottom text-center sm:px-4">
						<div
							class="mx-auto max-h-28 rotate-180 truncate [writing-mode:vertical-rl] sm:max-h-none sm:max-w-28 sm:rotate-0 sm:[writing-mode:horizontal-tb]"
							title={$universe.getPlayerPluralName(player.num)}
						>
							{$universe.getPlayerPluralName(player.num)}
						</div>
					</td>
				{/each}
			</tr>
		</thead>
		<tbody>
			{#if ($game.victoryConditions.conditions & VictoryConditionOwnPlanets) > 0}
				<tr>
					<th class="font-normal min-w-36 max-w-56 sm:max-w-none"
						>Owns {Math.ceil(
							($game.victoryConditions.ownPlanets / 100.0) * $universe.planetIntels.length
						).toFixed()}/{$universe.planets.length} planets.</th
					>
					{#each players as player (player.num)}
						<td class="px-2 text-center sm:px-4">
							{#if (($universe.getPlayerScore(player.num)?.achievedVictoryConditions ?? 0) & VictoryConditionOwnPlanets) > 0}
								<Icon src={CheckBadge} size="24" class="stroke-success" />
							{/if}
						</td>{/each}
				</tr>
			{/if}
			{#if ($game.victoryConditions.conditions & VictoryConditionAttainTechLevels) > 0}
				<tr>
					<th class="font-normal min-w-36 max-w-56 sm:max-w-none">
						Attains Tech {$game.victoryConditions.attainTechLevel} in {$game.victoryConditions
							.attainTechLevelNumFields} fields.
					</th>
					{#each players as player (player.num)}
						<td class="px-2 text-center sm:px-4">
							{#if (($universe.getPlayerScore(player.num)?.achievedVictoryConditions ?? 0) & VictoryConditionAttainTechLevels) > 0}
								<Icon src={CheckBadge} size="24" class="stroke-success" />
							{/if}
						</td>{/each}
				</tr>
			{/if}
			{#if ($game.victoryConditions.conditions & VictoryConditionExceedsScore) > 0}
				<tr>
					<th class="font-normal min-w-36 max-w-56 sm:max-w-none">
						Exceeds a score of {$game.victoryConditions.exceedsScore}.
					</th>
					{#each players as player (player.num)}
						<td class="px-2 text-center sm:px-4">
							{#if (($universe.getPlayerScore(player.num)?.achievedVictoryConditions ?? 0) & VictoryConditionExceedsScore) > 0}
								<Icon src={CheckBadge} size="24" class="stroke-success" />
							{/if}
						</td>{/each}
				</tr>
			{/if}
			{#if ($game.victoryConditions.conditions & VictoryConditionExceedsSecondPlaceScore) > 0}
				<tr>
					<th class="font-normal min-w-36 max-w-56 sm:max-w-none">
						Exceeds second place score by {$game.victoryConditions.exceedsSecondPlaceScore}%.
					</th>
					{#each players as player (player.num)}
						<td class="px-2 text-center sm:px-4">
							{#if (($universe.getPlayerScore(player.num)?.achievedVictoryConditions ?? 0) & VictoryConditionExceedsSecondPlaceScore) > 0}
								<Icon src={CheckBadge} size="24" class="stroke-success" />
							{/if}
						</td>{/each}
				</tr>
			{/if}
			{#if ($game.victoryConditions.conditions & VictoryConditionProductionCapacity) > 0}
				<tr>
					<th class="font-normal min-w-36 max-w-56 sm:max-w-none">
						Has a production capacity of {$game.victoryConditions.productionCapacity},000
						resources/yr.
					</th>
					{#each players as player (player.num)}
						<td class="px-2 text-center sm:px-4">
							{#if (($universe.getPlayerScore(player.num)?.achievedVictoryConditions ?? 0) & VictoryConditionProductionCapacity) > 0}
								<Icon src={CheckBadge} size="24" class="stroke-success" />
							{/if}
						</td>{/each}
				</tr>
			{/if}
			{#if ($game.victoryConditions.conditions & VictoryConditionOwnCapitalShips) > 0}
				<tr>
					<th class="font-normal min-w-36 max-w-56 sm:max-w-none">
						Owns {$game.victoryConditions.ownCapitalShips} capital ships.
					</th>
					{#each players as player (player.num)}
						<td class="px-2 text-center sm:px-4">
							{#if (($universe.getPlayerScore(player.num)?.achievedVictoryConditions ?? 0) & VictoryConditionOwnCapitalShips) > 0}
								<Icon src={CheckBadge} size="24" class="stroke-success" />
							{/if}
						</td>{/each}
				</tr>
			{/if}
			{#if ($game.victoryConditions.conditions & VictoryConditionHighestScoreAfterYears) > 0}
				<tr>
					<th class="font-normal min-w-36 max-w-56 sm:max-w-none">
						Has the highest score after {$game.victoryConditions.highestScoreAfterYears} years.
					</th>
					{#each players as player (player.num)}
						<td class="px-2 text-center sm:px-4">
							{#if (($universe.getPlayerScore(player.num)?.achievedVictoryConditions ?? 0) & VictoryConditionHighestScoreAfterYears) > 0}
								<Icon src={CheckBadge} size="24" class="stroke-success" />
							{/if}
						</td>{/each}
				</tr>
			{/if}
		</tbody>
	</table>
</div>

<div class="text-center">
	Winner must meet {$game.victoryConditions.numCriteriaRequired} of the above criteria after at least
	{$game.victoryConditions.yearsPassed} years have passed.
</div>
