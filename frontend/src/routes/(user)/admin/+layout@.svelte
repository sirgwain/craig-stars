<script lang="ts">
	import LayoutFullWidth from '$lib/components/LayoutFullWidth.svelte';
	import type { Snippet } from 'svelte';
	type Props = {
		children?: Snippet;
	};

	let { children }: Props = $props();
</script>

<LayoutFullWidth>
	{#if children}{@render children()}{:else}content{/if}
</LayoutFullWidth>
