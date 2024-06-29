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
	<path d="M128,256L384,256 M256,128L256,384" stroke-width="32" />
</svg>
