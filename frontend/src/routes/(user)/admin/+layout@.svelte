<script lang="ts">
	import LayoutFullWidth from '$lib/components/LayoutFullWidth.svelte';
	interface Props {
		children?: import('svelte').Snippet;
	}

	let { children }: Props = $props();
</script>

<LayoutFullWidth>
	{#if children}{@render children()}{:else}content{/if}
</LayoutFullWidth>
