<script lang="ts">
	import FilterFleets from '#lib/components/icons/FilterFleets.svelte';
	import { getGameContext } from '#lib/services/GameContext.js';
	import { ShipClass } from '#lib/types/cs-proto.js';
	import { enumToString } from '#lib/types/Enums.js';
	import { ShipClasses } from '#lib/types/Filter.js';
	import FleetFilterMenu from './FleetFilterMenu.svelte';

	const { settings } = getGameContext();

	const items = ShipClasses.map((c) => ({ value: c, label: enumToString(ShipClass, c) }));
</script>

<div class="flex flex-row">
	<div class="h-10 w-10">
		<a
			href="#filter-ally-designs"
			title="Ally Ship Class Filter"
			class:btn-primary={$settings.filterAllyDesigns}
			class:btn-ghost={!$settings.filterAllyDesigns}
			class="btn btn-ghost fill-current btn-xs w-full h-full p-0"
			onclick={(e) => {
				e.preventDefault();
				$settings.filterAllyDesigns = !$settings.filterAllyDesigns;
			}}><FilterFleets primary="fill-fleet-ally" class="w-6 h-6" /></a
		>
	</div>
	<div class="h-10 w-3">
		<FleetFilterMenu
			title="Ally Ship Class Filter Options"
			{items}
			bind:selected={$settings.filterAllyShipClasses}
			bind:enabled={$settings.filterAllyDesigns}
		/>
	</div>
</div>
