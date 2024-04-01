<script lang="ts" context="module">
	import type { CommandedFleet, Fleet } from '$lib/types/Fleet';
	export type MergeFleetsDialogEventDetails = {
		fleet: CommandedFleet;
		otherFleetsHere: Fleet[];
	};

	export type MergeFleetsDialogEvent = {
		'merge-fleets-dialog': MergeFleetsDialogEventDetails;
	};

	export type MergeFleetsEventDetails = {
		fleet: CommandedFleet;
		fleetNums: number[];
	};
	export type MergeFleetsEvent = {
		'merge-fleets': MergeFleetsEventDetails;
		cancel: void;
	};
</script>

<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import MergeFleets from './MergeFleets.svelte';

	const { merge } = getGameContext();

	export let show = false;
	export let props: MergeFleetsDialogEventDetails | undefined;

	const onOk = async (props: MergeFleetsEventDetails) => {
		if (props) {
			await merge(props.fleet, props.fleetNums);
		}

		// close the dialog
		show = false;
	};
</script>

<div class="modal" class:modal-open={show}>
	<div class="modal-box max-w-full max-h-max h-full w-full md:max-w-[32rem] md:max-h-[32rem]">
		{#if props && show}
			<MergeFleets
				fleet={props.fleet}
				otherFleetsHere={props.otherFleetsHere}
				on:merge-fleets={(e) => onOk(e.detail)}
				on:cancel={() => (show = false)}
			/>
		{/if}
	</div>
</div>
