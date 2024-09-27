<script lang="ts">
	import { emptyCargo, type Cargo } from '$lib/types/Cargo';

	export let cargo: Cargo = emptyCargo();
	export let transferAmount: Cargo = emptyCargo();
</script>

<div class="grid grid-cols-2">
	<div class="h-16 col-span-2" />

	<div class="col-span-2 mt-5" />

	<div class="text-right mr-1 h-8">Ironium</div>
	<div class="text-right">
		{(cargo?.ironium ?? 0) + (transferAmount.ironium ?? 0)}kT
	</div>
	<div class="text-right mr-1 h-8">Boranium</div>
	<div class="text-right">
		{(cargo?.boranium ?? 0) + (transferAmount.boranium ?? 0)}kT
	</div>
	<div class="text-right mr-1 h-8">Germanium</div>
	<div class="text-right">
		{(cargo?.germanium ?? 0) + (transferAmount.germanium ?? 0)}kT
	</div>

	<div class="text-right mr-1 h-8">Colonists</div>
	<div class="text-right">
		{(cargo?.colonists ?? 0) + (transferAmount.colonists ?? 0)}kT
	</div>
</div>
