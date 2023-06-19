<script lang="ts">
	import { selectedMapObject } from '$lib/services/Context';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType } from '$lib/types/MapObject';
	import type { Planet } from '$lib/types/Planet';
	import type { Player } from '$lib/types/Player';
	import { PaperAirplane } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import FleetSummary from './FleetSummary.svelte';
	import PlanetSummary from './PlanetSummary.svelte';
	import UnknownSummary from './UnknownSummary.svelte';
	import type { ShipDesign } from '$lib/types/ShipDesign';

	export let player: Player;
	export let designs: ShipDesign[];

	let selectedPlanet: Planet | undefined;
	let selectedFleet: Fleet | undefined;
	$: {
		selectedPlanet =
			$selectedMapObject?.type == MapObjectType.Planet ? ($selectedMapObject as Planet) : undefined;
		selectedFleet =
			$selectedMapObject?.type == MapObjectType.Fleet ? ($selectedMapObject as Fleet) : undefined;
	}
</script>

<div class="card bg-base-200 shadow-xl rounded-sm border-2 border-base-300">
	<div class="card-body p-2 gap-0">
		<div class="flex flex-row items-center">
			<div class="flex-1 text-center text-lg font-semibold text-secondary">
				{$selectedMapObject?.name ?? ''}
			</div>
			<Icon src={PaperAirplane} size="16" class="hover:stroke-accent" />
		</div>
		{#if selectedPlanet}
			<PlanetSummary planet={selectedPlanet} {player} />
		{:else if selectedFleet}
			<FleetSummary fleet={selectedFleet} {designs} />
		{:else}
			<UnknownSummary />
		{/if}
	</div>
</div>
