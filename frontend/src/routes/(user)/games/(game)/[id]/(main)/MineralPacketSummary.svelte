<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { AnyMineralPacket } from '$lib/services/Universe';
	import { distance } from '$lib/types/Vector';

	const { universe } = getGameContext();

	type Props = {
		mineralPacket: AnyMineralPacket;
	};

	let { mineralPacket }: Props = $props();
	let target = $derived($universe.getPlanet(mineralPacket.targetPlanetNum));
</script>

<div class="flex flex-row md:min-h-[11rem]">
	<div class="flex flex-col">
		<div class="avatar">
			<div class="mapobject-avatar-wrapper">
				<div class="mapobject-avatar mineral-packet"></div>
			</div>
		</div>
		<div class="text-center">
			{$universe.getPlayerPluralName(mineralPacket.mapObject?.playerNum)}
		</div>
	</div>

	<div class="flex flex-col grow">
		<div class="flex flex-row">
			<div class="w-28 mr-2">Location:</div>
			<div>
				({mineralPacket.mapObject?.position?.x?.toFixed() ?? 0}, {mineralPacket.mapObject?.position?.y?.toFixed() ??
					0})
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-28 mr-2">Traveling at Warp:</div>
			<div>
				{mineralPacket.warpSpeed}
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-28 mr-2">Destination:</div>
			<div>
				{$universe.getPlanet(mineralPacket.targetPlanetNum)?.mapObject?.name ?? 'Unknown'}
			</div>
		</div>
		{#if target}
			<div class="flex flex-row">
				<div class="w-28 mr-2">ETA:</div>
				<div>
					{Math.ceil(
						distance(mineralPacket.mapObject?.position, target.mapObject?.position) /
							(mineralPacket.warpSpeed * mineralPacket.warpSpeed)
					)} years
				</div>
			</div>
		{/if}
		<div class="flex flex-row mt-2">
			<div class="text-ironium w-28 mr-2">Ironium</div>
			<div>{mineralPacket.cargo?.ironium ?? 0}kT</div>
		</div>
		<div class="flex flex-row">
			<div class="text-boranium w-28 mr-2">Boranium</div>
			<div>{mineralPacket.cargo?.boranium ?? 0}kT</div>
		</div>
		<div class="flex flex-row">
			<div class="text-germanium w-28 mr-2">Germanium</div>
			<div>{mineralPacket.cargo?.germanium ?? 0}kT</div>
		</div>
	</div>
</div>
