<script lang="ts">
	interface Props {
		title?: string;
		children?: import('svelte').Snippet;
	}

	let { title = '', children }: Props = $props();
</script>

<h3 class="text-xl font-semibold mb-2 mt-4 border border-secondary rounded-md p-1 pl-2">
	{title}
	{@render children?.()}
</h3>
