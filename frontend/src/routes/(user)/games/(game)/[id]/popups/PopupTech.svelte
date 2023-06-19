<script>
	import TechSummary from '$lib/components/tech/TechSummary.svelte';
	import { popupTech, popupTechLocation } from '$lib/services/Context';
</script>

{#if $popupTech}
	<div
		class={`absolute z-50 w-full h-full sm:w-[400px] sm:h-[400px]`}
		style={`left: ${$popupTechLocation.x}px; top: ${$popupTechLocation.y - 250}px;`}
	>
		<div class="bg-base-300 rounded-sm p-2 shadow-md w-full h-full">
			<TechSummary tech={$popupTech} />
		</div>
	</div>
{/if}
