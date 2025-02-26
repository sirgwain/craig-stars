<script lang="ts">
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import NumberInput from '$lib/components/NumberInput.svelte';
	import {
		TransportActionLoadAll,
		TransportActionLoadDunnage,
		TransportActionLoadOptimal,
		TransportActionNone,
		TransportActionUnloadAll,
		type WaypointTaskTransportAction
	} from '$lib/types/cs';
	import { WaypointTransportTaskActions } from '$lib/types/Fleet';
	import { startCase } from 'lodash-es';

	type Props = {
		action: WaypointTaskTransportAction | undefined;
		amount: number | undefined;
		textClass?: string;
		title: string;
	};

	let { action = $bindable(), amount = $bindable(), textClass = '', title }: Props = $props();
</script>

<div class={textClass}>
	<div class="label"><span class="w-32 text-right">{title}</span></div>
</div>
<div class="col-span-2">
	<EnumSelect
		name={`action${title}`}
		options={WaypointTransportTaskActions}
		bind:value={action}
		titleClass="hidden"
		typeTitle={(value) => (!value || value === TransportActionNone ? 'None' : startCase(value))}
		showEmpty={true}
	/>
</div>
<div>
	<NumberInput
		titleClass="hidden"
		name={`amount${title}`}
		bind:value={amount}
		disabled={action == undefined ||
			action === TransportActionNone ||
			action === TransportActionLoadAll ||
			action === TransportActionUnloadAll ||
			action === TransportActionLoadDunnage ||
			action === TransportActionLoadOptimal}
	/>
</div>
