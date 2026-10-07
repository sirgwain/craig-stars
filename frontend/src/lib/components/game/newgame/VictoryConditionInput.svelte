<script lang="ts">
	type Props = {
		value: number;
		unit?: string | undefined;
		min?: number;
		max?: number;
		step?: number;
		required?: boolean;
	};

	let {
		value = $bindable(),
		unit = undefined,
		min = 0,
		max = 100,
		step = 1,
		required = true
	}: Props = $props();

	// keep the value in range and on a step boundary, i.e. years in increments of 10
	function onChange() {
		// a cleared input binds as null
		const current = Number.isFinite(value) ? value : min;
		const snapped = min + Math.round((current - min) / step) * step;
		value = Math.min(Math.max(snapped, min), max);
	}
</script>

<div class="join">
	<div class="join-item">
		<input
			class="input input-secondary input-sm pr-0 w-auto"
			type="number"
			{min}
			{max}
			{step}
			{required}
			bind:value
			onchange={onChange}
		/>
	</div>
	{#if unit}
		<div
			class="join-item border-l border-secondary bg-base-300 rounded-r-md px-2 flex items-center"
		>
			{unit}
		</div>
	{/if}
</div>
