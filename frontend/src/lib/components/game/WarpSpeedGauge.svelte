<script lang="ts">
	import { clamp } from '$lib/services/Math';
	import { createEventDispatcher } from 'svelte';

	export let value = 0;
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

	let percent = 0;
	let color = defaultColor;

	let pointerdown = false;

	$: percent = max > 0 ? (value / max) * 100 : 0;

	$: {
		color = defaultColor;

		if (useStargate && value >= stargateSpeed) {
			color = stargateColor;
		} else if (value >= dangerSpeed) {
			color = dangerColor;
		} else if (value >= warnSpeed) {
			color = warnColor;
		}
	}

	const dispatch = createEventDispatcher();

	let ref: HTMLDivElement;

	const getXFromPointerEvent = (e: PointerEvent) =>
		(e.clientX - ref.getBoundingClientRect().left) / ref.getBoundingClientRect().width;

	const onPointerDown = (x: number) => {
		pointerdown = true;
		updateValue(x);
		window.addEventListener('pointerup', onPointerUp);
		window.addEventListener('pointermove', onPointerMove);
	};

	function onPointerUp(e: PointerEvent) {
		e.preventDefault();
		window.removeEventListener('pointerup', onPointerUp);
		window.removeEventListener('pointermove', onPointerMove);
		pointerdown = false;
		dispatch('valuechanged', value);
	}

	const onPointerMove = (e: PointerEvent) => {
		if (pointerdown) {
			updateValue(getXFromPointerEvent(e));
		}
	};

	const updateValue = (x: number) => {
		const newValue = clamp(Math.round(x * max), min, max);
		if (newValue != value) {
			value = newValue;
		}
	};
</script>

<div
	bind:this={ref}
	class="border border-secondary w-full h-[1rem] text-[0rem] relative cursor-pointer select-none"
	on:pointerdown|preventDefault={(e) =>
		onPointerDown(
			(e.clientX - e.currentTarget.getBoundingClientRect().left) /
				e.currentTarget.getBoundingClientRect().width
		)}
>
	<div
		class="font-semibold text-sm text-center align-middle text-secondary w-full bg-blend-difference absolute"
	>
		{#if useStargate && value === stargateSpeed}
			Use Stargate
		{:else}
			Warp {value}
		{/if}
	</div>
	<div style={`width: ${percent.toFixed()}%`} class="{color} h-full" />
</div>
