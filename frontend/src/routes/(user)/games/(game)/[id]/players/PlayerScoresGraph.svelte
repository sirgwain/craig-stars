<script module lang="ts">
	export type ValueType = keyof PlayerScoreJson;
</script>

<script lang="ts">
	import AxisX from '#lib/components/graph/AxisX.html.svelte';
	import AxisY from '#lib/components/graph/AxisY.html.svelte';
	import MultiLine from '#lib/components/graph/MultiLine.svelte';
	import { getGameContext } from '#lib/services/GameContext.js';

	import type { PlayerScoreJson } from '#lib/types/cs-proto.js';
	import { scaleOrdinal } from 'd3-scale';
	import { Html, LayerCake, ScaledSvg } from 'layercake';

	const { game, player, universe } = getGameContext();

	type Props = {
		type?: ValueType;
	};

	let { type = 'score' }: Props = $props();

	type DataLongTurnValueType = { player: string; turn: number; value: number };
	type DataLongTypeItem = {
		player: string;
		playerName: string;
		playerNum: number;
		values: DataLongTurnValueType[];
	};

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
	let turnsPassed = $derived($game.year - ($game.rules.universeGenerationRules?.startingYear ?? 0));

	// other players' scores are hidden for the first few years of the game
	let hiddenYears = $derived($game.rules.showPublicScoresAfterYears);

	// get the highest value from the scores
	let highestValue = $derived(
		Math.max(
			...$universe.scoreIntels
				.map((score) => score.scoreHistory)
				.filter((scoreHistory) => scoreHistory.length > 0)
				.map((scoreHistory) => scoreHistory as PlayerScoreJson[]) // make the types happy
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
			const name = playerIntel.racePluralName || playerIntel.name;
			const num = playerIntel.num;
			const playerScores = $universe.scoreIntels[i].scoreHistory;

			return {
				[zKey]: String(playerIntel.num),
				playerName: name,
				playerNum: num,
				values: [...Array(turnsPassed).keys()]
					.filter((turn) => num === $player.num || turn >= hiddenYears)
					.map((turn) => ({
						[yKey]: playerScores[turn] ? playerScores[turn][type] || 0 : 0,
						[xKey]: turn,
						[zKey]: String(playerIntel.num)
					}))
			};
		})
	);
</script>

<div class="flex flex-col size-full">
	<div class="border border-base-300 bg-base-100 grow min-h-0">
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
					<!-- scores count planets, ships and so on, so there's no half of one -->
					<AxisY ticks={6} wholeNumbers={true} formatTick={(d) => `${d}`} />
				</Html>
				<ScaledSvg>
					<MultiLine zIsColorCode={true} />
				</ScaledSvg>

				<Html>
					<!-- TODO: get this working so we can see values on our graphs 
			 https://layercake.graphics/components/SharedTooltip.html.svelte
			 -->
					<!-- <SharedTooltip dataset={dataQuadTree} /> -->
				</Html>
			</LayerCake>
		{/if}
	</div>
	<!-- names on the lines run into each other, so each player's color is listed under the graph -->
	<ul
		class="flex flex-row flex-wrap justify-center gap-x-4 gap-y-1 mt-7 text-sm"
		data-type="legend"
	>
		{#each dataLong as series (series.playerNum)}
			<li class="flex flex-row items-center min-w-0 max-w-full">
				<span
					class="inline-block w-4 h-1 mr-1 shrink-0"
					style={`background-color: ${$universe.getPlayerColor(series.playerNum)}`}
				></span>
				<span class="truncate" title={series.playerName}>{series.playerName}</span>
			</li>
		{/each}
	</ul>
</div>
