<script lang="ts">
	import { preventDefault } from 'svelte/legacy';

	import { clamp } from '$lib/services/Math';
	import { createEventDispatcher } from 'svelte';
	import type { ValueChangedEvent } from '$lib/ValueChangedEvent';

	const dispatch = createEventDispatcher<ValueChangedEvent>();

	interface Props {
		value?: number;
		capacity?: number;
		min?: number;
		max?: any;
		editable?: boolean;
	}

	let {
		value = $bindable(0),
		capacity = 0,
		min = 0,
		max = capacity,
		editable = false
	}: Props = $props();

	let percent = $derived(capacity > 0 ? clamp((value / capacity) * 100, 0, 100) : 0);

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
	class="border border-secondary w-full h-[1rem] text-[0rem] relative select-none bg-gauge"
	class:cursor-pointer={editable}
	onpointerdown={preventDefault(
		(e) =>
			editable &&
			onPointerDown(
				(e.clientX - e.currentTarget.getBoundingClientRect().left) /
					e.currentTarget.getBoundingClientRect().width
			)
	)}
	onpointerup={preventDefault(
		(e) =>
			editable &&
			onPointerUp(
				(e.clientX - e.currentTarget.getBoundingClientRect().left) /
					e.currentTarget.getBoundingClientRect().width
			)
	)}
	onpointermove={preventDefault(
		(e) =>
			editable &&
			onPointerMove(
				(e.clientX - e.currentTarget.getBoundingClientRect().left) /
					e.currentTarget.getBoundingClientRect().width
			)
	)}
>
	<div class="font-extrabold text-sm text-center align-middle w-full absolute text-white">
		{value} of {capacity}mg
	</div>
	<div style={`width: ${percent.toFixed()}%`} class="fuel-bar h-full"></div>
</div>
