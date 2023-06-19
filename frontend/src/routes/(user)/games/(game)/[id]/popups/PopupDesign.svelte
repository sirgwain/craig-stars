<script>
	import Design from '$lib/components/game/design/Design.svelte';
	import { popupDesign, popupDesignLocation } from '$lib/services/Context';

	$: x = Math.min($popupDesignLocation.x, window.innerWidth - 400);
	$: y = Math.min($popupDesignLocation.y, 400);

</script>

{#if $popupDesign}
	<div
		class={`absolute z-50`}
		style={`left: ${x}px; top: ${y - 250}px;`}
	>
		<div class="bg-base-300 rounded-sm p-2 shadow-md w-full h-full">
			<Design design={$popupDesign} />
		</div>
	</div>
{/if}
