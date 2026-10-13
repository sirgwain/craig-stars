<script lang="ts">
	import { Grav, Rad, Temp, type HabType } from '#lib/types/Hab.js';
	import { getHabValueString, habTypeString } from '#lib/types/Hab.js';

	type Props = {
		habType: HabType;
		habLow: number | undefined;
		habHigh: number | undefined;
		immune: boolean | undefined;
		compact?: boolean;
	};

	let { habType, habLow, habHigh, immune, compact = false }: Props = $props();

	let habWidth = $derived((habHigh ?? 0) - (habLow ?? 0));
</script>

{#if compact}
	<!-- one row per hab; an immune hab fills the whole bar, dimmed, rather than leaving it empty -->
	<div class="flex items-center gap-2">
		<div class="w-20 shrink-0 text-base-content/70">{habTypeString(habType)}</div>
		<div class="grow min-w-0 bg-black overflow-hidden h-4 p-0.5">
			<div
				style={immune ? undefined : `width: ${habWidth.toFixed()}%; left: ${habLow}%;`}
				class="relative h-full"
				class:w-full={immune}
				class:opacity-40={immune}
				class:grav-bar={habType === Grav}
				class:temp-bar={habType === Temp}
				class:rad-bar={habType === Rad}
			></div>
		</div>
		<div class="w-28 shrink-0 text-right tabular-nums text-base-content/70">
			{#if immune}
				Immune
			{:else}
				{getHabValueString(habType, habLow ?? 0)} to {getHabValueString(habType, habHigh ?? 0)}
			{/if}
		</div>
	</div>
{:else}
	<div class="flex flex-col md:flex-row">
		<div class="text-center md:text-right md:w-32 h-full my-auto mr-2">
			{habTypeString(habType)}
		</div>
		<div class="grow flex flex-col">
			<div class="flex flex-row h-8 my-auto">
				<div class="grow bg-black mx-1 overflow-hidden h-full p-0.5">
					{#if !immune}
						<div
							style={`width: ${habWidth.toFixed()}%; left: ${habLow}%;`}
							class="relative h-full"
							class:grav-bar={habType === Grav}
							class:temp-bar={habType === Temp}
							class:rad-bar={habType === Rad}
						></div>
					{/if}
				</div>
			</div>
		</div>
		<div class="flex flex-row gap-1 justify-center md:flex-col md:text-center md:ml-2 md:w-32">
			{#if immune}
				Immune
			{:else}
				{getHabValueString(habType, habLow ?? 0)}
				to
				{getHabValueString(habType, habHigh ?? 0)}
			{/if}
		</div>
	</div>
{/if}
