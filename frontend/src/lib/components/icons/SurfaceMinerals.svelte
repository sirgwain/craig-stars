<script lang="ts">
	const ironiumSize = 400;
	const boraniumSize = 200;
	const germaniumSize = 300;
	const width = 140;
	const lineWidth = 20;
</script>

<svg
	viewBox="0 0 500 500"
	xmlns="http://www.w3.org/2000/svg"
	xmlns:xlink="http://www.w3.org/1999/xlink"
	{...$$restProps}
>
	<path class="stroke-white" stroke-width={lineWidth} fill="none" d={`M0,0L0,500 L500,500`} />

	<rect
		class="ironium-bar"
		x={lineWidth}
		y={500 - ironiumSize - lineWidth}
		{width}
		height={ironiumSize}
	/>
	<rect
		class="boranium-bar"
		x={lineWidth + width + 20}
		y={500 - boraniumSize - lineWidth}
		{width}
		height={boraniumSize}
	/>
	<rect
		class="germanium-bar"
		x={lineWidth + width*2 + 40}
		y={500 - germaniumSize - lineWidth}
		{width}
		height={germaniumSize}
	/>
	<text x=250 y=150 fill="white" class='text-[200px] text-white'>M</text>
</svg>
