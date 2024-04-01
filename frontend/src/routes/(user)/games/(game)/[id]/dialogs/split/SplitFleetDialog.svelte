<script lang="ts" context="module">
	import type { CommandedFleet, Fleet } from '$lib/types/Fleet';

	export type SplitFleetDialogEventDetails = {
		src: CommandedFleet;
		dest?: Fleet;
	};

	export type SplitFleetDialogEvent = {
		'split-fleet-dialog': SplitFleetDialogEventDetails;
	};
</script>

<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import SplitFleet, { type SplitFleetEventDetails } from './SplitFleet.svelte';

	const { split } = getGameContext();

	export let show = false;
	export let props: SplitFleetDialogEventDetails | undefined;

	const onSplitFleet = async (details: SplitFleetEventDetails) => {
		if (details) {
			await split(
				details.src,
				details.dest,
				details.srcTokens,
				details.destTokens,
				details.transferAmount
			);
		}

		// close the dialog
		show = false;
	};
</script>

<div class="modal" class:modal-open={show}>
	<div class="modal-box max-w-full max-h-max h-full w-full lg:max-w-[40rem] lg:max-h-[48rem]">
		{#if props && show}
			<SplitFleet
				src={props.src}
				dest={props.dest}
				on:split-fleet={(e) => onSplitFleet(e.detail)}
				on:cancel={() => (show = false)}
			/>
		{/if}
	</div>
</div>
