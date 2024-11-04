<script lang="ts" module>
	export type TransferButtons = {
		'transfer-to-source': MouseEvent;
		'transfer-to-dest': MouseEvent;
	};
</script>

<script lang="ts">
	import { ArrowLongLeft, ArrowLongRight } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher } from 'svelte';
	interface Props {
		[key: string]: any;
	}

	let { ...props }: Props = $props();

	const dispatch = createEventDispatcher();
</script>

<div class="flex justify-center {props.class}">
	<button
		onclick={(e) => dispatch('transfer-to-source', e)}
		class="btn btn-outline btn-xs normal-case btn-secondary inline-block p-1"
		><Icon src={ArrowLongLeft} size="16" class="hover:stroke-accent inline" />
	</button>
	<button
		onclick={(e) => dispatch('transfer-to-dest', e)}
		class="btn btn-outline btn-xs normal-case btn-secondary inline-block p-1"
		><Icon src={ArrowLongRight} size="16" class="hover:stroke-accent inline" /></button
	>
</div>
