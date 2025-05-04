<script lang="ts">
	import { type HabType, Grav, Temp, Rad } from '$lib/types/cs';
	import { clamp } from '$lib/services/Math';

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

	// immune can be undefined by default
	let isImmune = $derived(!!immune);

	let container: HTMLDivElement | undefined = $state();
	let width: number | undefined = $state();
	let height: number | undefined = $state();

	let insetHeight = $derived(height != null ? height - 4 : undefined);
	let insetWidth = $derived(width != null ? width - 4 : undefined);

	let low = $derived(
		insetWidth != null && habLow != null ? Math.trunc(insetWidth * (habLow / 100)) : undefined
	);
	let high = $derived(
		insetWidth != null && habHigh != null ? Math.trunc(insetWidth * (habHigh / 100)) : undefined
	);
	let actualWidth = $derived(
		insetWidth != null && low != null && high != null ? high - low : undefined
	);

	let mouseDrag = false;
	let mouseStart: number;
	let mouseDelta: number = $state(0);
	let habLowDrag: number | undefined;
	let habHighDrag: number | undefined;

	$effect(() => {
		if (!insetWidth) {
			return;
		}

		// delta in percentage
		const delta = Math.trunc((mouseDelta / insetWidth) * 100);

		console.log('mouse delta', mouseDelta, delta);

		if (mouseDrag && habLowDrag != null && habHighDrag != null && !isImmune) {
			if (habHigh != null && habLow != null) {
				const w = habHigh - habLow;
				if (habLowDrag + delta < 0) {
					habLow = 0;
					habHigh = w;
					return;
				} else if (habHighDrag + delta > 100) {
					habLow = 100 - w;
					habHigh = 100;
					return;
				}
			}

			habLow = clamp(habLowDrag + delta, 0, 100);
			habHigh = clamp(habHighDrag + delta, 0, 100);

			console.log('habLow', habLow);
		}
	});

	const onmousedown = (mouseEvent: MouseEvent) => {
		console.log('mouseDown', mouseEvent);
		mouseDrag = true;
		mouseStart = mouseEvent.clientX;
		mouseDelta = 0;
		habLowDrag = habLow;
		habHighDrag = habHigh;

		mouseEvent.preventDefault();
	};

	const onmouseup = (mouseEvent: MouseEvent) => {
		console.log('mouseUp', mouseEvent);
		mouseDrag = false;
	};

	const onmousemove = (mouseEvent: MouseEvent) => {
		if (mouseDrag) {
			console.log('mouseMove', mouseEvent);
			mouseDelta = mouseEvent.clientX - mouseStart;
			console.log(mouseDelta);
		}
	};

	const resizeObserver = new ResizeObserver((entries) => {
		width = entries[0].contentRect.width;
		height = entries[0].contentRect.height;
	});

	$effect(() => {
		if (container) {
			resizeObserver.disconnect();
			resizeObserver.observe(container);
		}
	});

	// $effect(() => {
	// 	console.log('actualWidth', actualWidth);
	// 	console.log('low', low);
	// 	console.log('high', high);
	// 	console.log('insetWidth', insetWidth);
	// 	console.log('insetHeight', insetHeight);
	// 	console.log('width', width);
	// 	console.log('height', height);
	// 	console.log('habLow', habLow);
	// 	console.log('habHigh', habHigh);
	// 	console.log('immune', immune);
	// 	console.log('habType', habType);
	// 	console.log('----------');
	// });
</script>

<svelte:document {onmouseup} {onmousemove} />

<div bind:this={container} class="grow px-1 overflow-hidden h-full">
	<svg {width} {height} viewBox={`0 0 ${width} ${height}`}>
		<rect x="0" y="0" {width} {height} fill="black"></rect>
		{#if actualWidth != null && low != null && !isImmune}
			<rect
				class="cursor-pointer"
				x={low + 2}
				y="2"
				width={actualWidth}
				height={insetHeight}
				fill="white"
				{onmousedown}
				{onmouseup}
				{onmousemove}
				role="menu"
				tabindex="-1"
				class:grav-bar={habType === Grav}
				class:temp-bar={habType === Temp}
				class:rad-bar={habType === Rad}
			/>
		{/if}
	</svg>
</div>
