<script lang="ts">
	import type { SVGAttributes } from 'svelte/elements';

	type Props = {
		style?: string | undefined;
	} & SVGAttributes<SVGElement>;

	let { style = undefined, ...rest }: Props = $props();
</script>

<svg
	viewBox="0 0 512 512"
	xmlns="http://www.w3.org/2000/svg"
	xmlns:xlink="http://www.w3.org/1999/xlink"
	{style}
	{...rest}
>
	<polygon stroke-width={50} points={`256,0 512,256 256,512 0,256`} />
</svg>
