<script lang="ts">
	import type { CommandedFleet, Fleet } from '$lib/types/Fleet';
	import hotkeys from 'hotkeys-js';
	import { createEventDispatcher, onMount } from 'svelte';
	import type { MergeFleetsEvent } from './MergeFleetsDialog.svelte';

	const dispatch = createEventDispatcher<MergeFleetsEvent>();

	export let fleet: CommandedFleet;
	export let otherFleetsHere: Fleet[];
	let selectedFleetIndexes: number[] = [];

	let fleetRefs: (HTMLLIElement | null)[] = [];

	function select(index: number) {
		if (selectedFleetIndexes.indexOf(index) == -1) {
			selectedFleetIndexes = [...selectedFleetIndexes, index];
		} else {
			selectedFleetIndexes = selectedFleetIndexes.filter((n) => n != index);
		}
	}

	function selectAll() {
		selectedFleetIndexes = Object.keys(otherFleetsHere).map((key) => parseInt(key));
	}
	function unselectAll() {
		selectedFleetIndexes = [];
	}

	function ok() {
		// TODO: otherFleetsHere[i] is sometimes undefined
		const fleetNums = selectedFleetIndexes.map((i) => otherFleetsHere[i].num);
		if (fleetNums.length > 0) {
			dispatch('merge-fleets', { fleet, fleetNums });
		}
	}

	function cancel() {
		dispatch('cancel');
	}

	onMount(() => {
		const originalScope = hotkeys.getScope();
		const scope = 'cargoTransfer';
		hotkeys('Esc', scope, cancel);
		hotkeys('Enter', scope, ok);
		hotkeys.setScope(scope);

		return () => {
			hotkeys.unbind('Esc', scope, cancel);
			hotkeys.unbind('Enter', scope, ok);
			hotkeys.deleteScope(scope);
			hotkeys.setScope(originalScope);
		};
	});
</script>

<div class="flex flex-row justify-center px-1 w-full">
	<div class="flex flex-col grow">
		<div class="text-xl font-semibold w-full text-center">Fleets to Merge</div>
		<div class="border border-secondary bg-base-300 min-w-fit max-h-[26rem] h-full overflow-y-auto">
			<ul class="w-full p-1">
				{#each otherFleetsHere as otherFleet, index}
					{#if otherFleet.num != fleet.num}
						<li
							bind:this={fleetRefs[index]}
							class="pl-1"
							class:bg-primary-focus={selectedFleetIndexes.indexOf(index) != -1}
						>
							<button class="w-full text-left" type="button" on:click={() => select(index)}>
								{otherFleet.name}
							</button>
						</li>
					{/if}
				{/each}
			</ul>
		</div>
	</div>
	<div class="flex flex-col mt-7 ml-2 gap-2">
		<button
			on:click|preventDefault={ok}
			type="submit"
			disabled={selectedFleetIndexes.length == 0}
			class="btn btn-sm normal-case btn-primary">OK</button
		>
		<button on:click={cancel} class="btn btn-outline btn-sm normal-case btn-secondary"
			>Cancel</button
		>
		<button
			type="button"
			on:click={selectAll}
			class="btn btn-outline btn-sm normal-case btn-secondary">Select All</button
		>
		<button
			type="button"
			on:click={unselectAll}
			class="btn btn-outline btn-sm normal-case btn-secondary">Unselect All</button
		>
	</div>
</div>
