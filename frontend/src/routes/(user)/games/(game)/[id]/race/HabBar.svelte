<script lang="ts">
	import { HabTypes, getHabValueString, type HabType } from '$lib/types/Hab';

	export let habType: HabType;
	export let habLow: number | undefined;
	export let habHigh: number | undefined;
	export let immune: boolean | undefined;

	$: habWidth = (habHigh ?? 0) - (habLow ?? 0);
</script>

<div class="flex flex-col md:flex-row">
	<div class="text-center md:text-right md:w-[5.5rem] h-full my-auto mr-2">{habType}</div>
	<div class="grow flex flex-col">
		<div class="flex flex-row h-8 my-auto">
			<div class="grow border-b border-base-300 bg-black mx-1 overflow-hidden h-full">
				{#if !immune}
					<div
						style={`width: ${habWidth.toFixed()}%; left: ${habLow}%;`}
						class="relative h-full"
						class:grav-bar={habType === HabTypes.Gravity}
						class:temp-bar={habType === HabTypes.Temperature}
						class:rad-bar={habType === HabTypes.Radiation}
					/>
				{/if}
			</div>
		</div>
	</div>
	<div class="flex flex-row gap-1 justify-center md:flex-col md:text-center md:ml-2 md:w-[5rem]">
		{#if immune}
			Immune
		{:else}
			<div>{getHabValueString(habType, habLow ?? 0)}</div>
			<div>to</div>
			<div>{getHabValueString(habType, habHigh ?? 0)}</div>
		{/if}
	</div>
</div>
