<!--
  @component
  Generates an HTML x-axis, useful for server-side rendered charts.  This component is also configured to detect if your x-scale is an ordinal scale. If so, it will place the markers in the middle of the bandwidth.
 -->
<script lang="ts">
	import { getContext } from 'svelte';

	const { xScale } = getContext('LayerCake');

	interface Props {
		gridlines?: Boolean;
		tickMarks?: Boolean;
		baseline?: Boolean;
		snapTicks?: Boolean;
		/** @type {Function} [formatTick=d => d] - A function that passes the current tick value and expects a nicely formatted value in return. */
		formatTick?: (d: any) => string;
		/** @type {Number|Array|Function} [ticks] - If this is a number, it passes that along to the [d3Scale.ticks](https://github.com/d3/d3-scale) function. If this is an array, hardcodes the ticks to those values. If it's a function, passes along the default tick values and expects an array of tick values in return. If nothing, it uses the default ticks supplied by the D3 function. */
		ticks?: Number | Array<any> | Function | undefined;
		yTick?: Number;
	}

	let {
		gridlines = true,
		tickMarks = false,
		baseline = false,
		snapTicks = false,
		formatTick = (d) => d,
		ticks = undefined,
		yTick = 7
	}: Props = $props();

	let isBandwidth = $derived(typeof $xScale.bandwidth === 'function');

	let tickVals = $derived(
		Array.isArray(ticks)
			? ticks
			: isBandwidth
				? $xScale.domain()
				: typeof ticks === 'function'
					? ticks($xScale.ticks())
					: $xScale.ticks(ticks)
	);
</script>

<div class="axis x-axis" class:snapTicks>
	{#each tickVals as tick, i (tick)}
		{#if gridlines !== false}
			<div
				class="border-l border-l-base-content"
				style="left:{$xScale(tick)}%;top: 0px;bottom: 0;"
			></div>
		{/if}
		{#if tickMarks === true}
			<div
				class="border-l-base-content"
				style="left:{$xScale(tick) +
					(isBandwidth ? $xScale.bandwidth() / 2 : 0)}%;height:6px;bottom: -6px;"
			></div>
		{/if}
		<div
			class="tick tick-{i}"
			style="left:{$xScale(tick) + (isBandwidth ? $xScale.bandwidth() / 2 : 0)}%;top:100%;"
		>
			<div class="text text-base-content" style="top:{yTick}px;">{formatTick(tick)}</div>
		</div>
	{/each}
	{#if baseline === true}
		<div class="border-t border-t-base-content" style="top: 100%;width: 100%;"></div>
	{/if}
</div>

<style>
	.axis,
	.tick,
	.tick-mark,
	.gridline,
	.baseline {
		position: absolute;
	}
	.axis {
		width: 100%;
		height: 100%;
	}
	.tick {
		font-size: 0.725em;
		font-weight: 200;
	}

	.tick .text {
		position: relative;
		white-space: nowrap;
		transform: translateX(-50%);
	}
	/* This looks a little better at 40 percent than 50 */
	.axis.snapTicks .tick:last-child {
		transform: translateX(-40%);
	}
	.axis.snapTicks .tick.tick-0 {
		transform: translateX(40%);
	}
</style>
