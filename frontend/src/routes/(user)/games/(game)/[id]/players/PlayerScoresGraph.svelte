<script lang="ts" context="module">
	import { type PlayerScore } from '$lib/types/Player';

	export type ValueType = keyof PlayerScore;
</script>

<script lang="ts">
	import AxisX from '$lib/components/graph/AxisX.html.svelte';
	import AxisY from '$lib/components/graph/AxisY.html.svelte';
	import GroupLabels from '$lib/components/graph/GroupLabels.svelte';
	import MultiLine from '$lib/components/graph/MultiLine.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	import { scaleOrdinal } from 'd3-scale';
	import { Html, LayerCake, ScaledSvg } from 'layercake';

	const { game, universe, player } = getGameContext();

	export let type: ValueType = 'score';

	type ScoreTurnValueType = { turn: number; value: number };
	type DataType = { [k: string]: ScoreTurnValueType[] };

	type DataLongTurnValueType = { player: string; turn: number; value: number };
	type DataLongType = { player: string; values: DataLongTurnValueType[] }[];

	let data: DataType = {};

	const xKey = 'turn';
	const yKey = 'value';
	const zKey = 'player';

	const seriesNames: string[] = $universe.players.map<string>((p) => p.racePluralName ?? p.name);
	const seriesColors: string[] = $universe.players.map<string>((p) => p.color);

	let dataLong: DataLongType;
	/* --------------------------------------------
	 * Make a flat array of the `values` of our nested series
	 * we can pluck the field set from `yKey` from each item
	 * in the array to measure the full extents
	 */
	const flatten = (data: any) =>
		data.reduce((memo: any, group: []) => {
			return memo.concat(group.values);
		}, []);

	// get the number of turns passed, i.e. 2 for 2402
	$: turnsPassed = $game.year - $game.rules.startingYear;

	// get the highest value from the scores
	$: highestValue = Math.max(
		...$universe.scores
			.filter((score) => score && score.length > 0)
			.flat()
			.map((score) => score[type] ?? 0)
	);

	$: {
		// init an empty array for each player name
		data = {};
		seriesNames.forEach((p) => (data[p] = []));

		const scores = $universe.scores;

		for (let turn = 0; turn < turnsPassed; turn += 1 /*turnsPassed / (numTicks - 1)*/) {
			$universe.players.forEach((playerIntel, i) => {
				const name = playerIntel.racePluralName ?? playerIntel.name;
				const playerScores = scores[i];
				data[name].push({
					turn: turn,
					value: playerScores && playerScores[turn] ? playerScores[turn][type] ?? 0 : 0
				});
			});
		}

		/* --------------------------------------------
		 * Create a "long" format that is a grouped series of data points
		 * Layer Cake uses this data structure and the key names
		 * set in xKey, yKey and zKey to map your data into each scale.
		 */
		dataLong = seriesNames.map((key) => {
			return {
				[zKey]: key,
				values: data[key].map((d) => {
					return {
						[yKey]: d.value,
						[xKey]: d.turn,
						[zKey]: key
					};
				})
			};
		});
	}
</script>

<div class="border border-base-300 bg-base-100 w-full h-full">
	{#if highestValue == 0}
		<div class="flex flex-row justify-center h-full">
			<div class="my-auto">No Data</div>
		</div>
	{:else}
		<LayerCake
			x={xKey}
			y={yKey}
			z={zKey}
			yDomain={[0, Math.ceil(highestValue + highestValue * 0.2)]}
			xDomain={[0, turnsPassed]}
			xRange={[0, 100]}
			yRange={[100, 0]}
			zScale={scaleOrdinal()}
			zRange={seriesColors}
			zDomain={seriesNames}
			flatData={flatten(dataLong)}
			data={dataLong}
		>
			<Html>
				<AxisX />
				<AxisY ticks={6} formatTick={(d) => `${d}`} />
			</Html>
			<ScaledSvg>
				<MultiLine colorIsCode={true} />
			</ScaledSvg>

			<Html>
				<GroupLabels />
				<!-- TODO: get this working so we can see values on our graphs -->
				<!-- <SharedTooltip dataset={dataQuadTree} /> -->
			</Html>
		</LayerCake>
	{/if}
</div>
