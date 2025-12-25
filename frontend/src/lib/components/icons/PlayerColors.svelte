<script lang="ts">
	import { getColor } from '$lib/components/game/newgame/playerColors';
	import type { SVGAttributes } from 'svelte/elements';

	type Props = SVGAttributes<SVGElement>;
	let { ...rest }: Props = $props();

	const size = 120;
	const positions = [
		[140, 100],
		[200, 230],
		[330, 330],
		[100, 400],
		[400, 220],
		[280, 100]
	];

	// Use the first 6 player colors
	const colors = [0, 1, 2, 3, 4, 5].map((i) => getColor(i));
</script>

<svg
	viewBox="0 0 500 500"
	xmlns="http://www.w3.org/2000/svg"
	xmlns:xlink="http://www.w3.org/1999/xlink"
	{...rest}
>
	{#each positions as pos, index (pos)}
		<polygon
			points={`0,0 0,${size} ${size},${size}`}
			fill={colors[index]}
			transform={`translate(${pos[0]} ${pos[1]}) translate(${-size / 2} ${-size / 2})`}
		/>
	{/each}
</svg>
