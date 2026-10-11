<!--
  @component
  The little dropdown beside a fleet filter button, for picking which designs or ship classes the filter shows.
 -->
<script lang="ts" generics="T extends string | number">
	import { ChevronDown } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	type Props = {
		title: string;
		items: { value: T; label: string }[];
		// the items the filter shows
		selected: T[];
		// true if the filter is turned on
		enabled: boolean;
		// open the menu to the left, for filters near the right edge of the screen
		end?: boolean;
	};

	let {
		title,
		items,
		selected = $bindable(),
		enabled = $bindable(),
		end = false
	}: Props = $props();

	// choosing something to show turns the filter on, but nothing in the menu turns it off
	function select(values: T[]) {
		selected = values;
		if (values.length > 0 && !enabled) {
			enabled = true;
		}
	}

	function toggle(value: T) {
		if (selected.includes(value)) {
			selected = selected.filter((v) => v !== value);
		} else {
			select([...selected, value]);
		}
	}
</script>

<div class="dropdown size-full" class:dropdown-end={end}>
	<div
		tabindex="0"
		role="button"
		{title}
		aria-label={title}
		class="btn btn-ghost btn-xs size-full p-0"
	>
		<Icon src={ChevronDown} class="size-3 " />
	</div>
	<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
	<ul
		tabindex="0"
		class="dropdown-content menu menu-sm flex-nowrap z-20 w-48 max-h-96 overflow-y-auto p-1 bg-base-200 rounded-box border border-base-300 shadow-sm"
	>
		<li>
			<button type="button" onclick={() => select(items.map((i) => i.value))}>All Designs</button>
		</li>
		<li>
			<button
				type="button"
				onclick={() =>
					select(items.map((i) => i.value).filter((value) => !selected.includes(value)))}
				>Invert Filter</button
			>
		</li>
		<li><button type="button" onclick={() => select([])}>No Designs</button></li>
		<li></li>
		{#each items as item (item.value)}
			<li>
				<button
					type="button"
					role="menuitemcheckbox"
					aria-checked={selected.includes(item.value)}
					onclick={() => toggle(item.value)}
				>
					<span class="w-4">{selected.includes(item.value) ? '✓' : ''}</span>
					{item.label}
				</button>
			</li>
		{/each}
	</ul>
</div>
