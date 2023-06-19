<script lang="ts">
	import { ChevronDown } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	export let title: string;
	export let items: any[];
	export let itemTitle: (item: any) => string = (i) => `${i}`;

	let detailsRef: HTMLDetailsElement;
</script>

<details class="dropdown dropdown-top dropdown-end" bind:this={detailsRef}>
	<summary class="btn btn-outline btn-sm btn-secondary normal-case"
		>{title} <Icon src={ChevronDown} size="16" class="hover:stroke-accent" /></summary
	>
	<div>
		<div
			class="shadow menu dropdown-content bg-base-100 rounded-box w-[13rem] sm:max-h-60 overflow-y-auto absolute"
		>
			<ul>
				{#each items as item}
					<li>
						<button
							type="button"
							on:click={() => {
								detailsRef?.removeAttribute('open');
								dispatch('selected', item);
							}}>{itemTitle(item)}</button
						>
					</li>
				{/each}
			</ul>
		</div>
	</div>
</details>
