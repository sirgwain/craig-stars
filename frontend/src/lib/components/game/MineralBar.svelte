<script lang="ts">
	import { clamp } from '$lib/services/Math';

	type Props = {
		value?: number;
		capacity?: number;
		min?: number;
		max?: number;
		color?: string;
		unit?: string;
		readonly?: boolean;
		onValueChanged?: (value: number) => void;
	};
	import { getXFromPointerEvent } from '$lib/services/Events';

	let {
		value = $bindable(0),
		capacity = 0,
		min = 0,
		max = capacity,
		color = 'ironium-bar',
		unit = 'kT',
		readonly = false,
		onValueChanged: onValueChanged
	}: Props = $props();

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
			ref.getBoundingClientRect()?.width
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
		let newValue = clamp(Math.round(x * capacity), min, max);
		if (newValue != value) {
			value = newValue;
			onValueChanged?.(value);
		}
	}
</script>

<div
	bind:this={ref}
	class="border border-secondary w-full h-[1rem] text-[0rem] relative bg-gauge select-none"
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
