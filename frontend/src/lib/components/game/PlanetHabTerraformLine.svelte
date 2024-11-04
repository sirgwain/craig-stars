<script lang="ts">
	interface Props {
		style?: string | undefined;
		x1?: number;
		x2?: number;
		[key: string]: any;
	}

	let { style = undefined, x1 = 50, x2 = 50, ...rest }: Props = $props();
</script>

<div class="absolute h-full w-full">
	<svg
		width="100%"
		height="100%"
		viewBox="0 0 100 100"
		preserveAspectRatio="none"
		{style}
		{...rest}
	>
		<line {x1} y1="50" {x2} y2="50" vector-effect="non-scaling-stroke" stroke-width="1" />
	</svg>
</div>
