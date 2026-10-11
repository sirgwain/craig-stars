<!--
  @component
  The minefields button. Like the original, it opens a menu for picking whose minefields the scanner shows.
 -->
<script lang="ts">
	import Minefield from '#lib/components/icons/Minefield.svelte';
	import { getGameContext } from '#lib/services/GameContext.js';
	import {
		getShownMinefieldOwners,
		MinefieldOwners,
		type MinefieldOwner
	} from '#lib/types/Filter.js';

	const { settings } = getGameContext();

	const labels: Record<MinefieldOwner, string> = {
		mine: 'Mine Fields of Yours',
		friends: 'Mine Fields of Friends',
		neutrals: 'Mine Fields of Neutrals',
		enemies: 'Mine Fields of Enemies'
	};

	let shown = $derived(getShownMinefieldOwners($settings));

	function show(owners: MinefieldOwner[]) {
		settings.update((s) => {
			s.setShownMinefieldOwners(owners);
			return s;
		});
	}

	function toggle(owner: MinefieldOwner) {
		// keep the owners in the order they are listed
		show(MinefieldOwners.filter((o) => shown.includes(o) !== (o === owner)));
	}
</script>

<div class="dropdown size-full">
	<div
		tabindex="0"
		role="button"
		title="Mine Fields Overlay"
		aria-label="Mine Fields Overlay"
		data-filter="minefields"
		class:btn-primary={shown.length > 0}
		class:btn-ghost={shown.length === 0}
		class="btn btn-xs size-full p-0"
	>
		<Minefield class="size-6 " />
	</div>
	<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
	<ul
		tabindex="0"
		class="dropdown-content menu menu-sm flex-nowrap z-20 w-56 p-1 bg-base-200 rounded-box border border-base-300 shadow-sm"
	>
		{#snippet item(label: string, checked: boolean, onclick: () => void)}
			<li>
				<button type="button" role="menuitemcheckbox" aria-checked={checked} {onclick}>
					<span class="w-4">{checked ? '✓' : ''}</span>
					{label}
				</button>
			</li>
		{/snippet}
		{@render item('All Mine Fields', shown.length === MinefieldOwners.length, () =>
			show([...MinefieldOwners])
		)}
		{@render item('No Mine Fields', shown.length === 0, () => show([]))}
		<li></li>
		{#each MinefieldOwners as owner (owner)}
			{@render item(labels[owner], shown.includes(owner), () => toggle(owner))}
		{/each}
	</ul>
</div>
