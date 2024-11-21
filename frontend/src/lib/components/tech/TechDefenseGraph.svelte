<script lang="ts">
	import { getDefenseCoverage, getSmartDefenseCoverage, type TechDefense } from '$lib/types/Tech';

	import { scaleOrdinal } from 'd3-scale';
	import { Html, LayerCake, ScaledSvg } from 'layercake';
	import AxisX from '../graph/AxisX.html.svelte';
	import AxisY from '../graph/AxisY.html.svelte';
	import GroupLabels from '../graph/GroupLabels.svelte';
	import MultiLine from '../graph/MultiLine.svelte';

	type Props = {
		defense: TechDefense;
	};

	let { defense }: Props = $props();

	const numTicks = 11;

	type CoverageType = { defenses: number; coverage: number };
	type DefenseType = 'Standard' | 'Smart';
	type DataType = { [k: string]: CoverageType[] };

	type DataLongCoverageType = { type: DefenseType; defenses: number; coverage: number };
	type DataLongType = { type: DefenseType; values: DataLongCoverageType[] };

	//type DataQuadTree = DataLongCoverageType[];

	let data: DataType = $derived.by(() => {
		const data: DataType = {
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
			}
		}

		return data;
	});

	/*let dataQuadTree: DataQuadTree = $derived.by(() => {
		const dataQuadTree: DataQuadTree = [];
		if (defense) {
			for (let i = 0; i <= 100; i += 100 / (numTicks - 1)) {
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

		return dataQuadTree;
	});*/

	const xKey = 'defenses';
	const yKey = 'coverage';
	const zKey = 'type';

	const seriesNames: DefenseType[] = ['Standard', 'Smart'];
	const seriesColors = ['stroke-primary', 'stroke-accent'];

	/* --------------------------------------------
	 * Create a "long" format that is a grouped series of data points
	 * Layer Cake uses this data structure and the key names
	 * set in xKey, yKey and zKey to map your data into each scale.
	 */
	let dataLong: DataLongType[] = $derived(
		seriesNames.map((key) => ({
			[zKey]: key,
			values: data[key].map((d) => {
				return {
					[yKey]: d.coverage,
					[xKey]: d.defenses,
					[zKey]: key
				};
			})
		}))
	);
	/* --------------------------------------------
	 * Make a flat array of the `values` of our nested series
	 * we can pluck the field set from `yKey` from each item
	 * in the array to measure the full extents
	 */
	// TODO - JD - SG Can you confirm that the flatMap on line 112 does the same as this function did?
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
		flatData={dataLong.flatMap((x) => x.values)}
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
			<!-- TODO: get this working so we can see values on our graphs -->
			<!-- <SharedTooltip dataset={dataQuadTree} /> -->
		</Html>
	</LayerCake>
</div>
