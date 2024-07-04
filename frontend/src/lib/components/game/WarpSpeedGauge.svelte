<script lang="ts">
	import { clamp } from '$lib/services/Math';
	import { createEventDispatcher } from 'svelte';

	export let value: number | undefined = 0;
	export let min = 0;
	export let max = 10;
	export let dangerSpeed = 11; // no danger speed unless doing packet warp bars
	export let warnSpeed = 10;
	export let stargateSpeed = 11;
	export let defaultColor = 'warp-bar';
	export let warnColor = 'warp-warn-bar';
	export let dangerColor = 'warp-danger-bar';
	export let stargateColor = 'warp-stargate-bar';
	export let useStargate = false;
	export let warp0Text = 'Warp 0';

	let percent = 0;
	let color = defaultColor;

	let pointerDown = false;
	let touchStarted = false;

	$: percent = max > 0 ? ((value ?? 0) / max) * 100 : 0;

	$: {
		color = defaultColor;

		if (useStargate && (value ?? 0) >= stargateSpeed) {
			color = stargateColor;
		} else if ((value ?? 0) >= dangerSpeed) {
			color = dangerColor;
		} else if ((value ?? 0) >= warnSpeed) {
			color = warnColor;
		}
	}

	const dispatch = createEventDispatcher();

	let ref: HTMLDivElement;

	const getXFromPointerEvent = (e: PointerEvent) =>
		(e.clientX - ref.getBoundingClientRect().left) / ref.getBoundingClientRect().width;

	function onPointerDown(e: PointerEvent) {
		if (touchStarted) {
			return;
		}
		pointerDown = true;
		updateValue(getXFromPointerEvent(e));
		window.addEventListener('pointerup', onPointerUp);
		window.addEventListener('pointermove', onPointerMove);
		document.body.classList.add('select-none', 'touch-none');
		document.body.classList.remove('touch-manipulation');
	}

	function onPointerUp() {
		window.removeEventListener('pointerup', onPointerUp);
		window.removeEventListener('pointermove', onPointerMove);
		document.body.classList.remove('select-none', 'touch-none');
		document.body.classList.add('touch-manipulation');
		pointerDown = false;
		dispatch('valuechanged', value);
	}

	const onPointerMove = (e: PointerEvent) => {
		if (pointerDown) {
			updateValue(getXFromPointerEvent(e));
		}
	};

	function getXFromTouchEvent(e: TouchEvent): number {
		return (
			(e.targetTouches[0].clientX - ref.getBoundingClientRect().left) /
			ref.getBoundingClientRect()?.width
		);
	}

	function onTouchStart(e: TouchEvent) {
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

	const updateValue = (x: number) => {
		const newValue = clamp(Math.round(x * max), min, max);
		if (newValue != value) {
			value = newValue;
			dispatch('valuedragged', value);
		}
	};
</script>

<div
	bind:this={ref}
	class="border border-secondary w-full h-[1rem] text-[0rem] relative cursor-pointer select-none"
	on:pointerdown={onPointerDown}
	on:touchstart={onTouchStart}
	on:touchmove={onTouchMove}
	on:touchend={onTouchEnd}
>
	<div
		class="font-semibold text-sm text-center align-middle text-secondary w-full bg-blend-difference absolute"
	>
		{#if useStargate && value === stargateSpeed}
			Use Stargate
		{:else if value === 0 || value == undefined}
			{warp0Text}
		{:else}
			Warp {value}
		{/if}
	</div>
	<div style={`width: ${percent.toFixed()}%`} class="{color} h-full" />
</div>
