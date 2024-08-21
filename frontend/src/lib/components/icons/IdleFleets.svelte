<script lang="ts">
	export let color = '#0000FF';

	const fleetSize = 120;
	const letterSize = 200;
	const smallLetterSize = 140;
</script>

<svg
	viewBox="0 0 500 500"
	xmlns="http://www.w3.org/2000/svg"
	xmlns:xlink="http://www.w3.org/1999/xlink"
	{...$$restProps}
>
	<text
		x={letterSize / 2 + 30}
		y={letterSize / 2 + 30}
		class={`text-[${letterSize}px] text-base-content`}>z</text
	>
	<text
		x={(letterSize / 2) * 2 + 60}
		y={letterSize / 2 + 60}
		class={`text-[${letterSize}px] text-base-content`}>z</text
	>
	<text
		x={(letterSize / 2) * 3 + 90}
		y={(letterSize / 2) * 2 + 60}
		class={`text-[${smallLetterSize}px] text-base-content`}>z</text
	>
	<text x={0} y={letterSize / 2} class={`text-[${letterSize}px] text-base-content`}>z</text>
	<text
		x={letterSize / 2 + 30}
		y={(letterSize / 2) * 4 + 30}
		class={`text-[${letterSize}px] text-base-content`}>z</text
	>
	<text
		x={(letterSize / 2) * 2 + 60}
		y={(letterSize / 2) * 4 }
		class={`text-[${letterSize}px] text-base-content`}>z</text
	>
	<text x={0} y={(letterSize / 2)*5} class={`text-[${letterSize}px] text-base-content`}>z</text>
	<polygon
		points={`0,0 0,${fleetSize} ${fleetSize},${fleetSize}`}
		fill={color}
		transform={`translate(90 220) translate(${-fleetSize / 2} ${-fleetSize / 2})`}
	/>

</svg>
