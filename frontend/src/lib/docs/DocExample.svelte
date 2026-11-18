<script lang="ts">
	import type { Snippet } from 'svelte';
	let { title, onReset, children }: { title: string; onReset: () => void; children: Snippet } =
		$props();
</script>

<section
	aria-label={title}
	class="my-6 min-w-0 rounded-box border border-base-300 bg-base-200 text-base/normal"
>
	<div class="flex items-center justify-between gap-3 border-b border-base-300 px-4 py-3">
		<div>
			<h3 class="text-xl font-semibold">{title}</h3>
			<p class="mt-1 text-secondary">Interactive example · Changes stay on this page.</p>
		</div>
		<button type="button" class="btn btn-outline btn-sm shrink-0" onclick={onReset}>Reset</button>
	</div>
	<div class="min-w-0 bg-base-100 p-3 sm:p-5">{@render children()}</div>
</section>
