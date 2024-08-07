<script lang="ts" context="module">
	import { type PlayerScore } from '$lib/types/Player';

	export type ValueType = keyof PlayerScore;
</script>

<script lang="ts">
	import AxisX from '$lib/components/graph/AxisX.html.svelte';
	import AxisY from '$lib/components/graph/AxisY.html.svelte';
	import MultiLine from '$lib/components/graph/MultiLine.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	import { scaleOrdinal } from 'd3-scale';
	import { Html, LayerCake, ScaledSvg } from 'layercake';
	import PlayerScoresGraphLabels from './PlayerScoresGraphLabels.svelte';

	const { game, universe } = getGameContext();

	export let type: ValueType = 'score';

	type DataLongTurnValueType = { player: string; turn: number; value: number };
	type DataLongType = { player: string; playerName: string; values: DataLongTurnValueType[] }[];

	const xKey = 'turn';
	const yKey = 'value';
	const zKey = 'player';

	const seriesNames: string[] = $universe.players.map<string>((p) => String(p.num));
	const seriesColors: string[] = $universe.players.map<string>((p) => p.color);

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

	let dataLong: DataLongType;

	$: {
		/* --------------------------------------------
		 * Create a "long" format that is a grouped series of data points
		 * Layer Cake uses this data structure and the key names
		 * set in xKey, yKey and zKey to map your data into each scale.
		 */
		dataLong = $universe.players.map((playerIntel, i) => {
			const name = playerIntel.racePluralName ?? playerIntel.name;
			const playerScores = $universe.scores[i];

			return {
				[zKey]: String(playerIntel.num),
				playerName: name,
				values: [...Array(turnsPassed).keys()].map((turn => ({
					[yKey]: playerScores && playerScores[turn] ? playerScores[turn][type] ?? 0 : 0,
					[xKey]: turn,
					[zKey]: String(playerIntel.num)
				})))
			}
		});
	}
</script>

<div class="border border-base-300 bg-base-100 w-full h-full">
	{#if highestValue === 0}
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
				<PlayerScoresGraphLabels />
				<!-- TODO: get this working so we can see values on our graphs -->
				<!-- <SharedTooltip dataset={dataQuadTree} /> -->
			</Html>
		</LayerCake>
	{/if}
</div>
