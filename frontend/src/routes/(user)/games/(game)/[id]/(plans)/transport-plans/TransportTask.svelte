<script lang="ts">
	import EnumSelect from '#lib/components/EnumSelect.svelte';
	import NumberInput from '#lib/components/NumberInput.svelte';
	import { WaypointTaskTransportAction } from '#lib/types/cs-proto.js';
	import { enumToString } from '#lib/types/Enums.js';

	type Props = {
		action: WaypointTaskTransportAction | undefined;
		amount: number | undefined;
		textClass?: string;
		title: string;
	};

	let { action = $bindable(), amount = $bindable(), textClass = '', title }: Props = $props();

	// fuel is measured in mg, everything else in kT, and the percent tasks take a percent
	let unit = $derived(
		action === WaypointTaskTransportAction.FILL_PERCENT ||
			action === WaypointTaskTransportAction.WAIT_FOR_PERCENT
			? '%'
			: title === 'Fuel'
				? 'mg'
				: 'kT'
	);

	let description = $derived.by(() => {
		const resource = title.toLowerCase();
		const value = amount ?? 0;
		const quantity =
			title === 'Colonists'
				? (value * 100).toLocaleString()
				: `${value.toLocaleString()} ${title === 'Fuel' ? 'mg' : 'kT'}`;
		const cargoAmount =
			title === 'Colonists' ? `${quantity} colonists` : `${quantity} of ${resource}`;
		const isFuel = title === 'Fuel';
		const capacity = isFuel ? 'fuel tanks' : 'cargo hold';
		// fuel only moves between fleets
		const target = isFuel ? 'target fleet' : 'target';
		const percent = `${value.toLocaleString()}%`;

		switch (action) {
			case WaypointTaskTransportAction.LOAD_OPTIMAL:
				return isFuel
					? 'Gives the target fleet any fuel not needed to finish the remaining waypoints'
					: `Fuel only, does nothing for ${resource}`;
			case WaypointTaskTransportAction.LOAD_ALL:
				return `Loads all the ${resource} the ${target} has, as far as the ${capacity} can hold it`;
			case WaypointTaskTransportAction.UNLOAD_ALL:
				return `Unloads all ${resource} to the ${target}`;
			case WaypointTaskTransportAction.LOAD_AMOUNT:
				return `Loads up to ${cargoAmount} from the ${target}`;
			case WaypointTaskTransportAction.UNLOAD_AMOUNT:
				return `Unloads up to ${cargoAmount} to the ${target}`;
			case WaypointTaskTransportAction.FILL_PERCENT:
				return `Loads ${resource} until it fills ${percent} of the ${capacity}, then moves on even if short`;
			case WaypointTaskTransportAction.WAIT_FOR_PERCENT:
				return isFuel
					? `Loads fuel until it fills ${percent} of the fuel tanks, waiting here until it does`
					: `Loads ${resource} until it fills ${percent} of the cargo hold, waiting here until it does or the hold is full`;
			case WaypointTaskTransportAction.LOAD_DUNNAGE:
				return isFuel
					? 'Loads as much fuel as fits once the other loads are done'
					: `Loads ${resource} into the space left once the other loads are done`;
			case WaypointTaskTransportAction.SET_AMOUNT_TO:
				return `Loads or unloads until the fleet carries ${cargoAmount}, waiting here if the ${target} doesn't have enough`;
			case WaypointTaskTransportAction.SET_WAYPOINT_TO:
				return `Loads or unloads until the ${target} has ${cargoAmount}, without waiting`;
			default:
				return '';
		}
	});
</script>

<div class={textClass}>
	<div class="cs-form-label"><span class="w-32 text-right">{title}</span></div>
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
		{unit}
		unitLabelClass="w-7 pl-1 text-sm"
		bind:value={amount}
		disabled={action == undefined ||
			action === WaypointTaskTransportAction.UNSPECIFIED ||
			action === WaypointTaskTransportAction.LOAD_ALL ||
			action === WaypointTaskTransportAction.UNLOAD_ALL ||
			action === WaypointTaskTransportAction.LOAD_DUNNAGE ||
			action === WaypointTaskTransportAction.LOAD_OPTIMAL}
	/>
</div>
{#if description}
	<div class="col-start-2 col-span-3 ml-2 mb-2 text-sm opacity-70">{description}</div>
{/if}
