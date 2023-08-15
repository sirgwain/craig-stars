<script lang="ts">
	import Cost from '$lib/components/game/Cost.svelte';
	import Hull from '$lib/components/game/design/Hull.svelte';
	import TechAvatar from '$lib/components/tech/TechAvatar.svelte';
	import { techs } from '$lib/services/Stores';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import DesignStats from '../DesignStats.svelte';

	export let design: ShipDesign;

	$: hull = design && $techs.getHull(design.hull);
</script>

<div class="flex flex-row justify-between">
	<div class="border border-secondary bg-black p-1">
		<TechAvatar tech={hull} hullSetNumber={design.hullSetNumber} />
	</div>
</div>
{#if hull}
	<div class="flex flex-row justify-center">
		<Hull
			{hull}
			cargoCapacity={design.spec.cargoCapacity ?? hull.cargoCapacity}
			shipDesignSlots={design?.slots ?? []}
		/>
	</div>
{/if}
{#if 'spec' in design}
	<div class="flex flex-col">
		<div class="flex justify-between">
			<div class="ml-2">
				<div>Cost of one {design.name}</div>
				<Cost cost={design.spec.cost} />
			</div>
			<div>
				<DesignStats spec={design.spec} />
			</div>
		</div>
	</div>
{/if}
