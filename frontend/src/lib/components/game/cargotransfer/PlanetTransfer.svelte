<script lang="ts">
	import { emptyCargo } from '#lib/types/Cargo.js';
	import type { Cargo } from '#lib/types/cs-proto.js';

	type Props = {
		cargo?: Cargo;
		transferAmount?: Cargo;
	};

	let { cargo = emptyCargo(), transferAmount = emptyCargo() }: Props = $props();
	let colonists = $derived(cargo.colonists + transferAmount.colonists);
</script>

<!-- on small screens, stack each label over its value so the rows line up with the fleet's bars -->
<div class="sm:grid sm:grid-cols-2">
	<div class="h-24 sm:h-16 col-span-2"></div>

	<div class="col-span-2 mt-10 sm:mt-5"></div>

	<div class="text-right mr-1 h-8 select-none text-ironium">Ironium</div>
	<div class="text-right h-4 leading-4 sm:h-auto sm:leading-6">
		{cargo.ironium + transferAmount.ironium}kT
	</div>
	<div class="text-right mr-1 h-8 select-none text-boranium">Boranium</div>
	<div class="text-right h-4 leading-4 sm:h-auto sm:leading-6">
		{cargo.boranium + transferAmount.boranium}kT
	</div>
	<div class="text-right mr-1 h-8 select-none text-germanium">Germanium</div>
	<div class="text-right h-4 leading-4 sm:h-auto sm:leading-6">
		{cargo.germanium + transferAmount.germanium}kT
	</div>

	<div class="text-right mr-1 h-8 select-none">Colonists</div>
	<div class="text-right h-4 leading-4 sm:h-auto sm:leading-6">
		{colonists}kT
	</div>

	{#if colonists}
		<div class="my-auto col-span-2 ml-auto pr-1">
			<span class="italic text-sm">{(colonists * 100).toLocaleString()} colonists</span>
		</div>
	{/if}
</div>
