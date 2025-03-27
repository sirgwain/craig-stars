<script lang="ts">
	import { onShipDesignTooltip } from '$lib/components/game/tooltips/ShipDesignTooltip.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import type { AnyFleet, AnyPlanet } from '$lib/services/Universe';
	import { getHullIcon } from '$lib/techicon';
	import { MapObjectTypeFleet, MapObjectTypePlanet, type MapObject } from '$lib/types/cs';
	import { QuestionMarkCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	const { universe } = getGameContext();

	type Props = {
		mapObject: MapObject | undefined;
	};

	let { mapObject }: Props = $props();

	let planet = $derived(
		mapObject?.type === MapObjectTypePlanet ? (mapObject as AnyPlanet) : undefined
	);
	let fleet = $derived(
		mapObject?.type === MapObjectTypeFleet ? (mapObject as AnyFleet) : undefined
	);
	let design = $derived.by(() => {
		if (fleet?.tokens && fleet.tokens.length > 0) {
			const designNum = fleet.tokens[0].designNum;
			return $universe.getDesign(fleet.playerNum, designNum);
		}
	});

	const icon = (planet: AnyPlanet) => (planet ? `planet-${(planet.num - 1) % 26}` : '');
</script>

{#if planet}
	<div class="avatar">
		<div class="border-2 border-neutral mr-2 p-2 bg-black">
			<div class="planet-avatar {icon(planet)} bg-black"></div>
		</div>
	</div>
{:else if fleet}
	<div class="avatar mr-2">
		<div
			class="border-2 border-neutral p-2 bg-black"
			style={`border-color: ${$universe.getPlayerColor(fleet.playerNum)};`}
		>
			{#if fleet.tokens && fleet.tokens.reduce((count, t) => count + t.quantity, 0) > 1}
				<div class="absolute -right-2 -top-1 text-xl w-6 h-6">+</div>
			{/if}

			<div class="fleet-avatar {getHullIcon(design)} bg-black">
				<button
					type="button"
					aria-label="Opens ship design tooltip"
					class="w-full h-full cursor-help"
					onpointerdown={(e) => onShipDesignTooltip(e, design)}
				></button>
			</div>
		</div>
	</div>
{:else}
	<div class="avatar">
		<Icon src={QuestionMarkCircle} size="64" class="hover:stroke-accent" />
	</div>
{/if}
