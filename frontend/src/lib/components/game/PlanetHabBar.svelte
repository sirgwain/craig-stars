<script lang="ts">
	import { clamp } from '$lib/services/Math';
	import type { HabType } from '$lib/types/Hab';
	import { getHabValueString, HabTypeShortString, habTypeString } from '$lib/types/Hab';
	import PlanetBaseHabPoint from './PlanetBaseHabPoint.svelte';
	import PlanetHabPoint from './PlanetHabPoint.svelte';
	import PlanetHabTerraformLine from './PlanetHabTerraformLine.svelte';

	type Props = {
		habType: HabType;
		value: number;
		baseValue: number;
		terraformValue: number;
		high: number;
		low: number;
		immune: boolean;
		onTooltip?: (e: PointerEvent) => void;
	};

	let { habType, value, baseValue, terraformValue, high, low, immune, onTooltip }: Props = $props();

	let width = $derived(high - low);

	let habPointPercent = $derived(clamp((value / 100) * 100, 0, 100));
	let baseHabPercent = $derived(clamp((baseValue / 100) * 100, 0, 100));
	let terraformHabPointPercent = $derived(
		clamp(((baseValue + terraformValue) / 100) * 100, 0, 100)
	);
	let habLowPercent = $derived(clamp((low / 100) * 100, 0, 100));
	let habWidthPercent = $derived(clamp((width / 100) * 100, 0, 100));
</script>

<div class="flex flex-row" class:cursor-help={!!onTooltip} onpointerdown={onTooltip}>
	<div class="text-right w-[5.5rem] text-tile-item-title">{habTypeString(habType)}</div>
	<div class="grow border-b border-base-300 bg-black mx-1 overflow-hidden">
		<div class="h-full relative">
			{#if !immune}
				<div
					style={`left: ${habLowPercent.toFixed()}%; width: ${habWidthPercent?.toFixed()}%`}
					class={`absolute h-full ${HabTypeShortString[habType]}-bar`}
				></div>
			{/if}
			<PlanetHabPoint
				style={`left: ${habPointPercent.toFixed()}%;`}
				class={`absolute h-full -translate-x-1/2 ${HabTypeShortString[habType]}-point`}
			/>
			<PlanetBaseHabPoint
				style={`left: ${baseHabPercent.toFixed()}%;`}
				class={`absolute h-full -translate-x-1/2 ${HabTypeShortString[habType]}-point`}
			/>
			<!-- Terraform line -->
			<PlanetHabTerraformLine
				x1={habPointPercent}
				x2={terraformHabPointPercent}
				class={`${HabTypeShortString[habType]}-point`}
			/>
		</div>
	</div>
	<div class="w-[3rem]">{getHabValueString(habType, value)}</div>
</div>
