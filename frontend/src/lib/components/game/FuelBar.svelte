<script lang="ts">
	import {
		getXFromPointerEvent,
		onPointerKeyDown,
		onPointerKeyUp,
		sliderValueForKey
	} from '#lib/services/Events.js';
	import { clamp } from '#lib/services/Math.js';

	type Props = {
		value: number;
		capacity: number;
		editable?: boolean;
		valuechanged?: (value: number) => void;
		// a bar that isn't editable can be pressed instead, to open the cargo transfer dialog
		onPress?: (e: PointerEvent) => void;
	};

	let { value = $bindable(0), capacity, editable = false, valuechanged, onPress }: Props = $props();

	let pressable = $derived(!editable && !!onPress);

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

	function onKeyDown(e: KeyboardEvent) {
		if (!editable) return;
		const next = sliderValueForKey(e.key, value, 0, capacity);
		if (next === undefined) return;
		e.preventDefault();
		e.stopPropagation();
		value = next;
		valuechanged?.(value);
	}
</script>

{#snippet bar()}
	<div class="font-extrabold text-sm text-center align-middle w-full absolute text-white">
		{value} of {capacity}mg
	</div>
	<div style={`width: ${percent.toFixed()}%`} class="fuel-bar h-full"></div>
{/snippet}

{#if pressable}
	<div
		role="button"
		aria-label="Transfer fuel"
		tabindex="0"
		onkeydown={onPointerKeyDown}
		onkeyup={onPointerKeyUp}
		onpointerdown={(e) => onPress?.(e)}
		class="border border-secondary w-full h-4 text-[0rem] relative select-none bg-gauge cursor-pointer"
	>
		{@render bar()}
	</div>
{:else}
	<div
		role="slider"
		aria-readonly={!editable}
		aria-label="Fuel amount"
		aria-valuemin={0}
		aria-valuemax={capacity}
		aria-valuenow={value}
		tabindex={editable ? 0 : undefined}
		onkeydown={onKeyDown}
		class="border border-secondary w-full h-4 text-[0rem] relative select-none bg-gauge"
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
		{@render bar()}
	</div>
{/if}
