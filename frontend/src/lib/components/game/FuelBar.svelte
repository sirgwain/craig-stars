<script lang="ts">
	import { getXFromPointerEvent } from '$lib/services/Events';
	import { clamp } from '$lib/services/Math';

	type Props = {
		value: number;
		capacity: number;
		editable?: boolean;
		valuechanged?: (value: number) => void;
	};

	let { value = $bindable(0), capacity, editable = false, valuechanged }: Props = $props();

	let percent = $derived(capacity > 0 ? clamp((value / capacity) * 100, 0, 100) : 0);

	let pointerdown = false;

	const onPointerDown = (x: number) => {
		pointerdown = true;
		updateValue(x);
	};

	const onPointerUp = () => {
		pointerdown = false;
		valuechanged?.(value);
	};

	const onPointerMove = (x: number) => {
		if (pointerdown) {
			updateValue(x);
		}
	};

	const updateValue = (x: number) => {
		const newValue = clamp(Math.round(x * capacity), 0, capacity);
		if (newValue != value) {
			value = newValue;
		}
	};
</script>

<div
	class="border border-secondary w-full h-[1rem] text-[0rem] relative select-none bg-gauge"
	class:cursor-pointer={editable}
	onpointerdown={(e) => {
		if (editable) {
			e.preventDefault();
			onPointerDown(getXFromPointerEvent(e, e.currentTarget));
		}
	}}
	onpointerup={(e) => {
		if (editable) {
			e.preventDefault();
			onPointerUp();
		}
	}}
	onpointermove={(e) => {
		if (editable) {
			e.preventDefault();
			onPointerMove(getXFromPointerEvent(e, e.currentTarget));
		}
	}}
>
	<div class="font-extrabold text-sm text-center align-middle w-full absolute text-white">
		{value} of {capacity}mg
	</div>
	<div style={`width: ${percent.toFixed()}%`} class="fuel-bar h-full"></div>
</div>
