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

	const dispatch = createEventDispatcher();
</script>

<div class="flex-none mx-0.5 md:w-32 px-1">
	<div class="flex-row flex-none gap-y-2">
		<button
			type="button"
			on:click={() => dispatch('add-item')}
			class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
			><span class="hidden sm:inline">Add </span><Icon
				src={ArrowLongRight}
				size="16"
				class="hover:stroke-accent inline"
			/></button
		>
		<button
			type="button"
			on:click={() => dispatch('remove-item')}
			class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
			><Icon src={ArrowLongLeft} size="16" class="hover:stroke-accent inline" /><span
				class="hidden sm:inline"
			>
				Remove</span
			>
		</button>
		<button
			type="button"
			on:click={() => dispatch('item-up')}
			class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
			><span class="hidden sm:inline">Item Up </span><Icon
				src={ArrowLongUp}
				size="16"
				class="hover:stroke-accent inline"
			/>
		</button>
		<button
			type="button"
			on:click={() => dispatch('item-down')}
			class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
			><span class="hidden sm:inline">Item Down </span><Icon
				src={ArrowLongDown}
				size="16"
				class="hover:stroke-accent inline"
			/>
		</button>
		<button
			type="button"
			on:click={() => dispatch('clear')}
			class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
			><span class="hidden sm:inline">Clear </span><Icon
				src={XCircle}
				size="16"
				class="hover:stroke-accent inline"
			/>
		</button>
	</div>
</div>
