<script lang="ts">
	import Menu from '$lib/components/Menu.svelte';
	import { me } from '$lib/services/Stores';
	type Props = {
		children?: import('svelte').Snippet;
	};

	let { children }: Props = $props();
</script>

<main class="p-3 flex flex-col">
	<div class="sticky top-0 z-10">
		<Menu user={$me} />
	</div>
	<div class="relative">
		<div class="w-full mx-auto">
			{#if children}{@render children()}{:else}This is the main content{/if}
		</div>
	</div>
</main>
