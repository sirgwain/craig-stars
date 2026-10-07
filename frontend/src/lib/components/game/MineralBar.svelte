<script lang="ts">
	import { clamp } from '#lib/services/Math.js';

	type Props = {
		value: number;
		capacity: number;
		color?: string;
		unit?: string;
		readonly?: boolean;
		onValueChanged?: (value: number) => number | undefined;
	};
	import { getXFromPointerEvent, sliderValueForKey } from '#lib/services/Events.js';

	let {
		value,
		capacity,
		color = 'ironium-bar',
		unit = 'kT',
		readonly = false,
		onValueChanged
	}: Props = $props();

	const min = 0;
	let percent = $derived(capacity > 0 ? (value / capacity) * 100 : 0);

	let pointerDown = false;
	let touchStarted = false;
	let ref: HTMLDivElement | undefined = $state();

	function onPointerDown(e: PointerEvent) {
		if (readonly) {
			return;
		}
		if (touchStarted) {
			return;
		}
		pointerDown = true;
		updateValue(getXFromPointerEvent(e, ref));
		window.addEventListener('pointerup', onPointerUp);
		window.addEventListener('pointermove', onPointerMove);
		document.body.classList.add('select-none', 'touch-none');
	}

	function onPointerUp() {
		window.removeEventListener('pointerup', onPointerUp);
		window.removeEventListener('pointermove', onPointerMove);
		document.body.classList.remove('select-none', 'touch-none');
		pointerDown = false;
	}

	function onPointerMove(e: PointerEvent) {
		if (pointerDown) {
			updateValue(getXFromPointerEvent(e, ref));
		}
	}

	function getXFromTouchEvent(e: TouchEvent): number {
		if (!ref) {
			return 0;
		}
		return (
			(e.targetTouches[0].clientX - ref.getBoundingClientRect().left) /
			ref.getBoundingClientRect().width
		);
	}

	function onTouchStart(e: TouchEvent) {
		if (readonly) {
			return;
		}
		if (e.cancelable) {
			e.preventDefault();
		}
		touchStarted = true;
		pointerDown = false;
		onPointerUp();
		updateValue(getXFromTouchEvent(e));
		document.body.classList.add('select-none', 'touch-none');
		document.body.classList.remove('touch-manipulation');
	}

	function onTouchEnd() {
		document.body.classList.remove('select-none', 'touch-none');
		document.body.classList.add('touch-manipulation');
		touchStarted = false;
	}

	function onTouchMove(e: TouchEvent) {
		if (touchStarted) {
			if (e.cancelable) {
				e.preventDefault();
			}
			updateValue(getXFromTouchEvent(e));
		}
	}

	function updateValue(x: number) {
		let newValue = clamp(Math.round(x * capacity), min, capacity);
		if (newValue != value) {
			value = onValueChanged?.(newValue) ?? newValue;
		}
	}

	function onKeyDown(e: KeyboardEvent) {
		if (readonly) return;
		const next = sliderValueForKey(e.key, value, min, capacity);
		if (next === undefined) return;
		e.preventDefault();
		e.stopPropagation();
		value = onValueChanged?.(next) ?? next;
	}
</script>

<div
	role="slider"
	aria-readonly={readonly}
	aria-label={`${color.replace('-bar', '')} amount`}
	aria-valuemin={min}
	aria-valuemax={capacity}
	aria-valuenow={value}
	tabindex={readonly ? undefined : 0}
	onkeydown={onKeyDown}
	bind:this={ref}
	class="border border-secondary w-full h-4 text-[0rem] relative bg-gauge select-none"
	class:cursor-pointer={!readonly}
	onpointerdown={onPointerDown}
	ontouchstart={onTouchStart}
	ontouchmove={onTouchMove}
	ontouchend={onTouchEnd}
>
	<div
		class="font-semibold text-sm text-center align-middle text-white mix-blend-difference w-full bg-blend-difference absolute"
	>
		{value} of {capacity}{unit}
	</div>

	<div style={`width: ${percent.toFixed()}%`} class="{color} h-full"></div>
</div>
