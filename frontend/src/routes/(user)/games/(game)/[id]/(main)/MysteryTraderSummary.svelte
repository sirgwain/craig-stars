<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { ownedBy } from '$lib/types/MapObject';
	import type { MysteryTrader } from '$lib/types/MysteryTrader';

	const { game, player, universe } = getGameContext();

	export let mysteryTrader: MysteryTrader;
</script>

<div class="flex flex-row min-h-[11rem]">
	<div class="flex flex-col">
		<div class="avatar">
			<div class="border-2 border-neutral mr-2 p-2 bg-black">
				<div class="mapobject-avatar mystery-trader bg-black" />
			</div>
		</div>
	</div>

	<div class="flex flex-col grow gap-y-1">
		<div class="flex flex-row">
			The trader requests interested parties to send it a feet with at least {mysteryTrader.requestedBoon}kT
			of minerals on board to be absorbed into the trader. It offers technological assistance in
			return.
		</div>
		<div class="flex flex-row">
			Trader is traveling at Warp {mysteryTrader.warpSpeed}.
		</div>
	</div>
</div>
