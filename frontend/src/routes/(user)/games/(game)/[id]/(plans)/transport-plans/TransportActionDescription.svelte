<script lang="ts">
	import { WaypointTaskTransportAction } from '$lib/types/Fleet';
	import { startCase } from 'lodash-es';
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
	{#if (action ?? WaypointTaskTransportAction.None) != WaypointTaskTransportAction.None}
		<div class="flex flex-row">
			<div class={`text-right font-semibold mr-2 w-28 ${titleTextClass}`}>{title}</div>
			<div>
				{startCase(action)}
				{#if amount != undefined}
					{amount}{units}
				{/if}
				<TransportActionIcon {action} />
			</div>
		</div>
	{/if}
</span>
