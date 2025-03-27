<script lang="ts">
	import { ChevronUp, ChevronDown } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import type { Snippet } from 'svelte';

	type Props = {
		open: boolean;
		children?: Snippet;
		onToggle: () => void;
	};

	let { open, children, onToggle }: Props = $props();
</script>

<!-- Use a whole button so the entire header is clickable -->
<button
	type="button"
	class="cursor-pointer w-full"
	aria-label="show command pane button"
	onclick={onToggle}
>
	<div class="flex flex-row justify-between">
		{@render children?.()}

		{#if open}
			<Icon src={ChevronUp} size="16" class="hover:stroke-accent" />
		{:else}
			<Icon src={ChevronDown} size="16" class="hover:stroke-accent" />
		{/if}
	</div>
</button>
