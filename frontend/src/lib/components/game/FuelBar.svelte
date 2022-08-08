<script lang="ts">
	export let value = 0;
	export let capacity = 0;

	let percent = 0;

	$: percent = capacity > 0 ? (value / capacity) * 100 : 0;
</script>

<div class="border border-secondary w-full h-[1rem] text-[0rem] relative">
	<div
		class="font-semibold text-sm text-center align-middle text-secondary w-full bg-blend-difference absolute"
	>
		{value} of {capacity}kT
	</div>
	<div style={`width: ${percent.toFixed()}%`} class="fuel-bar h-full" />
</div>
