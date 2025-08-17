<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { WormholeStability, type WormholeIntel } from '$lib/types/cs-proto';
	import { enumToString } from '$lib/types/Enums';

	const { universe } = getGameContext();

	type Props = {
		wormhole: WormholeIntel;
	};

	let { wormhole }: Props = $props();

	let destination = $derived(
		wormhole.destinationNum ? $universe.getWormhole(wormhole.destinationNum) : undefined
	);
</script>

<div class="flex flex-row md:min-h-[11rem]">
	<div class="flex flex-col">
		<div class="avatar">
			<div class="mapobject-avatar-wrapper">
				<div class="mapobject-avatar wormhole"></div>
			</div>
		</div>
	</div>

	<div class="flex flex-col grow">
		<div class="flex flex-row">
			<div class="w-28 mr-2">Location:</div>
			<div>
				({wormhole.mapObject?.position?.x.toFixed()}, {wormhole.mapObject?.position?.y.toFixed()})
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-28 mr-2">Destination:</div>
			<div>
				{destination
					? `(${destination.mapObject?.position?.x.toFixed()}, ${destination.mapObject?.position?.y.toFixed()})`
					: 'unknown'}
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-28 mr-2">Stability:</div>
			<div>{enumToString(WormholeStability, wormhole.stability)}</div>
		</div>
	</div>
</div>
