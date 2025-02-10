<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		children?: Snippet;
	};

	let { children }: Props = $props();
</script>

<!-- add a margin on small screens so it will scroll -->
<div class="mx-auto md:max-w-2xl mb-10 md:mb-0">
	{#if children}{@render children()}{:else}This is the main content{/if}
</div>
