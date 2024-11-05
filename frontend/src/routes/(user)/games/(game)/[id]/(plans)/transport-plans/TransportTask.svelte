<script lang="ts">
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import NumberInput from '$lib/components/NumberInput.svelte';
	import { WaypointTaskTransportAction } from '$lib/types/Fleet';
	import { startCase } from 'lodash-es';

	interface Props {
		action: WaypointTaskTransportAction | undefined;
		amount: number | undefined;
		textClass?: string;
		title: string;
	}

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
		typeTitle={(value) =>
			!value || value === WaypointTaskTransportAction.None ? 'None' : startCase(value)}
		showEmpty={true}
	/>
</div>
<div>
	<NumberInput
		titleClass="hidden"
		name={`amount${title}`}
		bind:value={amount}
		disabled={action == undefined ||
			action === WaypointTaskTransportAction.None ||
			action === WaypointTaskTransportAction.LoadAll ||
			action === WaypointTaskTransportAction.UnloadAll ||
			action === WaypointTaskTransportAction.LoadDunnage ||
			action === WaypointTaskTransportAction.LoadOptimal}
	/>
</div>
