<script lang="ts">
	import Menu from '$lib/components/Menu.svelte';
	import { me } from '$lib/services/Stores';
	import type { Snippet } from 'svelte';
	type Props = {
		children?: Snippet;
	};

	let { children }: Props = $props();
</script>

<main class="p-3 flex flex-col">
	<div class="sticky top-0 z-10">
		<Menu user={$me} />
	</div>
	<div class="relative">
		<div class="w-full mx-auto md:max-w-3xl">
			{@render children?.()}
		</div>
	</div>
</main>
