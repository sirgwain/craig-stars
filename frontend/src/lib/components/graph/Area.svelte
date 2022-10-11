<!--
  @component
  Generates an SVG area shape using the `area` function from [d3-shape](https://github.com/d3/d3-shape).
 -->
<script lang="ts">
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';

	const { data, xGet, yGet, xScale, yScale, extents, width, height } =
		getContext<LayerCake>('LayerCake');

	export let fill = 'fill-secondary-content opacity-50';

	$: path =
		'M' +
		$data
			.map((d: any) => {
				return $xGet(d) + ',' + $yGet(d);
			})
			.join('L');

	let area: string;

	$: {
		const yRange = $yScale.range();
		area =
			path +
			('L' +
				$xScale($extents.x ? $extents.x[1] : 0) +
				',' +
				yRange[0] +
				'L' +
				$xScale($extents.x ? $extents.x[0] : 0) +
				',' +
				yRange[0] +
				'Z');
	}
</script>

<path class="path-area {fill}" d={area} />
