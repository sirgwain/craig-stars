<script lang="ts" context="module">
	export type ProductionItemsButtonsEvent = {
		'add-item': MouseEvent;
		'remove-item': MouseEvent;
		'item-up': MouseEvent;
		'item-down': MouseEvent;
		clear: MouseEvent;
	};
</script>

<script lang="ts">
	import {
		ArrowLongRight,
		ArrowLongLeft,
		ArrowLongUp,
		ArrowLongDown,
		XCircle
	} from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher<ProductionItemsButtonsEvent>();
</script>

<div class="flex-none mx-0.5 md:w-32 px-1">
	<div class="flex-row flex-none gap-y-2">
		<button
			type="button"
			on:click|preventDefault={(e) => dispatch('add-item', e)}
			class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
			><span class="hidden sm:inline">Add </span><Icon
				src={ArrowLongRight}
				size="16"
				class="hover:stroke-accent inline"
			/></button
		>
		<button
			type="button"
			on:click|preventDefault={(e) => dispatch('remove-item', e)}
			class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
			><Icon src={ArrowLongLeft} size="16" class="hover:stroke-accent inline" /><span
				class="hidden sm:inline"
			>
				Remove</span
			>
		</button>
		<button
			type="button"
			on:click|preventDefault={(e) => dispatch('item-up', e)}
			class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
			><span class="hidden sm:inline">Item Up </span><Icon
				src={ArrowLongUp}
				size="16"
				class="hover:stroke-accent inline"
			/>
		</button>
		<button
			type="button"
			on:click|preventDefault={(e) => dispatch('item-down', e)}
			class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
			><span class="hidden sm:inline">Item Down </span><Icon
				src={ArrowLongDown}
				size="16"
				class="hover:stroke-accent inline"
			/>
		</button>
		<button
			type="button"
			on:click|preventDefault={(e) => dispatch('clear', e)}
			class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
			><span class="hidden sm:inline">Clear </span><Icon
				src={XCircle}
				size="16"
				class="hover:stroke-accent inline"
			/>
		</button>
	</div>
</div>
