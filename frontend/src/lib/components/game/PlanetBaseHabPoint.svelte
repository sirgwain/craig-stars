<script lang="ts">
	import type { SVGAttributes } from 'svelte/elements';

	type Props = {
		style?: string | undefined;
	} & SVGAttributes<SVGElement>;

	let { ...rest }: Props = $props();
</script>

<svg
	viewBox="0 0 512 512"
	xmlns="http://www.w3.org/2000/svg"
	xmlns:xlink="http://www.w3.org/1999/xlink"
	{...rest}
>
	<path d="M128,256L384,256 M256,128L256,384" stroke-width="32" />
</svg>
