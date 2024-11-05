<script lang="ts">
	interface Props {
		children?: import('svelte').Snippet;
	}

	let { children }: Props = $props();
</script>

<div class="w-full mx-auto md:max-w-2xl">
	{@render children?.()}
</div>
