<script lang="ts">
	import { clamp } from '$lib/services/Math';
	import { type HabType, Grav, Temp, Rad } from '$lib/types/cs';
	import { getHabValueString, HabTypeShortString, habTypeString } from '$lib/types/Hab';
	import { draggable, type DragEventData } from '@neodrag/svelte';
	import {
		ChevronDoubleLeft,
		ChevronDoubleRight,
		ChevronLeft,
		ChevronRight
	} from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	type Props = {
		habType: HabType;
		habLow: number | undefined;
		habHigh: number | undefined;
		immune: boolean | undefined;
	};

	let {
		habType,
		habLow = $bindable(),
		habHigh = $bindable(),
		immune = $bindable()
	}: Props = $props();

	let habTypeShortString = $derived(HabTypeShortString[habType]);

	let barContainerRef: HTMLDivElement | undefined = $state();
	let containerWidth = $derived(barContainerRef?.parentElement?.clientWidth ?? 0);

	let habWidth = $derived((habHigh ?? 0) - (habLow ?? 0));
	let position = $derived(
		barContainerRef
			? {
					x: Math.floor(((habLow ?? 0) / 100) * containerWidth),
					y: 0
				}
			: undefined
	);

	const onLeft = () => {
		const width = habWidth;
		habLow = clamp((habLow ?? 0) - 1, 0, 100 - width);
		habHigh = clamp((habHigh ?? 0) - 1, width, 100);
	};

	const onRight = () => {
		const width = habWidth;
		habLow = clamp((habLow ?? 0) + 1, 0, 100 - width);
		habHigh = clamp((habHigh ?? 0) + 1, width, 100);
	};

	const onGrow = () => {
		const width = clamp(habWidth + 2, 20, 100);
		habLow = clamp((habLow ?? 0) - 1, 0, 100 - width);
		habHigh = clamp((habHigh ?? 0) + 1, width, 100);
	};

	const onShrink = () => {
		const width = clamp(habWidth - 2, 20, 100);
		habLow = clamp((habLow ?? 0) + 1, 0, (habHigh ?? 0) - width);
		habHigh = clamp((habHigh ?? 0) - 1, habLow + width, 100);
	};

	const onDrag = (data: DragEventData) => {
		const width = habWidth;
		if (containerWidth && habLow) {
			const pixelOffsetInPercent = Math.floor((data.offsetX / containerWidth) * 100);
			habLow = clamp(pixelOffsetInPercent, 0, 100 - width);
			habHigh = clamp(habLow + width, width, 100);
		}
	};
</script>

<div class="flex flex-col md:flex-row">
	<div class="text-center md:text-right md:w-[5.5rem] h-full my-auto mr-2">
		{habTypeString(habType)}
	</div>
	<div class="grow flex flex-col">
		<div class="flex flex-row h-8">
			<button type="button" onclick={onLeft} class="btn btn-outline btn-sm"
				><Icon src={ChevronLeft} size="20" />
			</button>

			<div class="grow border-b border-base-300 bg-black mx-1 overflow-hidden h-full">
				<div class="h-full" class:hidden={immune} bind:this={barContainerRef}>
					{#if position}
						<div
							use:draggable={{ bounds: 'parent', position, onDrag }}
							style={`width: ${habWidth.toFixed()}%`}
							class="h-full"
							class:grav-bar={habType === Grav}
							class:temp-bar={habType === Temp}
							class:rad-bar={habType === Rad}
						></div>
					{/if}
				</div>
			</div>
			<button
				type="button"
				onclick={onRight}
				class="btn btn-outline btn-sm"
				data-type={`${habTypeShortString}-right-button`}
				><Icon src={ChevronRight} size="20" />
			</button>
		</div>
		<div class="flex flex-row grow mt-2">
			<div>
				<button
					type="button"
					onclick={onGrow}
					class="btn btn-outline btn-sm"
					data-type={`${habTypeShortString}-grow-button`}
					><Icon src={ChevronDoubleLeft} size="20" />
					<Icon src={ChevronDoubleRight} size="20" /></button
				>
			</div>
			<div class="grow ml-2">
				<label
					><input type="checkbox" bind:checked={immune} /> Immune to {habTypeString(habType)}</label
				>
			</div>
			<div>
				<button
					type="button"
					onclick={onShrink}
					class="btn btn-outline btn-sm"
					data-type={`${habTypeShortString}-left-button`}
					><Icon src={ChevronDoubleRight} size="20" />
					<Icon src={ChevronDoubleLeft} size="20" /></button
				>
			</div>
		</div>
	</div>
	<div class="flex flex-row gap-1 justify-center md:flex-col md:text-center md:ml-2 md:w-[5rem]">
		<div class:hidden={immune}>{getHabValueString(habType, habLow ?? 0)}</div>
		<div class:hidden={immune}>to</div>
		<div class:hidden={immune}>{getHabValueString(habType, habHigh ?? 0)}</div>
	</div>
</div>
