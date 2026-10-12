<script lang="ts">
	import Cost from '#lib/components/game/Cost.svelte';
	import Hull from '#lib/components/game/design/Hull.svelte';
	import TechAvatar from '#lib/components/tech/TechAvatar.svelte';
	import { techs } from '#lib/services/Stores.js';
	import { total } from '#lib/types/Cost.js';
	import type { TokenDamage } from '#lib/types/Fleet.js';
	import type { ShipDesign } from '#lib/types/cs-proto.js';
	import type { Snippet } from 'svelte';
	import DesignStats from '../DesignStats.svelte';

	type Props = {
		design: ShipDesign;
		damage?: TokenDamage;
		// shown in place of the player's notes on the design, i.e. to edit them
		description?: Snippet;
	};

	let { design, damage, description }: Props = $props();

	let hull = $derived($techs.getHull(design.hull));
</script>

<!-- the title is centered over everything, and the icon's top lines up with the notes beside it -->
<div class="grid grid-cols-[auto_1fr]">
	<h1 class="col-span-2 text-lg font-semibold min-w-0 text-center">
		{design.name}
	</h1>
	<!-- the notes have a little margin above them -->
	<div class="self-start mt-1 border border-secondary bg-black p-1">
		<TechAvatar tech={hull} hullSetNumber={design.hullSetNumber} />
	</div>
	<div class="min-w-0">
		{#if description}
			{@render description()}
		{:else if design.description}
			<p class="whitespace-pre-wrap wrap-break-word mx-2 my-1" data-type="design-description">
				{design.description}
			</p>
		{/if}
	</div>
</div>
{#if hull}
	<div class="flex flex-row justify-center">
		<Hull
			{hull}
			cargoCapacity={design.spec?.cargoCapacity ? design.spec.cargoCapacity : hull.cargoCapacity}
			shipDesignSlots={design.slots}
		/>
	</div>
{/if}
<div class="flex flex-col">
	<div class="flex flex-col sm:flex-row gap-1 justify-between">
		{#if design.spec?.cost && total(design.spec.cost)}
			<div class="mx-2">
				<div>Cost of one {design.name}</div>
				<Cost cost={design.spec.cost} />
			</div>
		{/if}
		<div class="mx-2 sm:mx-0">
			<DesignStats spec={design.spec} {damage} />
		</div>
	</div>
	{#if design.spec?.numBuilt}
		<div class="mx-2">
			{design.spec.numInstances} remaining of {design.spec.numBuilt}
			{design.originalPlayerNum ? 'received' : 'built'}
		</div>
	{/if}
</div>
