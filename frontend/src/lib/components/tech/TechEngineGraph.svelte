<script lang="ts">
	import { scaleSqrt } from 'd3-scale';
	import { Html, LayerCake, ScaledSvg } from 'layercake';
	import Area from '../graph/Area.svelte';
	import AxisX from '../graph/AxisX.html.svelte';
	import AxisY from '../graph/AxisY.html.svelte';
	import Line from '../graph/Line.svelte';
	import type { EngineJson } from '$lib/types/cs-proto';

	type Props = {
		engine: EngineJson;
	};

	let { engine }: Props = $props();

	type DataType = [number, number][];

	let data: DataType = $derived(
		engine.fuelUsage
			? engine.fuelUsage.map((usage: number, index: number): [number, number] => [index, usage])
			: []
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
			<AxisY ticks={[0, 25, 100, 200, 400, 800]} formatTick={(d) => `${d as number}%`} />
		</Html>
		<ScaledSvg>
			<Line />
			<Area />
		</ScaledSvg>
	</LayerCake>
</div>
