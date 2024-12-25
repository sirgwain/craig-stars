<script lang="ts">
	import type { MineField } from '$lib/types/MineField';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';

	const { xGet, yGet, xScale, yScale } = getContext<LayerCake>('LayerCake');

	type Props = {
		mineField: MineField;
		color?: string;
		selected?: boolean;
	};

	let { mineField, color = '#0900FF', selected = false }: Props = $props();
</script>

<circle
	cx={$xGet(mineField)}
	cy={$yGet(mineField)}
	r={$xScale(mineField.spec.radius)}
	mask="url(#mask-minefield)"
	fill={color}
	class:selected
/>

{#if selected}
	<rect
		width={$xScale(2)}
		height={$yScale(2)}
		rx={0.5}
		x={$xGet(mineField) - $xScale(1)}
		y={$yGet(mineField) - $yScale(1)}
		fill={color}
	/>
{/if}

<style>
	.selected {
		filter: brightness(0.6);
	}
</style>
