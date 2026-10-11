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
	<div class="size-10">
		<a
			href="#filter-my-designs"
			title="Design Filter"
			class:btn-primary={$settings.filterMyDesigns}
			class:btn-ghost={!$settings.filterMyDesigns}
			class="btn btn-ghost fill-current btn-xs size-full p-0"
			onclick={(e) => {
				e.preventDefault();
				$settings.filterMyDesigns = !$settings.filterMyDesigns;
			}}><FilterFleets class="size-6 " /></a
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
