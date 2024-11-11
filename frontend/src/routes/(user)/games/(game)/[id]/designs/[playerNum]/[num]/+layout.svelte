<script lang="ts">
	type Props = {
		children?: import('svelte').Snippet;
	};

	let { children }: Props = $props();
</script>

<div class="w-full mx-auto md:max-w-3xl">
	{@render children?.()}
</div>
