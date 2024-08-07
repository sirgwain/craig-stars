<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { ownedBy } from '$lib/types/MapObject';
	import type { MineralPacket } from '$lib/types/MineralPacket';
	import { distance } from '$lib/types/Vector';

	const { game, player, universe } = getGameContext();

	export let mineralPacket: MineralPacket;
	$: target = $universe.getPlanet(mineralPacket.targetPlanetNum);
</script>

<div class="flex flex-row min-h-[11rem]">
	<div class="flex flex-col">
		<div class="avatar">
			<div class="border-2 border-neutral mr-2 p-2 bg-black">
				<div class="mapobject-avatar mineral-packet bg-black" />
			</div>
		</div>
		<div class="text-center">{$universe.getPlayerName(mineralPacket.playerNum)}</div>
	</div>

	<div class="flex flex-col grow">
		<div class="flex flex-row">
			<div class="w-28 mr-2">Location:</div>
			<div>
				({mineralPacket.position.x.toFixed()}, {mineralPacket.position.y.toFixed()})
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
				{$universe.getPlanet(mineralPacket.targetPlanetNum ?? 0)?.name ?? 'Unknown'}
			</div>
		</div>
		{#if target}
			<div class="flex flex-row">
				<div class="w-28 mr-2">ETA:</div>
				<div>
					{Math.ceil(
						distance(mineralPacket.position, target.position) /
							(mineralPacket.warpSpeed * mineralPacket.warpSpeed)
					)} years
				</div>
			</div>
		{/if}
		<div class="flex flex-row mt-2">
			<div class="text-ironium w-28 mr-2">Ironium</div>
			<div>{mineralPacket.cargo.ironium ?? 0}kT</div>
		</div>
		<div class="flex flex-row">
			<div class="text-boranium w-28 mr-2">Boranium</div>
			<div>{mineralPacket.cargo.boranium ?? 0}kT</div>
		</div>
		<div class="flex flex-row">
			<div class="text-germanium w-28 mr-2">Germanium</div>
			<div>{mineralPacket.cargo.germanium ?? 0}kT</div>
		</div>
	</div>
</div>
