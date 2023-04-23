<script lang="ts">
	import Cost from '$lib/components/game/Cost.svelte';
	import Hull from '$lib/components/game/design/Hull.svelte';
	import TechAvatar from '$lib/components/tech/TechAvatar.svelte';
	import { techs } from '$lib/services/Context';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import DesignStats from '../DesignStats.svelte';


	export let design: ShipDesign;

	$: hull = design && $techs.getHull(design.hull);
</script>

<div class="flex flex-row justify-between">
	<div>
		<TechAvatar tech={hull} hullSetNumber={design.hullSetNumber} />
	</div>
</div>
{#if hull}
	<div class="flex flex-row justify-center">
		<Hull {hull} shipDesignSlots={design?.slots ?? []} />
	</div>
{/if}
<div class="flex flex-col">
	<div>Cost of one {design.name}</div>
	<div class="flex justify-between">
		<div class="ml-2">
			<Cost cost={design.spec.cost} />
		</div>
		<div>
			<DesignStats {design} />
		</div>
	</div>
</div>
