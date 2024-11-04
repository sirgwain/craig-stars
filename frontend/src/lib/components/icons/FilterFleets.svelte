<script lang="ts">
	interface Props {
		primary?: string;
		[key: string]: any;
	}

	let { primary = 'fill-fleet-owned', ...rest }: Props = $props();
	export const secondary = 'fill-fleet-neutral';

	const size = 120;
	const primaryPositions = [
		[140, 100],
		[200, 230],
		[330, 330]
	];
	const secondaryPositions = [
		[100, 400],
		[400, 220],
		[280, 100]
	];
</script>

<svg
	viewBox="0 0 500 500"
	xmlns="http://www.w3.org/2000/svg"
	xmlns:xlink="http://www.w3.org/1999/xlink"
	{...rest}
>
	{#each primaryPositions as pos}
		<polygon
			points={`0,0 0,${size} ${size},${size}`}
			class={`${primary}`}
			transform={`translate(${pos[0]} ${pos[1]}) translate(${-size / 2} ${-size / 2})`}
		/>
	{/each}

	{#each secondaryPositions as pos}
		<polygon
			points={`0,0 0,${size} ${size},${size}`}
			transform={`translate(${pos[0]} ${pos[1]}) translate(${-size / 2} ${-size / 2})`}
		/>
	{/each}
</svg>
