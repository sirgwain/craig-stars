<!--
  @component
  Generates an SVG area shape using the `area` function from [d3-shape](https://github.com/d3/d3-shape).
 -->
<script lang="ts">
	import { getContext } from 'svelte';

	const { data, xGet, yGet } = getContext('LayerCake');

	export let stroke = 'stroke-primary';

	$: path =
		'M' +
		$data
			.map((d: any) => {
				return $xGet(d) + ',' + $yGet(d);
			})
			.join('L');
</script>

<!-- <path class='path-line' d='{path}' {stroke}></path> -->
<path class="path-line {stroke}" d={path} />

<style>
	.path-line {
		fill: none;
		stroke-linejoin: round;
		stroke-linecap: round;
		stroke-width: 3;
	}
</style>
