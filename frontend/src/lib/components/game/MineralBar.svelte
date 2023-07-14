<script lang="ts">
	import { clamp } from '$lib/services/Math';
	import { createEventDispatcher } from 'svelte';
	import type { ValueChangedEvent } from '$lib/ValueChangedEvent';

	export let value = 0;
	export let capacity = 0;
	export let min = 0;
	export let max = capacity;
	export let color = 'ironium-bar';

	const dispatch = createEventDispatcher<ValueChangedEvent>();

	$: percent = capacity > 0 ? (value / capacity) * 100 : 0;

	let pointerdown = false;
	let ref: HTMLDivElement;

	const getXFromPointerEvent = (e: PointerEvent) =>
		(e.clientX - ref.getBoundingClientRect().left) / ref.getBoundingClientRect()?.width;

	const onPointerDown = (x: number) => {
		pointerdown = true;
		updateValue(x);
		window.addEventListener('pointerup', onPointerUp);
		window.addEventListener('pointermove', onPointerMove);
		document.body.classList.remove('select-none', 'touch-none');
	};

	function onPointerUp(e: PointerEvent) {
		e.preventDefault();
		window.removeEventListener('pointerup', onPointerUp);
		window.removeEventListener('pointermove', onPointerMove);
		document.body.classList.add('select-none', 'touch-none');
		pointerdown = false;
	}

	const onPointerMove = (e: PointerEvent) => {
		if (pointerdown) {
			updateValue(getXFromPointerEvent(e));
		}
	};

	const updateValue = (x: number) => {
		let newValue = clamp(Math.round(x * capacity), min, max);
		if (newValue != value) {
			value = newValue;
			dispatch('valuechanged', value);
		}
	};
</script>

<div
	bind:this={ref}
	class="border border-secondary w-full h-[1rem] text-[0rem] relative bg-base-200 cursor-pointer select-none"
	on:pointerdown|preventDefault={(e) => onPointerDown(getXFromPointerEvent(e))}
>
	<div
		class="font-semibold text-sm text-center align-middle text-secondary w-full bg-blend-difference absolute"
	>
		{value} of {capacity}kT
	</div>

	<div style={`width: ${percent.toFixed()}%`} class="{color} h-full" />
</div>
