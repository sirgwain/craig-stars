<!--
  @component
  The notes button. It opens a menu for picking which of the player's notes the scanner shows: labels
  placed on the map, and notes on planets by their tags.
 -->
<script lang="ts">
	import { getGameContext } from '#lib/services/GameContext.js';
	import { NoteKindLabel, NoteKindUntagged, noteKindOfTag } from '#lib/types/Filter.js';
	import { PencilSquare } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	const { settings, universe } = getGameContext();

	// every kind of note that can be shown, the player's tags after the kinds that are always there
	let kinds = $derived([
		{ kind: NoteKindLabel, label: 'Map Labels' },
		{ kind: NoteKindUntagged, label: 'Notes Without Tags' },
		...$universe.tags.map(({ tag }) => ({ kind: noteKindOfTag(tag), label: tag }))
	]);

	// kinds that were hidden stay hidden if their tag goes away and comes back
	let hidden = $derived($settings.showNotes ? $settings.hiddenNoteKinds : kinds.map((k) => k.kind));
	let shown = $derived(kinds.filter((k) => !hidden.includes(k.kind)));

	function hide(kinds: string[], showNotes = true) {
		settings.update((s) => {
			s.setHiddenNoteKinds(kinds, showNotes);
			return s;
		});
	}

	function toggle(kind: string) {
		if (!$settings.showNotes) {
			// picking one turns notes back on, with only that kind showing
			hide(kinds.map((k) => k.kind).filter((k) => k !== kind));
		} else if (hidden.includes(kind)) {
			hide(hidden.filter((k) => k !== kind));
		} else {
			hide([...hidden, kind]);
		}
	}
</script>

<div class="dropdown size-full">
	<div
		tabindex="0"
		role="button"
		title="Notes Overlay"
		aria-label="Notes Overlay"
		data-filter="notes"
		class:btn-primary={shown.length > 0}
		class:btn-ghost={shown.length === 0}
		class="btn btn-xs size-full p-0"
	>
		<Icon src={PencilSquare} class="size-6" />
	</div>
	<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
	<ul
		tabindex="0"
		class="dropdown-content menu menu-sm flex-nowrap z-20 w-56 max-h-96 overflow-y-auto p-1 bg-base-200 rounded-box border border-base-300 shadow-sm"
	>
		{#snippet item(label: string, checked: boolean, onclick: () => void)}
			<li>
				<button type="button" role="menuitemcheckbox" aria-checked={checked} {onclick}>
					<span class="w-4 flex-none">{checked ? '✓' : ''}</span>
					<span class="truncate">{label}</span>
				</button>
			</li>
		{/snippet}
		{@render item('All Notes', shown.length === kinds.length, () => hide([]))}
		{@render item('No Notes', shown.length === 0, () => hide($settings.hiddenNoteKinds, false))}
		<li></li>
		{#each kinds as { kind, label } (kind)}
			{@render item(label, !hidden.includes(kind), () => toggle(kind))}
		{/each}
	</ul>
</div>
