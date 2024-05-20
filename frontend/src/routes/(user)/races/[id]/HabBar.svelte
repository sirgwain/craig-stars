<script lang="ts">
	import { clamp } from '$lib/services/Math';
	import { HabTypes, getHabValueString, type HabType, habTypeString } from '$lib/types/Hab';
	import { draggable, type DragEventData } from '@neodrag/svelte';
	import {
		ChevronDoubleLeft,
		ChevronDoubleRight,
		ChevronLeft,
		ChevronRight
	} from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	export let habType: HabType;
	export let habLow: number | undefined;
	export let habHigh: number | undefined;
	export let immune: boolean | undefined;

	let barContainerRef: HTMLDivElement | null = null;
	$: containerWidth = barContainerRef?.parentElement?.clientWidth ?? 0;

	$: habWidth = (habHigh ?? 0) - (habLow ?? 0);
	$: position = barContainerRef
		? {
				x: Math.floor(((habLow ?? 0) / 100) * containerWidth),
				y: 0
			}
		: undefined;

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
			<button type="button" on:click|preventDefault={() => onLeft()} class="btn btn-outline btn-sm"
				><Icon src={ChevronLeft} size="20" />
			</button>

			<div class="grow border-b border-base-300 bg-black mx-1 overflow-hidden h-full">
				<div class="h-full" class:hidden={immune} bind:this={barContainerRef}>
					{#if position}
						<div
							use:draggable={{ bounds: 'parent', position }}
							on:neodrag={(e) => onDrag(e.detail)}
							style={`width: ${habWidth.toFixed()}%`}
							class="h-full"
							class:grav-bar={habType === HabTypes.Gravity}
							class:temp-bar={habType === HabTypes.Temperature}
							class:rad-bar={habType === HabTypes.Radiation}
						/>
					{/if}
				</div>
			</div>
			<button type="button" on:click|preventDefault={() => onRight()} class="btn btn-outline btn-sm"
				><Icon src={ChevronRight} size="20" />
			</button>
		</div>
		<div class="flex flex-row grow mt-2">
			<div>
				<button
					type="button"
					on:click|preventDefault={() => onGrow()}
					class="btn btn-outline btn-sm"
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
					on:click|preventDefault={() => onShrink()}
					class="btn btn-outline btn-sm"
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
