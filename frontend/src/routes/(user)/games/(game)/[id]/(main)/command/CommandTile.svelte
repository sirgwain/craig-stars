<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		title?: string;
		children?: Snippet;
	};

	let { title = '', children }: Props = $props();
</script>

<div
	class="md:w-[14rem] card bg-base-200 shadow rounded-sm border-2 border-base-300"
	data-type="command-tile"
	data-id={title}
>
	<div class="card-body p-3 gap-0">
		<button class="w-full">
			<div class="flex flex-row items-center mb-1">
				<div class="flex-1 text-center text-lg font-semibold text-secondary">
					{title}
				</div>
			</div>
		</button>
		{@render children?.()}
	</div>
</div>
