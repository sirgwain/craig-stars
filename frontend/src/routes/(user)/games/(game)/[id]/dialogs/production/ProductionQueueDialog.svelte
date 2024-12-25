<script lang="ts">
	import type { OnCancel, OnOk } from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import type { CommandedPlanet } from '$lib/types/Planet';
	import ProductionQueue from './ProductionQueue.svelte';

	const { commandedPlanet } = getGameContext();

	type Props = {
		show?: boolean;
		onOk: OnOk<CommandedPlanet>;
		onCancel: OnCancel;
		onNext: () => Promise<void>;
		onPrev: () => Promise<void>;
	};
	let { show = false, onOk, onCancel, onNext, onPrev }: Props = $props();
</script>

<div class="modal" class:modal-open={show}>
	<div
		class="modal-box max-w-full max-h-max h-full w-full lg:max-w-[40rem] lg:max-h-[48rem] p-0 md:p-[1.25rem] pb-5"
	>
		{#if $commandedPlanet && show}
			<ProductionQueue planet={$commandedPlanet} {onOk} {onCancel} {onNext} {onPrev} />
		{/if}
	</div>
</div>
