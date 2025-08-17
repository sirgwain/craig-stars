<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import {
		VictoryConditionAttainTechLevels,
		VictoryConditionExceedsScore,
		VictoryConditionExceedsSecondPlaceScore,
		VictoryConditionHighestScoreAfterYears,
		VictoryConditionOwnCapitalShips,
		VictoryConditionOwnPlanets,
		VictoryConditionProductionCapacity
	} from '$lib/types/Consts';
	import { CheckBadge } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	const { game, universe } = getGameContext();

	let players = $derived($universe.playerIntels);
</script>

<table class="table table-zebra table-fixed mx-auto w-full sm:w-auto">
	<thead>
		<tr>
			<th></th>
			{#each players as player (player.num)}
				<th class="h-20 w-20"
					><div class="py-4 -rotate-45">{$universe.getPlayerPluralName(player.num)}</div></th
				>
			{/each}
		</tr>
	</thead>
	<tbody>
		{#if ($game.victoryConditions.conditions & VictoryConditionOwnPlanets) > 0}
			<tr>
				<td
					>Owns {Math.ceil(
						($game.victoryConditions.ownPlanets / 100.0) * $universe.planetIntels.length
					).toFixed()}/{$universe.planets.length} planets.</td
				>
				{#each players as player (player.num)}
					<td>
						{#if (($universe.getPlayerScore(player.num)?.achievedVictoryConditions ?? 0) & VictoryConditionOwnPlanets) > 0}
							<Icon src={CheckBadge} size="24" class="stroke-success" />
						{/if}
					</td>{/each}
			</tr>
		{/if}
		{#if ($game.victoryConditions.conditions & VictoryConditionAttainTechLevels) > 0}
			<tr>
				<td>
					Attains Tech {$game.victoryConditions.attainTechLevel} in {$game.victoryConditions
						.attainTechLevelNumFields} fields.
				</td>
				{#each players as player (player.num)}
					<td>
						{#if (($universe.getPlayerScore(player.num)?.achievedVictoryConditions ?? 0) & VictoryConditionAttainTechLevels) > 0}
							<Icon src={CheckBadge} size="24" class="stroke-success" />
						{/if}
					</td>{/each}
			</tr>
		{/if}
		{#if ($game.victoryConditions.conditions & VictoryConditionExceedsScore) > 0}
			<tr>
				<td>
					Exceeds a score of {$game.victoryConditions.exceedsScore}.
				</td>
				{#each players as player (player.num)}
					<td>
						{#if (($universe.getPlayerScore(player.num)?.achievedVictoryConditions ?? 0) & VictoryConditionExceedsScore) > 0}
							<Icon src={CheckBadge} size="24" class="stroke-success" />
						{/if}
					</td>{/each}
			</tr>
		{/if}
		{#if ($game.victoryConditions.conditions & VictoryConditionExceedsSecondPlaceScore) > 0}
			<tr>
				<td>
					Exceeds second place score by {$game.victoryConditions.exceedsSecondPlaceScore}%.
				</td>
				{#each players as player (player.num)}
					<td>
						{#if (($universe.getPlayerScore(player.num)?.achievedVictoryConditions ?? 0) & VictoryConditionExceedsSecondPlaceScore) > 0}
							<Icon src={CheckBadge} size="24" class="stroke-success" />
						{/if}
					</td>{/each}
			</tr>
		{/if}
		{#if ($game.victoryConditions.conditions & VictoryConditionProductionCapacity) > 0}
			<tr>
				<td>
					Has a production capacity of {$game.victoryConditions.productionCapacity},000
					resources/yr.
				</td>
				{#each players as player (player.num)}
					<td>
						{#if (($universe.getPlayerScore(player.num)?.achievedVictoryConditions ?? 0) & VictoryConditionProductionCapacity) > 0}
							<Icon src={CheckBadge} size="24" class="stroke-success" />
						{/if}
					</td>{/each}
			</tr>
		{/if}
		{#if ($game.victoryConditions.conditions & VictoryConditionOwnCapitalShips) > 0}
			<tr>
				<td>
					Owns {$game.victoryConditions.ownCapitalShips} capital ships.
				</td>
				{#each players as player (player.num)}
					<td>
						{#if (($universe.getPlayerScore(player.num)?.achievedVictoryConditions ?? 0) & VictoryConditionOwnCapitalShips) > 0}
							<Icon src={CheckBadge} size="24" class="stroke-success" />
						{/if}
					</td>{/each}
			</tr>
		{/if}
		{#if ($game.victoryConditions.conditions & VictoryConditionHighestScoreAfterYears) > 0}
			<tr>
				<td>
					Has the highest score after {$game.victoryConditions.highestScoreAfterYears} years.
				</td>
				{#each players as player (player.num)}
					<td>
						{#if (($universe.getPlayerScore(player.num)?.achievedVictoryConditions ?? 0) & VictoryConditionHighestScoreAfterYears) > 0}
							<Icon src={CheckBadge} size="24" class="stroke-success" />
						{/if}
					</td>{/each}
			</tr>
		{/if}
	</tbody>
</table>

<div class="text-center">
	Winner must meet {$game.victoryConditions.numCriteriaRequired} of the above criteria after at least
	{$game.victoryConditions.yearsPassed} years have passed.
</div>
