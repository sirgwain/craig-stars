<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let value = 0;
	export let dangerSpeed = 11; // no danger speed unless doing packet warp bars
	export let warnSpeed = 10;
	export let stargateSpeed = 11;
	export let defaultColor = 'warp-bar';
	export let warnColor = 'warp-warn-bar';
	export let dangerColor = 'warp-danger-bar';
	export let stargateColor = 'warp-stargate-bar';

	// set to 11 if the planet has a stargate
	export let capacity = 10;

	let percent = 0;
	let color = defaultColor;

	let mousedown = false;

	$: percent = capacity > 0 ? (value / capacity) * 100 : 0;

	$: {
		color = defaultColor;

		if (value >= dangerSpeed) {
			color = dangerColor;
		} else if (value >= warnSpeed) {
			color = warnColor;
		} else if (value >= stargateSpeed) {
			color = stargateColor;
		}
	}

	const dispatch = createEventDispatcher();

	const onMouseDown = (x: number) => {
		mousedown = true;
		updateValue(x);
	};

	const onMouseUp = (x: number) => {
		mousedown = false;
	};

	const onMouseMove = (x: number) => {
		if (mousedown) {
			updateValue(x);
		}
	};

	const updateValue = (x: number) => {
		const newValue = Math.round(x * capacity);
		if (newValue != value) {
			value = newValue;
			dispatch('valuechanged', value);
		}
	};
</script>

<div
	class="border border-secondary w-full h-[1rem] text-[0rem] relative cursor-pointer select-none"
	on:mousedown={(e) =>
		onMouseDown(
			(e.clientX - e.currentTarget.getBoundingClientRect().left) /
				e.currentTarget.getBoundingClientRect().width
		)}
	on:mouseup={(e) =>
		onMouseUp(
			(e.clientX - e.currentTarget.getBoundingClientRect().left) /
				e.currentTarget.getBoundingClientRect().width
		)}
	on:mousemove={(e) =>
		onMouseMove(
			(e.clientX - e.currentTarget.getBoundingClientRect().left) /
				e.currentTarget.getBoundingClientRect().width
		)}
>
	<div
		class="font-semibold text-sm text-center align-middle text-secondary w-full bg-blend-difference absolute"
	>
		Warp {value}
	</div>
	<div style={`width: ${percent.toFixed()}%`} class="{color} h-full" />
</div>
