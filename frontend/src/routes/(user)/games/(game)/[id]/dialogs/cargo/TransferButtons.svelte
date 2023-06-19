<script lang="ts">
	import { ArrowLongLeft, ArrowLongRight } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();
</script>

<div class="flex justify-center {$$props.class}">
	<button
		on:click={() => dispatch('transfer-to-source')}
		class="btn btn-outline btn-xs normal-case btn-secondary inline-block p-1"
		><Icon src={ArrowLongLeft} size="16" class="hover:stroke-accent inline" />
	</button>
	<button
		on:click={() => dispatch('transfer-to-dest')}
		class="btn btn-outline btn-xs normal-case btn-secondary inline-block p-1"
		><Icon src={ArrowLongRight} size="16" class="hover:stroke-accent inline" /></button
	>
</div>
