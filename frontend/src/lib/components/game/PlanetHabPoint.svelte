<script lang="ts">
	export let style: string | undefined = undefined;
</script>

<svg
	viewBox="0 0 512 512"
	xmlns="http://www.w3.org/2000/svg"
	xmlns:xlink="http://www.w3.org/1999/xlink"
	{style}
	{...$$restProps}
>
	<polygon points={`256,0 512,256 256,512 0,256`} />
</svg>
