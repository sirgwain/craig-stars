<script lang="ts">
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import NumberInput from '$lib/components/NumberInput.svelte';
	import { WaypointTaskTransportAction } from '$lib/types/cs-proto';
	import { enumToString } from '$lib/types/Enums';

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
		enumType={WaypointTaskTransportAction}
		bind:value={action}
		titleClass="hidden"
		typeTitle={(value) => (!value ? 'None' : enumToString(WaypointTaskTransportAction, value))}
		showEmpty={true}
	/>
</div>
<div>
	<NumberInput
		titleClass="hidden"
		name={`amount${title}`}
		bind:value={amount}
		disabled={action == undefined ||
			action === WaypointTaskTransportAction.UNSPECIFIED ||
			action === WaypointTaskTransportAction.LOAD_ALL ||
			action === WaypointTaskTransportAction.UNLOAD_ALL ||
			action === WaypointTaskTransportAction.LOAD_DUNNAGE ||
			action === WaypointTaskTransportAction.LOAD_OPTIMAL}
	/>
</div>
