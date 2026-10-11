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
	<div class="size-10">
		<a
			href="#filter-enemy-designs"
			title="Enemy Ship Class Filter"
			class:btn-primary={$settings.filterEnemyDesigns}
			class:btn-ghost={!$settings.filterEnemyDesigns}
			class="btn btn-ghost fill-current btn-xs size-full p-0"
			onclick={(e) => {
				e.preventDefault();
				$settings.filterEnemyDesigns = !$settings.filterEnemyDesigns;
			}}><FilterFleets primary="fill-fleet-enemy" class="size-6 " /></a
		>
	</div>
	<div class="h-10 w-3">
		<FleetFilterMenu
			title="Enemy Ship Class Filter Options"
			end
			{items}
			bind:selected={$settings.filterEnemyShipClasses}
			bind:enabled={$settings.filterEnemyDesigns}
		/>
	</div>
</div>
