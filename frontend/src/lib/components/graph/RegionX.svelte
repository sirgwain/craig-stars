<!--
  @component
  Generates a full height SVG rect shading the region of the chart between two x values.
 -->
<script lang="ts">
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';

	const { xScale, yScale } = getContext<LayerCake>('LayerCake');

	type Props = {
		from: number;
		to: number;
		fill?: string;
	};

	let { from, to, fill = 'fill-warning opacity-25' }: Props = $props();

	let x = $derived(Math.min($xScale(from), $xScale(to)));
	let width = $derived(Math.abs($xScale(to) - $xScale(from)));
	let yRange = $derived($yScale.range());
	let y = $derived(Math.min(yRange[0], yRange[1]));
	let height = $derived(Math.abs(yRange[0] - yRange[1]));
</script>

<rect class="region-x {fill}" {x} {y} {width} {height} />
