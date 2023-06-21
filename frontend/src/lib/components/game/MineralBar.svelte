<script lang="ts">
	import { clamp } from '$lib/services/Math';
	import { createEventDispatcher } from 'svelte';

	export let value = 0;
	export let capacity = 0;
	export let min = 0;
	export let max = capacity;
	export let color = 'ironium-bar';

	const dispatch = createEventDispatcher();

	$: percent = capacity > 0 ? (value / capacity) * 100 : 0;

	let pointerdown = false;

	const onPointerDown = (x: number) => {
		pointerdown = true;
		updateValue(x);
	};

	const onPointerUp = (x: number) => {
		pointerdown = false;
	};

	const onPointerMove = (x: number) => {
		if (pointerdown) {
			updateValue(x);
		}
	};

	const updateValue = (x: number) => {
		const newValue = clamp(Math.round(x * capacity), min, max);
		if (newValue != value) {
			value = newValue;
			dispatch('valuechanged', value);
		}
	};
</script>

<div
	class="border border-secondary w-full h-[1rem] text-[0rem] relative bg-base-200 cursor-pointer select-none"
	on:pointerdown|preventDefault={(e) =>
		onPointerDown(
			(e.clientX - e.currentTarget.getBoundingClientRect().left) /
				e.currentTarget.getBoundingClientRect().width
		)}
	on:pointerup|preventDefault={(e) =>
		onPointerUp(
			(e.clientX - e.currentTarget.getBoundingClientRect().left) /
				e.currentTarget.getBoundingClientRect().width
		)}
	on:pointermove|preventDefault={(e) =>
		onPointerMove(
			(e.clientX - e.currentTarget.getBoundingClientRect().left) /
				e.currentTarget.getBoundingClientRect().width
		)}
>
	<div
		class="font-semibold text-sm text-center align-middle text-secondary w-full bg-blend-difference absolute"
	>
		{value} of {capacity}kT
	</div>

	<div style={`width: ${percent.toFixed()}%`} class="{color} h-full" />
</div>
