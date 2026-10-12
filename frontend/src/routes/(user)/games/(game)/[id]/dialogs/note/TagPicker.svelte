<script lang="ts">
	import TagBadge from '#lib/components/game/TagBadge.svelte';
	import type { TagCount } from '#lib/services/Universe.js';
	import { MaxTagLength, MaxTags } from '#lib/types/Consts.js';
	import { cleanTag, sameTag } from '#lib/types/Tags.js';
	import { Check, Plus } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	type Props = {
		// the tags on the map object
		tags: string[];
		// the tags the player has already used on other map objects
		known: TagCount[];
		disabled?: boolean;
	};

	let { tags = $bindable(), known, disabled = false }: Props = $props();

	let search = $state('');
	let query = $derived(cleanTag(search));

	// every tag the player can toggle: the ones they've used before, most used first, followed by
	// any new ones they've added to this map object
	let all = $derived([
		...known.map((k) => k.tag),
		...tags.filter((t) => !known.some((k) => sameTag(k.tag, t)))
	]);
	let matches = $derived(
		query ? all.filter((t) => t.toLowerCase().includes(query.toLowerCase())) : all
	);
	let canCreate = $derived(!!query && !all.some((t) => sameTag(t, query)));
	let full = $derived(tags.length >= MaxTags);

	function has(tag: string) {
		return tags.some((t) => sameTag(t, tag));
	}

	function toggle(tag: string) {
		if (disabled) {
			return;
		}
		if (has(tag)) {
			tags = tags.filter((t) => !sameTag(t, tag));
		} else if (!full) {
			tags = [...tags, tag];
		}
	}

	// enter toggles the tag that was typed, creating it if it's new, or the only match
	function onKeyDown(e: KeyboardEvent) {
		if (e.key !== 'Enter' || e.metaKey || e.ctrlKey) {
			return;
		}
		e.preventDefault();
		e.stopPropagation();
		if (!query) {
			return;
		}
		const exact = all.find((t) => sameTag(t, query));
		const tag = exact ?? (matches.length === 1 ? matches[0] : query);
		toggle(tag);
		search = '';
	}
</script>

<div class="flex flex-col gap-2" data-type="tag-picker">
	{#if !disabled}
		<input
			type="text"
			class="input input-sm w-full text-base"
			bind:value={search}
			onkeydown={onKeyDown}
			maxlength={MaxTagLength}
			placeholder={all.length ? 'Find or create a tag' : 'Create a tag, i.e. "target"'}
			aria-label="Find or create a tag"
			autocomplete="off"
			autocapitalize="off"
			data-type="tag-search"
		/>
	{/if}
	<div
		class="flex flex-row flex-wrap gap-1 max-h-16 md:max-h-32 overflow-y-auto"
		data-type="tag-list"
	>
		{#each matches as tag (tag.toLowerCase())}
			{@const selected = has(tag)}
			<button
				type="button"
				class="max-w-full rounded-full disabled:opacity-50"
				aria-pressed={selected}
				aria-label={tag}
				disabled={disabled || (!selected && full)}
				onclick={() => toggle(tag)}
				data-type="tag-toggle"
				data-id={tag}
			>
				<TagBadge {tag} outline={!selected} class="badge-md gap-1 cursor-pointer">
					{#if selected}
						<Icon src={Check} class="size-3 flex-none" />
					{/if}
				</TagBadge>
			</button>
		{/each}
		{#if canCreate && !disabled && !full}
			<button
				type="button"
				class="badge badge-md badge-outline gap-1 cursor-pointer max-w-full"
				onclick={() => {
					toggle(query);
					search = '';
				}}
				data-type="tag-create"
			>
				<Icon src={Plus} class="size-3 flex-none" />
				<span class="truncate">Create “{query}”</span>
			</button>
		{/if}
		{#if !matches.length && !canCreate && disabled}
			<span class="text-sm opacity-70">No tags</span>
		{/if}
	</div>
	{#if full}
		<div class="text-xs opacity-70">A map object can have at most {MaxTags} tags</div>
	{/if}
</div>
