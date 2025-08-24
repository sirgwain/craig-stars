<script lang="ts">
	import { emptyCargo } from '$lib/types/Cargo';
	import type { Cargo } from '$lib/types/cs-proto';

	type Props = {
		cargo?: Cargo;
		transferAmount?: Cargo;
	};

	let { cargo = emptyCargo(), transferAmount = emptyCargo() }: Props = $props();
	let colonists = $derived(cargo.colonists + transferAmount.colonists);
</script>

<div class="grid grid-cols-2">
	<div class="h-16 col-span-2"></div>

	<div class="col-span-2 mt-5"></div>

	<div class="text-right mr-1 h-8 select-none">Ironium</div>
	<div class="text-right">
		{cargo.ironium + transferAmount.ironium}kT
	</div>
	<div class="text-right mr-1 h-8 select-none">Boranium</div>
	<div class="text-right">
		{cargo.boranium + transferAmount.boranium}kT
	</div>
	<div class="text-right mr-1 h-8 select-none">Germanium</div>
	<div class="text-right">
		{cargo.germanium + transferAmount.germanium}kT
	</div>

	<div class="text-right mr-1 h-8 select-none">Colonists</div>
	<div class="text-right">
		{colonists}kT
	</div>

	{#if colonists}
		<div class="my-auto col-span-2 ml-auto pr-1">
			<span class="italic text-sm">{(colonists * 100).toLocaleString()} colonists</span>
		</div>
	{/if}
</div>
