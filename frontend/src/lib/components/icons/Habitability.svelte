<script>
	/** @type {{ [key: string]: any }} */
	let { ...rest } = $props();
</script>

<svg
	viewBox="0 0 500 500"
	xmlns="http://www.w3.org/2000/svg"
	xmlns:xlink="http://www.w3.org/1999/xlink"
	{...rest}
>
	<text
		x="50%"
		y="50%"
		fill="white"
		dominant-baseline="middle"
		text-anchor="middle"
		class="text-[300px] text-white">%</text
	>
</svg>
