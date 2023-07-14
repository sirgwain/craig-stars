
<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { commandedPlanet } from '$lib/services/Stores';
	import ProductionQueue from './ProductionQueue.svelte';

	const { game, player, universe } = getGameContext();

	export let show = false;

	async function onOk() {
		if ($commandedPlanet) {
			await $game.updatePlanetOrders($commandedPlanet);
		}
		show = false;
	}
</script>

<div class="modal" class:modal-open={show}>
	<div
		class="modal-box max-w-full max-h-max h-full w-full lg:max-w-[40rem] lg:max-h-[48rem] p-0 md:p-[1.25rem]"
	>
		{#if $commandedPlanet}
			<ProductionQueue
				planet={$commandedPlanet}
				on:ok={onOk}
				on:cancel={() => (show = false)}
			/>
		{/if}
	</div>
</div>
