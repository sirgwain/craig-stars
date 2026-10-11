<!--
  @component
  Generates an HTML y-axis.
 -->
<script lang="ts">
	import { LayerCake } from 'layercake';
	import { getContext } from 'svelte';

	const { padding, xRange, yScale } = getContext<LayerCake>('LayerCake');

	type Props = {
		gridlines?: boolean;
		formatTick?: (d: unknown) => string;
		tickClass?: (d: unknown) => string;
		gridlineClass?: (d: unknown) => string;
		ticks?: number | Array<unknown> | undefined;
		xTick?: number;
		yTick?: number;
	};

	let {
		gridlines = true,
		formatTick = (d) => `${d}`,
		tickClass = () => 'text-base-content',
		gridlineClass = () => 'border-dashed',
		ticks = 4,
		xTick = -4,
		yTick = -1
	}: Props = $props();

	let isBandwidth = $derived(typeof $yScale.bandwidth === 'function');

	let tickVals = $derived(
		Array.isArray(ticks) ? ticks : isBandwidth ? $yScale.domain() : $yScale.ticks(ticks)
	);
</script>

<div class="axis y-axis" style="transform:translate(-{$padding.left}px, 0)">
	{#each tickVals as tick, i (tick)}
		<div
			class="tick tick-{i}"
			style="top:{$yScale(tick) + (isBandwidth ? $yScale.bandwidth() / 2 : 0)}%;left:{$xRange[0]}%;"
		>
			{#if gridlines !== false}
				<div
					class="border-t border-base-content {gridlineClass(tick)}"
					style="top:0;left:{isBandwidth ? $padding.left : 0}px;right:-{$padding.left +
						$padding.right}px;"
				></div>
			{/if}
			<div
				class={tickClass(tick)}
				style="
          top:{yTick}px;
          left:{isBandwidth ? $padding.left + xTick - 4 : 0}px;
          transform: translate({isBandwidth ? '-100%' : 0}, {isBandwidth
					? -50 - Math.floor($yScale.bandwidth() / -2)
					: '-100'}%);
        "
			>
				{formatTick(tick)}
			</div>
		</div>
	{/each}
</div>

<style>
	.axis,
	.tick,
	.tick-mark,
	.gridline,
	.text {
		position: absolute;
	}
	.axis {
		width: 100%;
		height: 100%;
	}
	.tick {
		font-size: 12px;
		width: 100%;
		font-weight: 100;
	}
</style>
