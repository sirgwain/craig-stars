<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { SalvageIntel } from '$lib/types/cs';

	const { universe } = getGameContext();

	type Props = {
		salvage: SalvageIntel;
	};

	let { salvage }: Props = $props();
</script>

<div class="flex flex-row md:min-h-[11rem]">
	<div class="flex flex-col">
		<div class="avatar">
			<div class="mapobject-avatar-wrapper">
				<div class="mapobject-avatar salvage"></div>
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
