<script lang="ts">
	import type { Snippet } from "svelte";

	type Props = {
		children?: Snippet;
	};

	let { children }: Props = $props();
</script>

<div class="w-full mx-auto md:max-w-3xl">
	{@render children?.()}
</div>
