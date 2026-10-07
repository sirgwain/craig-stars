<script lang="ts">
	import type { SVGAttributes } from 'svelte/elements';

	type Props = SVGAttributes<SVGElement>;

	let { ...rest }: Props = $props();
</script>

<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" {...rest}>
	<path fill="#00e600" d="M15 5h4v4h-4z" />
	<path fill="#ffff00" d="M19 2h4v4h-4z" />
	<text
		x="1"
		y="22"
		fill="currentColor"
		font-family="Arial, sans-serif"
		font-size="15"
		textLength="22"
		lengthAdjust="spacingAndGlyphs">Abc</text
	>
</svg>
