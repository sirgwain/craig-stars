<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		title?: string;
		// small buttons right after the title
		actions?: Snippet;
		children?: Snippet;
	};

	let { title = '', actions, children }: Props = $props();
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
	class="desk:w-56 phone-landscape:w-56 card bg-base-200 shadow-sm rounded-xs border-2 border-base-300 select-none"
	data-type="command-tile"
	data-id={title}
	oncontextmenu={(e) => e.preventDefault()}
>
	<div class="card-body p-3 gap-0">
		<!-- actions sit right beside the title, and the two are centered together -->
		<div class="flex flex-row items-center justify-center gap-1.5 mb-1">
			<button class="min-w-0" class:w-full={!actions}>
				<div class="truncate text-center text-lg font-semibold text-secondary" {title}>
					{title}
				</div>
			</button>
			{#if actions}
				<div class="flex-none flex flex-row items-center">
					{@render actions()}
				</div>
			{/if}
		</div>
		{@render children?.()}
	</div>
</div>
