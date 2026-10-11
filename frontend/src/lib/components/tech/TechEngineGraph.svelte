<script lang="ts">
	import { scaleSqrt } from 'd3-scale';
	import { Html, LayerCake, ScaledSvg } from 'layercake';
	import AxisX from '../graph/AxisX.html.svelte';
	import AxisY from '../graph/AxisY.html.svelte';
	import Line from '../graph/Line.svelte';
	import RegionX from '../graph/RegionX.svelte';
	import type { EngineJson } from '#lib/types/cs-proto.js';

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

	// warp speeds above the engine's max safe speed damage the fleet
	let maxSpeed = $derived(data.length - 1);
	let maxSafeSpeed = $derived(engine.maxSafeSpeed ?? 0);
	let hasUnsafeSpeeds = $derived(maxSafeSpeed > 0 && maxSafeSpeed < maxSpeed);

	const xGetter = (d: DataType) => d[0];
	const yGetter = (d: DataType) => d[1];
</script>

<div class="border border-base-300 bg-base-100 size-full mt-5 pb-7">
	<LayerCake
		percentRange={true}
		x={xGetter}
		y={yGetter}
		yDomain={[1, 1200]}
		yScale={scaleSqrt()}
		{data}
	>
		<Html>
			<AxisX
				tickClass={(d) =>
					d === engine.idealSpeed ? 'text-primary font-bold' : 'text-base-content'}
			/>
			<AxisY
				ticks={[0, 25, 100, 200, 400, 800]}
				formatTick={(d) => `${d as number}%`}
				tickClass={(d) => ((d as number) > 100 ? 'text-error' : 'text-base-content')}
				gridlineClass={(d) => (d === 100 ? 'border-solid' : 'border-dashed')}
			/>
		</Html>
		<ScaledSvg>
			{#if hasUnsafeSpeeds}
				<RegionX from={maxSafeSpeed} to={maxSpeed} />
			{/if}
			<Line />
		</ScaledSvg>
	</LayerCake>
</div>
