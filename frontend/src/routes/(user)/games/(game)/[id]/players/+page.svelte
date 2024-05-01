<script lang="ts">
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { startCase } from 'lodash-es';
	import PlayerScores from './PlayerScores.svelte';
	import PlayerScoresGraph, { type ValueType } from './PlayerScoresGraph.svelte';
	import PlayerVictoryConditions from './PlayerVictoryConditions.svelte';
	import PlayersStatus from './PlayersStatus.svelte';

	const { game, player, universe } = getGameContext();

	const graphTypes: ValueType[] = [
		'planets',
		'starbases',
		'unarmedShips',
		'escortShips',
		'capitalShips',
		'techLevels',
		'resources',
		'score'
	];

	let type: ValueType = 'score';
</script>

<div class="w-full mx-auto md:max-w-2xl">
	<ItemTitle>Players</ItemTitle>

	<PlayersStatus />
</div>
<div class="w-full">
	<ItemTitle>Score</ItemTitle>
	<div>
		<PlayerScores />
	</div>

	<div class="h-[400px] w-full sm:w-[42rem] mx-auto my-5 py-5 px-1">
		<PlayerScoresGraph {type} />
	</div>
	<div class="flex flex-row flex-wrap gap-1 justify-center">
		{#each graphTypes as graphType}
			<div class="form-control">
				<label
					class="label cursor-pointer btn w-[11rem]"
					class:bg-primary={type == graphType}
				>
					<span class="label-text text-center w-full">{startCase(graphType)}</span>
					<input
						type="radio"
						name="score-graph-value-type"
						class="hidden"
						value={graphType}
						bind:group={type}
					/>
				</label>
			</div>
		{/each}
	</div>
</div>

<div class="w-full">
	<ItemTitle>Victory Conditions</ItemTitle>
	<PlayerVictoryConditions />
</div>

<div class="mb-10" />
