<script lang="ts">
	import { ChevronDown } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	export let title: string;
	export let items: any[];
	export let itemTitle: (item: any) => string = (i) => `${i}`;

	function onSelect(item: any) {
		(document.activeElement as HTMLElement)?.blur();
		dispatch('selected', item);
	}

	let divRef: HTMLDivElement;
</script>

<div class="dropdown dropdown-top dropdown-end">
	<!-- svelte-ignore a11y-no-noninteractive-tabindex -->
	<!-- svelte-ignore a11y-label-has-associated-control -->
	<label tabindex="0" class="btn btn-outline btn-sm btn-secondary normal-case"
		>{title} <Icon src={ChevronDown} size="16" class="hover:stroke-accent" /></label
	>
	<!-- svelte-ignore a11y-no-noninteractive-tabindex -->
	<div
		tabindex="0"
		class="shadow menu dropdown-content bg-base-100 rounded-box w-[13rem] sm:max-h-60 overflow-y-auto absolute"
	>
		<ul>
			{#each items as item}
				<li>
					<button type="button" on:click={() => onSelect(item)}>{itemTitle(item)}</button>
				</li>
			{/each}
		</ul>
	</div>
</div>
