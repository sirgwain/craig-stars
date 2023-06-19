<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import type { Wormhole } from '$lib/types/Wormhole';
	import { startCase } from 'lodash-es';

	const { game, player, universe } = getGameContext();

	export let wormhole: Wormhole;

	$: destination = wormhole.destinationNum
		? $universe.getWormhole(wormhole.destinationNum)
		: undefined;
</script>

<div class="flex flex-row min-h-[11rem]">
	<div class="flex flex-col">
		<div class="avatar ">
			<div class="border-2 border-neutral mr-2 p-2 bg-black">
				<div class="mapobject-avatar wormhole bg-black" />
			</div>
		</div>
	</div>

	<div class="flex flex-col grow">
		<div class="flex flex-row">
			<div class="w-28 mr-2">Location:</div>
			<div>
				({wormhole.position.x.toFixed()}, {wormhole.position.y.toFixed()})
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-28 mr-2">Destination:</div>
			<div>
				{destination
					? `(${destination.position.x.toFixed()}, ${destination.position.y.toFixed()})`
					: 'unknown'}
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-28 mr-2">Stability:</div>
			<div>{startCase(wormhole.stability)}</div>
		</div>
	</div>
</div>
