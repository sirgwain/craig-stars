<script lang="ts">
	
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
		onvaluedragged?: (value: number) => void;
		onvaluechanged?: (value: number) => void;
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
		onvaluedragged,
		onvaluechanged
	}: Props = $props();

	let percent = $derived(max > 0 ? ((value ?? 0) / max) * 100 : 0);
	let color = $derived.by(() => {
		let color = defaultColor;

		if (useStargate && (value ?? 0) >= stargateSpeed) {
			color = stargateColor;
		} else if (isPacket && (value ?? 0) < warnSpeed) {
			color = packetColor;
		} else if ((value ?? 0) >= dangerSpeed) {
			color = dangerColor;
		} else if ((value ?? 0) >= warnSpeed) {
			color = warnColor;
		}
		return color;
	});

	let pointerDown = false;
	let touchStarted = false;

	let ref: HTMLDivElement | undefined = $state();

	function getXFromPointerEvent(e: PointerEvent): number {
		if (!ref) {
			return 0;
		}
		return (e.clientX - ref.getBoundingClientRect().left) / ref.getBoundingClientRect()?.width;
	}

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
		onvaluechanged?.(value);
	}

	const onPointerMove = (e: PointerEvent) => {
		if (pointerDown) {
			updateValue(getXFromPointerEvent(e));
		}
	};

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
			onvaluedragged?.(value);
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
		{:else if value === 0 || value == undefined}
			{warp0Text}
		{:else}
			Warp {value}
		{/if}
	</div>
	<div style={`width: ${percent.toFixed()}%`} class="{color} h-full"></div>
</div>
