<script lang="ts">
	import type { SVGAttributes } from 'svelte/elements';

	let props: SVGAttributes<SVGElement> = $props();
</script>

<svg
	viewBox="0 0 512 512"
	xmlns="http://www.w3.org/2000/svg"
	xmlns:xlink="http://www.w3.org/1999/xlink"
	{...props}
>
	<polygon stroke-width={50} points={`256,0 512,256 256,512 0,256`} />
</svg>
