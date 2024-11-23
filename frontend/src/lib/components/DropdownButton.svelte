<script lang="ts" generics="T extends object">
	import { ChevronDown } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	type Props = {
		title: string;
		items: T[];
		itemTitle?: (item: T) => string;
		onSelected?: (item: T) => void;
	};

	let { title, items, itemTitle = (i) => `${i}`, onSelected }: Props = $props();

	function onSelect(item: T) {
		(document.activeElement as HTMLElement)?.blur();
		onSelected?.(item);
	}

	let divRef: HTMLDivElement;
</script>

<div class="dropdown dropdown-top dropdown-end">
	<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
	<!-- svelte-ignore a11y_label_has_associated_control -->
	<label tabindex="0" class="btn btn-outline btn-sm btn-secondary normal-case"
		>{title} <Icon src={ChevronDown} size="16" class="hover:stroke-accent" /></label
	>
	<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
	<div
		tabindex="0"
		class="shadow menu dropdown-content bg-base-100 rounded-box w-[13rem] sm:max-h-60 overflow-y-auto absolute"
	>
		<ul>
			{#each items as item}
				<li>
					<button type="button" onclick={() => onSelect(item)}>{itemTitle(item)}</button>
				</li>
			{/each}
		</ul>
	</div>
</div>
