<script lang="ts">
	type Props = {
		style?: string | undefined;
		[key: string]: unknown;
	};

	let { style = undefined, ...rest }: Props = $props();
</script>

<svg viewBox="0 0 12 6" fill="none" stroke="white" stroke-width="1.5" {style} {...rest}>
	<g transform="translate(0,0)">
		<!-- start at 5,15 draw line from there to 5,-5, then 5,5 -->
		<path stroke-linecap="round" stroke-linejoin="round" d="M1 6l5-5 5 5" />
	</g>
</svg>
