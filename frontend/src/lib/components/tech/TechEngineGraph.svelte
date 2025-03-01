<script lang="ts">
	import { scaleSqrt } from 'd3-scale';
	import { Html, LayerCake, ScaledSvg } from 'layercake';
	import Area from '../graph/Area.svelte';
	import AxisX from '../graph/AxisX.html.svelte';
	import AxisY from '../graph/AxisY.html.svelte';
	import Line from '../graph/Line.svelte';
	import type { TechEngine } from '$lib/types/cs';

	type Props = {
		engine: TechEngine;
	};

	let { engine }: Props = $props();

	type DataType = [number, number][];

	let data: DataType = $derived(
		engine?.fuelUsage ? engine.fuelUsage.map((usage, index) => [index, usage]) : []
	);

	const xGetter = (d: DataType) => d[0];
	const yGetter = (d: DataType) => d[1];
</script>

<div class="border border-base-300 bg-base-100 w-full h-full mt-5 pb-7">
	<LayerCake
		percentRange={true}
		x={xGetter}
		y={yGetter}
		yDomain={[1, 1200]}
		yScale={scaleSqrt()}
		{data}
	>
		<Html>
			<AxisX />
			<AxisY ticks={[0, 25, 100, 200, 400, 800]} formatTick={(d) => `${d}%`} />
		</Html>
		<ScaledSvg>
			<Line />
			<Area />
		</ScaledSvg>
	</LayerCake>
</div>
