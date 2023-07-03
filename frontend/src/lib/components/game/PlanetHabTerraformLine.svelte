<script lang="ts">
	export let style: string | undefined = undefined;
	export let x1 = 50;
	export let x2 = 50;
</script>


<div class="absolute h-full w-full">
	<svg
		width="100%"
		height="100%"
		viewBox="0 0 100 100"
		preserveAspectRatio="none"
		{style}
		{...$$restProps}
	>
		<line
			x1={x1}
			y1="50"
			x2={x2}
			y2="50"
			vector-effect="non-scaling-stroke"
			stroke-width="1"			
		/>
	</svg>
</div>
