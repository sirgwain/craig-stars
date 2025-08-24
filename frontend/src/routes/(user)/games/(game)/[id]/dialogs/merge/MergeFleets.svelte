<script lang="ts">
	import type { Fleet } from '$lib/types/cs-proto';
	import type { MergeFleetsEvent, OnCancel, OnOk } from '$lib/services/Events';
	import { type CommandedFleet } from '$lib/types/Fleet';
	import { getMapObjectName, key } from '$lib/types/MapObject';
	import hotkeys from 'hotkeys-js';
	import { onMount } from 'svelte';

	type Props = {
		fleet: CommandedFleet;
		otherFleetsHere: Fleet[];
		onOk?: OnOk<MergeFleetsEvent>;
		onCancel?: OnCancel;
	};

	let { fleet, otherFleetsHere, onOk, onCancel }: Props = $props();
	let selectedFleetIndexes: number[] = $state([]);

	let fleetRefs: (HTMLLIElement | null)[] = $state([]);

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
		const fleetNums = selectedFleetIndexes.map((i) => otherFleetsHere[i].mapObject?.num ?? 0);
		if (fleetNums.length > 0) {
			onOk?.({ fleet, fleetNums });
		}
	}

	function cancel() {
		onCancel?.();
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
				{#each otherFleetsHere as otherFleet, index (key(otherFleet))}
					{#if otherFleet.mapObject?.num !== fleet.mapObject.num}
						<li
							bind:this={fleetRefs[index]}
							class="pl-1"
							class:bg-primary-focus={selectedFleetIndexes.indexOf(index) != -1}
						>
							<button class="w-full text-left" type="button" onclick={() => select(index)}>
								{getMapObjectName(otherFleet)}
							</button>
						</li>
					{/if}
				{/each}
			</ul>
		</div>
	</div>
	<div class="flex flex-col mt-7 ml-2 gap-2">
		<button
			type="submit"
			onclick={(e) => {
				e.preventDefault();
				ok();
			}}
			disabled={selectedFleetIndexes.length == 0}
			class="btn btn-sm normal-case btn-primary">OK</button
		>
		<button onclick={onCancel} class="btn btn-outline btn-sm normal-case btn-secondary"
			>Cancel</button
		>
		<button
			type="button"
			onclick={selectAll}
			class="btn btn-outline btn-sm normal-case btn-secondary">Select All</button
		>
		<button
			type="button"
			onclick={unselectAll}
			class="btn btn-outline btn-sm normal-case btn-secondary">Unselect All</button
		>
	</div>
</div>
