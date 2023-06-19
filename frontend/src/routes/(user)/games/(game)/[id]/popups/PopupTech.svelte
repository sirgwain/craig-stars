<script>
	import TechSummary from '$lib/components/tech/TechSummary.svelte';
	import { popupTech, popupTechLocation } from '$lib/services/Context';

	$: x = Math.min($popupTechLocation.x, window.innerWidth - 400);
	$: y = Math.min($popupTechLocation.y, 400);

</script>

{#if $popupTech}
	<div
		class={`absolute z-50 w-full sm:w-[400px] h-[400px]`}
		style={`left: ${x}px; top: ${y - 250}px;`}
	>
		<div class="bg-base-300 rounded-sm p-2 shadow-md w-full h-full">
			<TechSummary tech={$popupTech} />
		</div>
	</div>
{/if}
