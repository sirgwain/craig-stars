<script lang="ts">
	import { WaypointTaskTransportAction } from '$lib/types/cs-proto';
	import { enumToString } from '$lib/types/Enums';
	import TransportActionIcon from './TransportActionIcon.svelte';

	type Props = {
		action: WaypointTaskTransportAction | undefined;
		amount: number | undefined;
		units?: string;
		title: string;
		titleTextClass?: string;
	};

	let { action, amount, units = '', title, titleTextClass = '' }: Props = $props();
</script>

<span class="inline-block">
	{#if (action ?? WaypointTaskTransportAction.UNSPECIFIED) !== WaypointTaskTransportAction.UNSPECIFIED}
		<div class="flex flex-row">
			<div class={`text-right font-semibold mr-2 w-28 ${titleTextClass}`}>{title}</div>
			<div>
				{enumToString(
					WaypointTaskTransportAction,
					action ?? WaypointTaskTransportAction.UNSPECIFIED
				)}
				{#if amount != undefined}
					{amount}{units}
				{/if}
				<TransportActionIcon {action} />
			</div>
		</div>
	{/if}
</span>
