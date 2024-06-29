<script lang="ts">
	import VictoryConditions from '$lib/components/game/newgame/VictoryConditions.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { VictoryCondition } from '$lib/types/Game';
	import { CheckBadge } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	const { game, player, universe } = getGameContext();

	$: players = $universe.players;
</script>

<table class="table table-zebra table-fixed mx-auto w-full sm:w-auto">
	<thead>
		<th />
		{#each players as player}
			<th class="h-20 w-20"
				><div class="py-4 -rotate-45">{$universe.getPlayerName(player.num)}</div></th
			>
		{/each}
	</thead>
	<tbody>
		{#if ($game.victoryConditions.conditions & VictoryCondition.OwnPlanets) > 0}
			<tr>
				<td
					>Own {Math.ceil(
						($game.victoryConditions.ownPlanets / 100.0) * $universe.planets.length
					).toFixed()} of {$universe.planets.length} planets</td
				>
				{#each players as player}
					<td>
						{#if (($universe.getPlayerScore(player.num)?.achievedVictoryConditions ?? 0) & VictoryCondition.OwnPlanets) > 0}
							<Icon src={CheckBadge} size="24" class="stroke-success" />
						{/if}
					</td>{/each}
			</tr>
		{/if}
		{#if ($game.victoryConditions.conditions & VictoryCondition.AttainTechLevels) > 0}
			<tr>
				<td>
					Attain Tech {$game.victoryConditions.attainTechLevel} in {$game.victoryConditions
						.attainTechLevelNumFields} fields.
				</td>
				{#each players as player}
					<td>
						{#if (($universe.getPlayerScore(player.num)?.achievedVictoryConditions ?? 0) & VictoryCondition.AttainTechLevels) > 0}
							<Icon src={CheckBadge} size="24" class="stroke-success" />
						{/if}
					</td>{/each}
			</tr>
		{/if}
		{#if ($game.victoryConditions.conditions & VictoryCondition.ExceedsScore) > 0}
			<tr>
				<td>
					Exceeds a score of {$game.victoryConditions.exceedsScore}.
				</td>
				{#each players as player}
					<td>
						{#if (($universe.getPlayerScore(player.num)?.achievedVictoryConditions ?? 0) & VictoryCondition.ExceedsScore) > 0}
							<Icon src={CheckBadge} size="24" class="stroke-success" />
						{/if}
					</td>{/each}
			</tr>
		{/if}
		{#if ($game.victoryConditions.conditions & VictoryCondition.ExceedsSecondPlaceScore) > 0}
			<tr>
				<td>
					Exceeds a second place score by {$game.victoryConditions.exceedsSecondPlaceScore}%
				</td>
				{#each players as player}
					<td>
						{#if (($universe.getPlayerScore(player.num)?.achievedVictoryConditions ?? 0) & VictoryCondition.ExceedsSecondPlaceScore) > 0}
							<Icon src={CheckBadge} size="24" class="stroke-success" />
						{/if}
					</td>{/each}
			</tr>
		{/if}
		{#if ($game.victoryConditions.conditions & VictoryCondition.ProductionCapacity) > 0}
			<tr>
				<td>
					Has a production capacity of {$game.victoryConditions.productionCapacity} thousand.
				</td>
				{#each players as player}
					<td>
						{#if (($universe.getPlayerScore(player.num)?.achievedVictoryConditions ?? 0) & VictoryCondition.ProductionCapacity) > 0}
							<Icon src={CheckBadge} size="24" class="stroke-success" />
						{/if}
					</td>{/each}
			</tr>
		{/if}
		{#if ($game.victoryConditions.conditions & VictoryCondition.OwnCapitalShips) > 0}
			<tr>
				<td>
					Owns {$game.victoryConditions.ownCapitalShips} capital ships.
				</td>
				{#each players as player}
					<td>
						{#if (($universe.getPlayerScore(player.num)?.achievedVictoryConditions ?? 0) & VictoryCondition.OwnCapitalShips) > 0}
							<Icon src={CheckBadge} size="24" class="stroke-success" />
						{/if}
					</td>{/each}
			</tr>
		{/if}
		{#if ($game.victoryConditions.conditions & VictoryCondition.HighestScoreAfterYears) > 0}
			<tr>
				<td>
					Has the highest score after {$game.victoryConditions.highestScoreAfterYears} years.
				</td>
				{#each players as player}
					<td>
						{#if (($universe.getPlayerScore(player.num)?.achievedVictoryConditions ?? 0) & VictoryCondition.HighestScoreAfterYears) > 0}
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
