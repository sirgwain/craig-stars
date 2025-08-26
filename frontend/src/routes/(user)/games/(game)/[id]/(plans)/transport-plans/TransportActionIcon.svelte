<script lang="ts">
	import { WaypointTaskTransportAction } from '$lib/types/cs-proto';
	import { enumToString } from '$lib/types/Enums';
	import { isLoadAction, isUnloadAction } from '$lib/types/Fleet';
	import { ArrowDown, ArrowUp, ArrowsUpDown, XMark } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	type Props = {
		action: WaypointTaskTransportAction | undefined;
	};

	let { action }: Props = $props();
</script>

<span class="inline-block my-auto">
	<div
		class="tooltip tooltip-left lg:tooltip-top"
		data-tip={`${enumToString(WaypointTaskTransportAction, action ?? WaypointTaskTransportAction.UNSPECIFIED)}`}
	>
		{#if action && isLoadAction(action)}
			<Icon src={ArrowUp} size="16" class="hover:stroke-accent stroke-2" />
		{:else if action && isUnloadAction(action)}
			<Icon src={ArrowDown} size="16" class="hover:stroke-accent stroke-2" />
		{:else if action == WaypointTaskTransportAction.WAIT_FOR_PERCENT || action == WaypointTaskTransportAction.SET_AMOUNT_TO || action == WaypointTaskTransportAction.SET_WAYPOINT_TO}
			<Icon src={ArrowsUpDown} size="16" class="hover:stroke-accent stroke-2" />
		{:else}
			<Icon src={XMark} size="16" class="hover:stroke-accent stroke-2" />
		{/if}
	</div>
</span>
