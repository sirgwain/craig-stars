<script lang="ts">
	import { run } from 'svelte/legacy';

	import type { TechEngine } from '$lib/types/Tech';

	import { scaleSqrt } from 'd3-scale';
	import { Html, LayerCake, ScaledSvg } from 'layercake';
	import Area from '../graph/Area.svelte';
	import AxisX from '../graph/AxisX.html.svelte';
	import AxisY from '../graph/AxisY.html.svelte';
	import Line from '../graph/Line.svelte';

	interface Props {
		engine: TechEngine;
	}

	let { engine }: Props = $props();

	type DataType = [number, number][];

	let data: DataType = $state([]);

	const xGetter = (d: DataType) => d[0];
	const yGetter = (d: DataType) => d[1];

	run(() => {
		data = engine?.fuelUsage ? engine.fuelUsage.map((usage, index) => [index, usage]) : [];
	});
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
