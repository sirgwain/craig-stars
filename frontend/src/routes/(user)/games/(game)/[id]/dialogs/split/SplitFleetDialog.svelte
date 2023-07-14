<script lang="ts" context="module">
	import type { CommandedFleet, Fleet, ShipToken } from '$lib/types/Fleet';
	export type SplitFleetDialogEventDetails = {
		fleet: CommandedFleet;
		target?: Fleet | undefined; // sometimes we "split" tokens into another fleet. this is also a merge
	};
	export type SplitAllEventDetails = {
		fleet: CommandedFleet;
	};
	export type SplitFleetEventDetails = {
		fleet: CommandedFleet;
		tokens: ShipToken[];
	};

	export type SplitFleetEvent = {
		'split-fleet-dialog': SplitFleetDialogEventDetails;
		'split-fleet': SplitFleetEventDetails;
		'split-all': SplitAllEventDetails;
	};
</script>

<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import SplitFleet from './SplitFleet.svelte';

	const { game, player, universe } = getGameContext();

	export let show = false;
	export let props: SplitFleetEventDetails | undefined;

	const onOk = async () => {
		if (props) {
			// TODO: do the server side split fleet
		}

		// close the dialog
		show = false;
	};
</script>

<div class="modal" class:modal-open={show}>
	<div class="modal-box max-w-full max-h-max h-full w-full lg:max-w-[40rem] lg:max-h-[48rem]">
		{#if props}
			<SplitFleet fleet={props.fleet} on:ok={onOk} on:cancel={() => (show = false)} />
		{/if}
	</div>
</div>
