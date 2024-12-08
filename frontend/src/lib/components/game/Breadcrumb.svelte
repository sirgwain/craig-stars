<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		crumbs?: Snippet;
		end?: Snippet;
	};

	let { crumbs, end }: Props = $props();
</script>

<div class="w-full flex flex-row justify-between border-primary border-b-2 mb-2">
	<div class="breadcrumbs ml-2">
		<ul>
			{@render crumbs?.()}
		</ul>
	</div>
	<div class="grow">
		{@render end?.()}
	</div>
</div>
