<script lang="ts">
	import FilterFleets from '#lib/components/icons/FilterFleets.svelte';
	import { getGameContext } from '#lib/services/GameContext.js';
	import FleetFilterMenu from './FleetFilterMenu.svelte';

	const { settings, universe } = getGameContext();

	let items = $derived(
		$universe
			.getMyDesigns()
			.filter((d) => !d.spec?.starbase)
			.map((d) => ({ value: d.num, label: d.name }))
	);
</script>

<div class="flex flex-row">
	<div class="h-10 w-10">
		<a
			href="#filter-my-designs"
			title="Design Filter"
			class:btn-primary={$settings.filterMyDesigns}
			class:btn-ghost={!$settings.filterMyDesigns}
			class="btn btn-ghost fill-current btn-xs w-full h-full p-0"
			onclick={(e) => {
				e.preventDefault();
				$settings.filterMyDesigns = !$settings.filterMyDesigns;
			}}><FilterFleets class="w-6 h-6" /></a
		>
	</div>
	<div class="h-10 w-3">
		<FleetFilterMenu
			title="Design Filter Options"
			{items}
			bind:selected={$settings.filterDesigns}
			bind:enabled={$settings.filterMyDesigns}
		/>
	</div>
</div>
