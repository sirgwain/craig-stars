<script lang="ts">
	import type { TechEngine } from '$lib/types/Tech';

	import { Html, LayerCake, ScaledSvg } from 'layercake';
	import Area from '../graph/Area.svelte';
	import AxisX from '../graph/AxisX.html.svelte';
	import AxisY from '../graph/AxisY.html.svelte';
	import Line from '../graph/Line.svelte';

	export let engine: TechEngine;

	type DataType = [number, number][];

	let data: DataType = [];

	const xGetter = (d: DataType) => d[0];
	const yGetter = (d: DataType) => d[1];

	$: data = engine?.fuelUsage ? engine.fuelUsage.map((usage, index) => [index, usage]) : [];
</script>

<div class="border border-base-300 bg-base-100 w-full h-full mt-5 pb-7">
	<LayerCake
		ssr={true}
		percentRange={true}		
		x={xGetter}
		y={yGetter}
		yDomain={[0, 1200]}
		{data}
	>
		<Html>
			<AxisX />
			<AxisY ticks={4} formatTick={(d) => `${d}%`} />
		</Html>
		<ScaledSvg>
			<Line />
			<Area />
		</ScaledSvg>
	</LayerCake>
</div>
