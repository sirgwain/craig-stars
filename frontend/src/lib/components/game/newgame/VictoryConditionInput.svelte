<script lang="ts">
	import type { VictoryCondition } from '$lib/types/Game';

	interface Props {
		value: number;
		unit?: string | undefined;
		min?: number;
		max?: number;
		step?: number;
		required?: boolean;
	}

	let {
		value = $bindable(),
		unit = undefined,
		min = 0,
		max = 100,
		step = 1,
		required = true
	}: Props = $props();
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
