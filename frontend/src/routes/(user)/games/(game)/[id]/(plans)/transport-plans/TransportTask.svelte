<script lang="ts">
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import NumberInput from '$lib/components/NumberInput.svelte';
	import { WaypointTaskTransportAction } from '$lib/types/Fleet';

	export let action: WaypointTaskTransportAction | undefined;
	export let amount: number | undefined;
	export let textClass = '';
	export let title: string;
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
	/>
</div>
<div>
	<NumberInput titleClass="hidden" name={`amount${title}`} bind:value={amount} />
</div>
