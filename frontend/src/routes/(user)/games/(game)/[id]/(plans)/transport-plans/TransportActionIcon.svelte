<script lang="ts">
	import { WaypointTaskTransportAction, isLoadAction, isUnloadAction } from '$lib/types/Fleet';
	import { ArrowDown, ArrowUp, ArrowsUpDown, XMark } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { startCase } from 'lodash-es';

	export let action: WaypointTaskTransportAction | undefined;
</script>

<span class="inline-block my-auto">
	<div class="tooltip" data-tip={`${startCase(action)}`}>
		{#if action && isLoadAction(action)}
			<Icon src={ArrowUp} size="16" class="hover:stroke-accent stroke-2" />
		{:else if action && isUnloadAction(action)}
			<Icon src={ArrowDown} size="16" class="hover:stroke-accent stroke-2" />
		{:else if action == WaypointTaskTransportAction.WaitForPercent || action == WaypointTaskTransportAction.SetAmountTo || action == WaypointTaskTransportAction.SetWaypointTo}
			<Icon src={ArrowsUpDown} size="16" class="hover:stroke-accent stroke-2" />
		{:else}
			<Icon src={XMark} size="16" class="hover:stroke-accent stroke-2" />
		{/if}
	</div>
</span>
