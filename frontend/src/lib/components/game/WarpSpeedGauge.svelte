<script lang="ts">
	import { getXFromPointerEvent } from '$lib/services/Events';

	import { clamp } from '$lib/services/Math';

	type Props = {
		value?: number | undefined;
		min?: number;
		max?: number;
		dangerSpeed?: number;
		warnSpeed?: number;
		stargateSpeed?: number;
		defaultColor?: string;
		warnColor?: string;
		dangerColor?: string;
		stargateColor?: string;
		packetColor?: string;
		useStargate?: boolean;
		isPacket?: boolean;
		warp0Text?: string;
		onValueDragged?: (value: number) => void;
		onValueChanged?: (value: number) => void;
	};

	let {
		value = $bindable(0),
		min = 0,
		max = 10,
		dangerSpeed = 11,
		warnSpeed = 10,
		stargateSpeed = 11,
		defaultColor = 'warp-bar',
		warnColor = 'warp-warn-bar',
		dangerColor = 'warp-danger-bar',
		stargateColor = 'warp-stargate-bar',
		packetColor = 'warp-packet-bar',
		useStargate = false,
		isPacket = false,
		warp0Text = 'Warp 0',
		onValueDragged: onValueDragged,
		onValueChanged: onValueChanged
	}: Props = $props();

	let percent = $derived(max > 0 ? (value / max) * 100 : 0);
	let color = $derived.by(() => {
		let color = defaultColor;

		if (useStargate && value >= stargateSpeed) {
			color = stargateColor;
		} else if (isPacket && value < warnSpeed) {
			color = packetColor;
		} else if (value >= dangerSpeed) {
			color = dangerColor;
		} else if (value >= warnSpeed) {
			color = warnColor;
		}
		return color;
	});

	let pointerDown = false;
	let touchStarted = false;

	let ref: HTMLDivElement | undefined = $state();

	function onPointerDown(e: PointerEvent) {
		if (touchStarted) {
			return;
		}
		pointerDown = true;
		updateValue(getXFromPointerEvent(e, ref));
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
		onValueChanged?.(value);
	}

	const onPointerMove = (e: PointerEvent) => {
		if (pointerDown) {
			updateValue(getXFromPointerEvent(e, ref));
		}
	};

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
			onValueDragged?.(value);
		}
	};
</script>

<div
	bind:this={ref}
	class="border border-secondary w-full h-[1rem] text-[0rem] relative cursor-pointer select-none"
	onpointerdown={onPointerDown}
	ontouchstart={onTouchStart}
	ontouchmove={onTouchMove}
	ontouchend={onTouchEnd}
>
	<div
		class="font-semibold text-sm text-center align-middle text-secondary w-full bg-blend-difference absolute"
	>
		{#if useStargate && value === stargateSpeed}
			Use Stargate
		{:else if value === 0}
			{warp0Text}
		{:else}
			Warp {value}
		{/if}
	</div>
	<div style={`width: ${percent.toFixed()}%`} class="{color} h-full"></div>
</div>
