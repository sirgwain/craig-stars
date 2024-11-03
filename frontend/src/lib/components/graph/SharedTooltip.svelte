<!--
  @component
  Generates a tooltip that works on multiseries datasets, like multiline charts. It creates a tooltip showing the name of the series and the current value. It finds the nearest data point using the [QuadTree.html.svelte](https://layercake.graphics/components/QuadTree.html.svelte) component.
 -->
<script lang="ts">
	import { format } from 'd3-format';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import QuadTree from './QuadTree.svelte';

	const { data, xGet, yGet, zGet, xScale, yScale, width, height, config } =
		getContext<LayerCake>('LayerCake');

	const commas = format(',');
	const titleCase = (d: any) => d.replace(/^\w/, (w: string) => w.toUpperCase());

	

	

	

	

	
	interface Props {
		formatTitle?: Function;
		formatValue?: Function;
		formatKey?: Function;
		offset?: Number;
		/** @type {Array} [dataset] - The dataset to work off of—defaults to $data if left unset. You can pass something custom in here in case you don't want to use the main data or it's in a strange format. */
		dataset?: any;
	}

	let {
		formatTitle = (d: any) => d,
		formatValue = (d: any) => (isNaN(+d) ? d : commas(d)),
		formatKey = (d: any) => titleCase(d),
		offset = -20,
		dataset = undefined
	}: Props = $props();

	const w = 150;
	const w2 = w / 2;

	/* --------------------------------------------
	 * Sort the keys by the highest value
	 */
	function sortResult(result: any) {
		if (Object.keys(result).length === 0) return [];
		const rows = Object.keys(result)
			.filter((d) => d !== $config.x)
			.map((key) => {
				return {
					key,
					value: result[key]
				};
			})
			.sort((a, b) => b.value - a.value);
		return rows;
	}
</script>

<QuadTree dataset={dataset || $data} y="x"     >
	{#snippet children({ x, y, visible, found, e })}
		{@const foundSorted = sortResult(found)}
		{#if visible === true}
			<div style="left:{x}px;" class="line"></div>
			<div
				class="tooltip"
				style="
	          width:{w}px;
	          display: {visible ? 'block' : 'none'};
	          top:{$yScale(foundSorted[0].value) + offset}px;
	          left:{Math.min(Math.max(w2, x), $width - w2)}px;"
			>
				<div class="title">{formatTitle(found[$config.x])}</div>
				{#each foundSorted as row}
					<div class="row">
						<span class="key">{formatKey(row.key)}:</span>
						{formatValue(row.value)}
					</div>
				{/each}
			</div>
		{/if}
	{/snippet}
</QuadTree>

<style>
	.tooltip {
		position: absolute;
		font-size: 13px;
		pointer-events: none;
		border: 1px solid #ccc;
		background: rgba(255, 255, 255, 0.85);
		transform: translate(-50%, -100%);
		padding: 5px;
		z-index: 15;
		pointer-events: none;
	}
	.line {
		position: absolute;
		top: 0;
		bottom: 0;
		width: 1px;
		border-left: 1px dotted #666;
		pointer-events: none;
	}
	.tooltip,
	.line {
		transition: left 250ms ease-out, top 250ms ease-out;
	}
	.title {
		font-weight: bold;
	}
	.key {
		color: #999;
	}
</style>
