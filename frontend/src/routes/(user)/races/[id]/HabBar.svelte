<script lang="ts">
	import { type HabType, Grav, Temp, Rad } from '$lib/types/cs';
	import { clamp } from '$lib/services/Math';

	type Props = {
		habType: HabType;
		habLow: number | undefined;
		habHigh: number | undefined;
		immune: boolean | undefined;
		onValueChanged?: (low: number, high: number) => void | undefined;
	};

	let { habType, habLow, habHigh, immune, onValueChanged: onValueChanged }: Props = $props();

	// immune can be undefined by default
	let isImmune = $derived(!!immune);

	let container: HTMLDivElement | undefined = $state();
	let width: number | undefined = $state();
	let height: number | undefined = $state();

	// insetMargin specifies how much space to leave around the inset rectangle.
	const insetMargin = 2;

	let insetHeight = $derived(height != null ? height - insetMargin * 2 : undefined);
	let insetWidth = $derived(width != null ? width - insetMargin * 2 : undefined);

	let low = $derived(
		insetWidth != null && habLow != null ? Math.trunc(insetWidth * (habLow / 100)) : undefined
	);
	let high = $derived(
		insetWidth != null && habHigh != null ? Math.trunc(insetWidth * (habHigh / 100)) : undefined
	);
	let actualWidth = $derived(
		insetWidth != null && low != null && high != null ? high - low : undefined
	);

	let pointerDown = false;
	let ref: SVGRectElement | undefined = $state();
	let startValue: number | undefined;
	let habLowStart: number | undefined;
	let habHighStart: number | undefined;

	function onPointerDown(e: PointerEvent) {
		if (e.cancelable) {
			e.preventDefault();
		}

		if (!ref || !container) return;

		pointerDown = true;

		startValue = getPercentFromPointerEvent(e, container);
		habLowStart = habLow;
		habHighStart = habHigh;

		updateValue(getPercentFromPointerEvent(e, container));

		document.body.classList.add('select-none', 'touch-none');
		document.body.classList.remove('touch-manipulation');

		ref.setPointerCapture(e.pointerId);
	}

	function onPointerUp(e: PointerEvent) {
		if (!ref) return;

		document.body.classList.remove('select-none', 'touch-none');
		document.body.classList.add('touch-manipulation');
		pointerDown = false;
		ref.releasePointerCapture(e.pointerId);
	}

	function onPointerMove(e: PointerEvent) {
		if (e.cancelable) {
			e.preventDefault();
		}

		if (pointerDown && container) {
			updateValue(getPercentFromPointerEvent(e, container));
		}
	}

	// This is invoked when the pointer cannot give any more events,
	// such as a touch event where scrolling has started to engage
	function onPointerCancel(e: PointerEvent) {
		if (!ref) return;

		document.body.classList.remove('select-none', 'touch-none');
		document.body.classList.add('touch-manipulation');
		pointerDown = false;
		ref.releasePointerCapture(e.pointerId);
	}

	function getPercentFromPointerEvent(e: PointerEvent, container: HTMLElement): number {
		const boundingRect = container.getBoundingClientRect();

		const x = e.clientX - boundingRect.left + insetMargin;
		const w = boundingRect.width - insetMargin * 2;
		const p = (x / w) * 100;

		return clamp(p, 0, 100);
	}

	function updateValue(x: number) {
		if (startValue != null && habLowStart != null && habHighStart != null) {
			const delta = x - startValue;

			const habWidth = habHighStart - habLowStart;

			let newHabLow = clamp(Math.trunc(habLowStart + delta), 0, 100 - habWidth);
			let newHabHigh = clamp(Math.trunc(newHabLow + habWidth), habWidth, 100);

			onValueChanged?.(newHabLow, newHabHigh);
		}
	}

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
</script>

<svelte:document {onmouseup} {onmousemove} />

<div bind:this={container} class="grow px-1 overflow-hidden h-full">
	<svg {width} {height} viewBox={`0 0 ${width} ${height}`}>
		<rect x="0" y="0" {width} {height} fill="black"></rect>
		{#if actualWidth != null && low != null && !isImmune}
			<rect
				bind:this={ref}
				class="cursor-pointer focus:outline-none"
				x={low + 2}
				y="2"
				width={actualWidth}
				height={insetHeight}
				fill="white"
				onpointerdown={onPointerDown}
				onpointerup={onPointerUp}
				onpointermove={onPointerMove}
				onpointercancel={onPointerCancel}
				role="menu"
				tabindex="-1"
				class:grav-bar={habType === Grav}
				class:temp-bar={habType === Temp}
				class:rad-bar={habType === Rad}
			/>
		{/if}
	</svg>
</div>
