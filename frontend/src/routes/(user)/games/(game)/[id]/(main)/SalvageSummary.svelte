<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { ownedBy } from '$lib/types/MapObject';
	import type { Salvage } from '$lib/types/Salvage';

	const { game, player, universe } = getGameContext();

	export let salvage: Salvage;
</script>

<div class="flex flex-row min-h-[11rem]">
	<div class="flex flex-col">
		<div class="avatar ">
			<div class="border-2 border-neutral mr-2 p-2 bg-black">
				<div class="mapobject-avatar salvage bg-black" />
			</div>
		</div>
		<div class="text-center">{$universe.getPlayerPluralName(salvage.playerNum)}</div>
	</div>

	<div class="flex flex-col grow">
		<div class="flex flex-row">
			<div class="w-28 mr-2">Location:</div>
			<div>
				({salvage.position.x.toFixed()}, {salvage.position.y.toFixed()})
			</div>
		</div>
		<div class="flex flex-row mt-2">
			<div class="text-ironium w-28 mr-2">Ironium</div>
			<div>{salvage.cargo.ironium ?? 0}kT</div>
		</div>
		<div class="flex flex-row">
			<div class="text-boranium w-28 mr-2">Boranium</div>
			<div>{salvage.cargo.boranium ?? 0}kT</div>
		</div>
		<div class="flex flex-row">
			<div class="text-germanium w-28 mr-2">Germanium</div>
			<div>{salvage.cargo.germanium ?? 0}kT</div>
		</div>
	</div>
</div>
