<script lang="ts">
	import { clamp } from '$lib/services/Math';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	export let value = 0;
	export let capacity = 0;
	export let min = 0;
	export let max = capacity;
	export let editable = false;

	$: percent = capacity > 0 ? (value / capacity) * 100 : 0;

	let pointerdown = false;

	const onPointerDown = (x: number) => {
		pointerdown = true;
		updateValue(x);
	};

	const onPointerUp = (x: number) => {
		pointerdown = false;
		dispatch('valuechanged', value);
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
		}
	};
</script>

<div
	class="border border-secondary w-full h-[1rem] text-[0rem] relative select-none"
	class:cursor-pointer={editable}
	on:pointerdown|preventDefault={(e) =>
		editable &&
		onPointerDown(
			(e.clientX - e.currentTarget.getBoundingClientRect().left) /
				e.currentTarget.getBoundingClientRect().width
		)}
	on:pointerup|preventDefault={(e) =>
		editable &&
		onPointerUp(
			(e.clientX - e.currentTarget.getBoundingClientRect().left) /
				e.currentTarget.getBoundingClientRect().width
		)}
	on:pointermove|preventDefault={(e) =>
		editable &&
		onPointerMove(
			(e.clientX - e.currentTarget.getBoundingClientRect().left) /
				e.currentTarget.getBoundingClientRect().width
		)}
>
	<div
		class="font-semibold text-sm text-center align-middle text-secondary w-full bg-blend-difference absolute"
	>
		{value} of {capacity}mg
	</div>
	<div style={`width: ${percent.toFixed()}%`} class="fuel-bar h-full" />
</div>
