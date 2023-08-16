<script lang="ts">
	export let color = '#0000FF';

	$: size = 256;
</script>

<svg
	viewBox="0 0 500 500"
	xmlns="http://www.w3.org/2000/svg"
	xmlns:xlink="http://www.w3.org/1999/xlink"
	{...$$restProps}
>
	<polygon
		points={`0,0 0,${size} ${size},${size}`}
		fill={color}
		transform={`translate(350 400) translate(${
			-size / 2
		} ${-size / 2})`}
	/>
	<text x="200" y="150" class="text-[200px] text-base-content">32</text>
</svg>
