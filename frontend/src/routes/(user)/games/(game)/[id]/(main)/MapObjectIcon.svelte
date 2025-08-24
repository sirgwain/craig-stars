<script lang="ts">
	import { onShipDesignTooltip } from '$lib/components/game/tooltips/ShipDesignTooltip.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { getHullIcon } from '$lib/techicon';
	import { MinefieldType, type Planet } from '$lib/types/cs-proto';
	import { getUnderlyingMapObject, type MapObjectLike } from '$lib/types/MapObject';
	import { QuestionMarkCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	const { universe } = getGameContext();

	type Props = {
		mapObject: MapObjectLike | undefined;
	};

	let { mapObject }: Props = $props();

	let { planet, fleet, wormhole, minefield, mysteryTrader, salvage, mineralPacket } = $derived(
		getUnderlyingMapObject(mapObject)
	);

	let design = $derived.by(() => {
		if (fleet?.tokens && fleet.tokens.length > 0) {
			const designNum = fleet.tokens[0].designNum;
			return $universe.getDesign(fleet.mapObject?.playerNum ?? 0, designNum);
		}
	});

	const icon = (planet: Planet) =>
		planet.mapObject?.num ? `planet-${(planet.mapObject.num - 1) % 26}` : '';
</script>

{#if planet}
	<div class="avatar">
		<div class="mapobject-avatar-wrapper">
			<div class="planet-avatar {icon(planet)} bg-black"></div>
		</div>
	</div>
{:else if fleet}
	<div class="avatar mr-2">
		<div
			class="border-2 border-neutral p-2 bg-black"
			style={`border-color: ${$universe.getPlayerColor(fleet.mapObject?.playerNum)};`}
		>
			{#if fleet.tokens.reduce((count, t) => count + t.quantity, 0) > 1}
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
{:else if mineralPacket}
	<div class="avatar">
		<div class="mapobject-avatar-wrapper">
			<div class="mapobject-avatar mineral-packet"></div>
		</div>
	</div>
{:else if salvage}
	<div class="avatar">
		<div class="mapobject-avatar-wrapper">
			<div class="mapobject-avatar salvage"></div>
		</div>
	</div>
{:else if minefield}
	<div class="avatar">
		<div class="mapobject-avatar-wrapper">
			<div
				class:standard-minefield={minefield.minefieldType === MinefieldType.STANDARD}
				class:heavy-minefield={minefield.minefieldType === MinefieldType.HEAVY}
				class:speed-bump-minefield={minefield.minefieldType === MinefieldType.SPEED_BUMP}
				class="mapobject-avatar"
			></div>
		</div>
	</div>
{:else if wormhole}
	<div class="avatar">
		<div class="mapobject-avatar-wrapper">
			<div class="mapobject-avatar wormhole"></div>
		</div>
	</div>
{:else if mysteryTrader}
	<div class="avatar">
		<div class="mapobject-avatar-wrapper">
			<div class="mapobject-avatar mystery-trader"></div>
		</div>
	</div>
{:else}
	<div class="avatar">
		<Icon src={QuestionMarkCircle} size="64" class="hover:stroke-accent" />
	</div>
{/if}
