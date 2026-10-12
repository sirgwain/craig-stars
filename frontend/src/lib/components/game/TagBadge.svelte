<script lang="ts">
	import { tagHue } from '#lib/types/Tags.js';
	import type { Snippet } from 'svelte';

	type Props = {
		tag: string;
		// an outline instead of a filled badge, for a tag that isn't on the map object
		outline?: boolean;
		class?: string;
		// shown before the tag, i.e. a checkmark
		children?: Snippet;
	};

	let { tag, outline = false, class: className = '', children }: Props = $props();

	let hue = $derived(tagHue(tag));
</script>

<span
	class="badge badge-sm max-w-full border {className}"
	class:text-white={!outline}
	style={outline
		? `border-color: hsl(${hue} 55% 45%); background-color: transparent;`
		: `border-color: hsl(${hue} 55% 34%); background-color: hsl(${hue} 55% 34%);`}
	data-type="tag"
	data-id={tag}
>
	{@render children?.()}
	<span class="truncate">{tag}</span>
</span>
