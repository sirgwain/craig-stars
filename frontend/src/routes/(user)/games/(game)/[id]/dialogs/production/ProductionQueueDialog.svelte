<script lang="ts" context="module">
	import type { CommandedPlanet } from '$lib/types/Planet';

	export type ProductionQueueDialogEvent = {
		'change-production': CommandedPlanet;
	};
</script>

<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import ProductionQueue from './ProductionQueue.svelte';

	const { commandedPlanet, nextMapObject, previousMapObject, updatePlanetOrders } = getGameContext();

	export let show = false;

	async function onNext() {
		if ($commandedPlanet) {
			await updatePlanetOrders($commandedPlanet);
		}

		nextMapObject();
	}

	async function onPrev() {
		if ($commandedPlanet) {
			await updatePlanetOrders($commandedPlanet);
		}

		previousMapObject();
	}

	async function onOk() {
		if ($commandedPlanet) {
			await updatePlanetOrders($commandedPlanet);
		}
		show = false;
	}
</script>

<div class="modal" class:modal-open={show}>
	<div
		class="modal-box max-w-full max-h-max h-full w-full lg:max-w-[40rem] lg:max-h-[48rem] p-0 md:p-[1.25rem] pb-5"
	>
		{#if $commandedPlanet && show}
			<ProductionQueue
				planet={$commandedPlanet}
				on:ok={onOk}
				on:next={onNext}
				on:prev={onPrev}
				on:cancel={() => (show = false)}
			/>
		{/if}
	</div>
</div>
