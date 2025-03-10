<script lang="ts">
	import { ArrowLongLeft, ArrowLongRight } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import type { HTMLAttributes } from 'svelte/elements';
	type Props = {
		onTransferToSource: () => void;
		onTransferToDest: () => void;
	} & HTMLAttributes<HTMLDivElement>;

	let {
		onTransferToSource: onTransferToSource,
		onTransferToDest: onTransferToDest,
		...rest
	}: Props = $props();
</script>

<div class="flex justify-center {rest.class}">
	<button
		onclick={onTransferToSource}
		class="btn btn-outline btn-xs normal-case btn-secondary inline-block p-1"
		data-id={rest['data-id']}
		data-type="transfer-to-source-button"
		><Icon src={ArrowLongLeft} size="16" class="hover:stroke-accent inline" />
	</button>
	<button
		onclick={onTransferToDest}
		class="btn btn-outline btn-xs normal-case btn-secondary inline-block p-1"
		data-id={rest['data-id']}
		data-type="transfer-to-dest-button"
		><Icon src={ArrowLongRight} size="16" class="hover:stroke-accent inline" /></button
	>
</div>
