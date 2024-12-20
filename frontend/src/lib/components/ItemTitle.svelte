<script lang="ts">
	import type { Snippet } from 'svelte';
	let { children }: { children?: Snippet } = $props();
</script>

<h3 class="text-2xl font-semibold mb-2 mt-4 border-primary border-b-2 p-2 text-center">
	{@render children?.()}
</h3>
