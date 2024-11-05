<script>
	/**
	 * @typedef {Object} Props
	 * @property {import('svelte').Snippet} [children]
	 */

	/** @type {Props} */
	let { children } = $props();
</script>

<h3 class="text-2xl font-semibold mb-2 mt-4 border-primary border-b-2 p-2 text-center">
	{@render children?.()}
</h3>
