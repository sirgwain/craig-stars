<script lang="ts">
	import type { VictoryCondition } from '$lib/types/Game';

	export let value: number;
	export let unit: string | undefined = undefined;
	export let min = 0;
	export let max = 100;
	export let step = 1;
	export let required = true;
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
