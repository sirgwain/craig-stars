<script lang="ts">
	export let style: string | undefined = undefined;
</script>

<svg viewBox="0 0 12 6" fill="none" stroke="white" stroke-width="1.5" {style} {...$$restProps}>
	<g transform="translate(0,0)">
		<!-- start at 5,15 draw line from there to 5,-5, then 5,5 -->
		<path stroke-linecap="round" stroke-linejoin="round" d="M1 6l5-5 5 5" />
	</g>
</svg>
