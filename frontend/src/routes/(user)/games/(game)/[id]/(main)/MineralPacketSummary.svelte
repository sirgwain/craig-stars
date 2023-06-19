<script lang="ts">
	import type { FullGame } from '$lib/services/FullGame';
	import { ownedBy } from '$lib/types/MapObject';
	import type { MineralPacket } from '$lib/types/MineralPacket';
	import type { Player } from '$lib/types/Player';

	export let game: FullGame;
	export let mineralPacket: MineralPacket;
	export let player: Player;
</script>

<div class="flex flex-row min-h-[11rem]">
	<div class="flex flex-col">
		<div class="avatar ">
			<div class="border-2 border-neutral mr-2 p-2 bg-black">
				<div class="mapobject-avatar mineral-packet bg-black" />
			</div>
		</div>
		<div class="text-center">{game.getPlayerName(mineralPacket.playerNum)}</div>
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
				{mineralPacket.warpFactor}
			</div>
		</div>
		{#if ownedBy(mineralPacket, player.num)}
			<div class="flex flex-row">
				<div class="w-28 mr-2">Destination:</div>
				<div>
					{game.getPlanet(mineralPacket.targetPlanetNum ?? 0)?.name ?? 'Unknown'}
				</div>
			</div>
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
		{/if}
	</div>
</div>
