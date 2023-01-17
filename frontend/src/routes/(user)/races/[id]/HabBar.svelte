<script lang="ts">
	import { draggable, type DragEventData } from '@neodrag/svelte';
	import { getHabValueString, HabType } from '$lib/types/Hab';
	import { Icon } from '@steeze-ui/svelte-icon';
	import {
		ChevronLeft,
		ChevronRight,
		ChevronDoubleLeft,
		ChevronDoubleRight
	} from '@steeze-ui/heroicons';
	import { clamp } from '$lib/services/Math';

	export let habType: HabType;
	export let habLow: number | undefined;
	export let habHigh: number | undefined;
	export let immune: boolean | undefined;

	$: dragHabLow = habLow ?? 0;
	$: dragHabHigh = habHigh ?? 0;
	$: habWidth = (habHigh ?? 0) - (habLow ?? 0);
	$: habLowPercent = clamp(((habLow ?? 0) / 100) * 100, 0, 100);
	$: habWidthPercent = clamp((habWidth / 100) * 100, 0, 100);

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
		const width = clamp(habWidth + 2, 0, 100);
		habLow = clamp((habLow ?? 0) - 1, 0, 100 - width);
		habHigh = clamp((habHigh ?? 0) + 1, width, 100);
	};

	const onShrink = () => {
		const width = clamp(habWidth - 2, 2, 100);
		habLow = clamp((habLow ?? 0) + 1, 0, 100 - width);
		habHigh = clamp((habHigh ?? 0) - 1, width, 100);
	};

	const onDrag = (data: DragEventData) => {
		const width = habWidth;
		dragHabLow = clamp(
			(habLow ?? 0) + Math.floor((data.offsetX / data.rootNode.clientWidth) * 100),
			0,
			100 - width
		);
		dragHabHigh = clamp(dragHabLow + width, width, 100);
	};

	const onDragEnd = (data: DragEventData) => {
		const width = habWidth;
		habLow = clamp(
			(habLow ?? 0) + Math.floor((data.offsetX / data.rootNode.clientWidth) * 100),
			0,
			100 - width
		);
		habHigh = clamp(habLow + width, width, 100);
		dragHabLow = habLow;
		dragHabHigh = habHigh;
	};
</script>

<div class="flex flex-col md:flex-row">
	<div class="text-center md:text-right md:w-[5.5rem] h-full my-auto mr-2">{habType}</div>
	<div class="grow flex flex-col">
		<div class="flex flex-row h-8">
			<button on:click|preventDefault={() => onLeft()} class="btn btn-outline btn-sm"
				><Icon src={ChevronLeft} size="20" />
			</button>

			<div class="grow border-b border-base-300 bg-black mx-1 overflow-hidden h-full">
				<div class="relative h-full" class:hidden={immune}>
					<div
						use:draggable={{ bounds: 'parent' }}
						on:neodrag={(e) => onDrag(e.detail)}
						on:neodrag:end={(e) => onDragEnd(e.detail)}
						style={`left: ${habLowPercent.toFixed()}%; width: ${habWidthPercent.toFixed()}%`}
						class="absolute h-full"
						class:grav-bar={habType === HabType.Gravity}
						class:temp-bar={habType === HabType.Temperature}
						class:rad-bar={habType === HabType.Radiation}
					/>
				</div>
			</div>
			<button on:click|preventDefault={() => onRight()} class="btn btn-outline btn-sm"
				><Icon src={ChevronRight} size="20" />
			</button>
		</div>
		<div class="flex flex-row grow mt-2">
			<div>
				<button on:click|preventDefault={() => onGrow()} class="btn btn-outline btn-sm"
					><Icon src={ChevronDoubleLeft} size="20" />
					<Icon src={ChevronDoubleRight} size="20" /></button
				>
			</div>
			<div class="grow ml-2">
				<label><input type="checkbox" bind:checked={immune} /> Immune to {habType}</label>
			</div>
			<div>
				<button on:click|preventDefault={() => onShrink()} class="btn btn-outline btn-sm"
					><Icon src={ChevronDoubleRight} size="20" />
					<Icon src={ChevronDoubleLeft} size="20" /></button
				>
			</div>
		</div>
	</div>
	<div class="flex flex-row gap-1 justify-center md:flex-col md:text-center md:ml-2 md:w-[5rem]">
		<div class:hidden={immune}>{getHabValueString(habType, dragHabLow ?? 0)}</div>
		<div class:hidden={immune}>to</div>
		<div class:hidden={immune}>{getHabValueString(habType, dragHabHigh ?? 0)}</div>
	</div>
</div>
