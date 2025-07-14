<script lang="ts" module>
	import { type PlayerScore } from '$lib/types/cs';
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

	type Props = {
		type?: ValueType;
	};

	let { type = 'score' }: Props = $props();

	type DataLongTurnValueType = { player: string; turn: number; value: number };
	type DataLongTypeItem = { player: string; playerName: string; values: DataLongTurnValueType[] };

	const xKey = 'turn';
	const yKey = 'value';
	const zKey = 'player';

	const seriesNames: string[] = $universe.playerIntels.map<string>((p) => String(p.num));
	const seriesColors: string[] = $universe.playerIntels.map<string>((p) => p.color);

	/* --------------------------------------------
	 * Make a flat array of the `values` of our nested series
	 * we can pluck the field set from `yKey` from each item
	 * in the array to measure the full extents
	 */
	function flatten(data: DataLongTypeItem[]) {
		return data.reduce((memo: DataLongTurnValueType[], group: DataLongTypeItem) => {
			return memo.concat(group.values);
		}, []);
	}

	// get the number of turns passed, i.e. 2 for 2402
	let turnsPassed = $derived($game.year - $game.rules.startingYear);

	// get the highest value from the scores
	let highestValue = $derived(
		Math.max(
			...$universe.scoreIntels
				.map((score) => score.scoreHistory)
				.filter((scoreHistory) => scoreHistory && scoreHistory.length > 0)
				.map((scoreHistory) => scoreHistory as PlayerScore[]) // make the types happy
				.flat()
				.map((scoreHistory) => scoreHistory[type] ?? 0)
		)
	);

	/* --------------------------------------------
	 * Create a "long" format that is a grouped series of data points
	 * Layer Cake uses this data structure and the key names
	 * set in xKey, yKey and zKey to map your data into each scale.
	 */
	let dataLong: DataLongTypeItem[] = $derived(
		$universe.playerIntels.map((playerIntel, i) => {
			const name = playerIntel.racePluralName ?? playerIntel.name;
			const playerScores = $universe.scoreIntels[i].scoreHistory;

			return {
				[zKey]: String(playerIntel.num),
				playerName: name,
				values: [...Array(turnsPassed).keys()].map((turn) => ({
					[yKey]: playerScores && playerScores[turn] ? (playerScores[turn][type] ?? 0) : 0,
					[xKey]: turn,
					[zKey]: String(playerIntel.num)
				}))
			};
		})
	);
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
				<MultiLine zIsColorCode={true} />
			</ScaledSvg>

			<Html>
				<PlayerScoresGraphLabels />
				<!-- TODO: get this working so we can see values on our graphs 
			 https://layercake.graphics/components/SharedTooltip.html.svelte
			 -->
				<!-- <SharedTooltip dataset={dataQuadTree} /> -->
			</Html>
		</LayerCake>
	{/if}
</div>
