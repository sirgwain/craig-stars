<script lang="ts">
	import { getDefenseCoverage, getSmartDefenseCoverage, type TechDefense } from '$lib/types/Tech';

	import { scaleOrdinal } from 'd3-scale';
	import { Html, LayerCake, ScaledSvg } from 'layercake';
	import AxisX from '../graph/AxisX.html.svelte';
	import AxisY from '../graph/AxisY.html.svelte';
	import GroupLabels from '../graph/GroupLabels.svelte';
	import MultiLine from '../graph/MultiLine.svelte';

	export let defense: TechDefense;

	const numTicks = 11;

	type CoverageType = { defenses: number; coverage: number };
	type DefenseType = 'Standard' | 'Smart';
	type DataType = { [k: string]: CoverageType[] };

	type DataLongCoverageType = { type: DefenseType; defenses: number; coverage: number };
	type DataLongType = { type: DefenseType; values: DataLongCoverageType[] }[];

	type DataQuadTree = DataLongCoverageType[];

	let data: DataType = {};
	let dataQuadTree: DataQuadTree = [];

	const xKey = 'defenses';
	const yKey = 'coverage';
	const zKey = 'type';

	const seriesNames: DefenseType[] = ['Standard', 'Smart'];
	const seriesColors = ['stroke-primary', 'stroke-accent'];

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

	$: {
		data = {
			Standard: [],
			Smart: []
		};
		if (defense) {
			for (let i = 0; i <= 100; i += 100 / (numTicks - 1)) {
				data['Standard'].push({
					defenses: i,
					coverage: getDefenseCoverage(defense, i) * 100
				});
				data['Smart'].push({
					defenses: i,
					coverage: getSmartDefenseCoverage(defense, i) * 100
				});

				dataQuadTree.push({
					type: 'Standard',
					defenses: i,
					coverage: getDefenseCoverage(defense, i) * 100
				});
				dataQuadTree.push({
					type: 'Smart',
					defenses: i,
					coverage: getSmartDefenseCoverage(defense, i) * 100
				});
			}
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
						[yKey]: d.coverage,
						[xKey]: d.defenses,
						[zKey]: key
					};
				})
			};
		});
	}
</script>

<div class="border border-base-300 bg-base-100 w-full h-full mt-5 pb-7">
	<LayerCake
		x={xKey}
		y={yKey}
		z={zKey}
		yDomain={[0, 100]}
		xDomain={[0, 100]}
		xRange={[0, 100]}
		yRange={[100, 0]}
		zScale={scaleOrdinal()}
		zRange={seriesColors}
		flatData={flatten(dataLong)}
		data={dataLong}
		ssr={true}
	>
		<Html>
			<AxisX />
			<AxisY ticks={5} formatTick={(d) => `${d}%`} />
		</Html>
		<ScaledSvg>
			<MultiLine />
		</ScaledSvg>

		<Html>
			<GroupLabels />
			<!-- <SharedTooltip dataset={dataQuadTree} /> -->
		</Html>
	</LayerCake>
</div>
