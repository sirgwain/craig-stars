<script lang="ts">
	import type { MysteryTrader } from '$lib/types/cs-proto';

	type Props = {
		mysteryTrader: MysteryTrader;
	};

	let { mysteryTrader }: Props = $props();
</script>

<div class="flex flex-row md:min-h-[11rem]">
	<div class="flex flex-col">
		<div class="avatar">
			<div class="mapbject-avatar-wrapper">
				<div class="mapobject-avatar mystery-trader"></div>
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
