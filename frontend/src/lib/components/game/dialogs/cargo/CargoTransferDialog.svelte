<script lang="ts">
	import { PlanetService } from '$lib/services/PlanetService';
	import {
		negativeCargo,
		emptyCargo,
		type Cargo,
		totalCargo,
		subtract,
		add
	} from '$lib/types/Cargo';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType } from '$lib/types/MapObject';
	import type { Planet } from '$lib/types/Planet';
	import hotkeys from 'hotkeys-js';
	import { createEventDispatcher } from 'svelte';
	import FleetTransfer from './FleetTransfer.svelte';
	import PlanetTransfer from './PlanetTransfer.svelte';
	import TransferButtons from './TransferButtons.svelte';
	import { getQuantityModifier } from '$lib/quantityModifier';
	import { clamp } from '$lib/services/Math';

	export let src: Fleet | undefined;
	export let dest: Fleet | Planet | undefined;

	let transferAmount: Cargo = emptyCargo();
	let cargo: Cargo = emptyCargo();
	let fuelTransferAmount: number = 0;
	let fuel: number = 0;

	const ok = async () => {
		dispatch('ok');
	};

	const cancel = () => {
		dispatch('cancel');
	};

	const transferIronium = (amount: number) => {
		const available = (src?.spec.cargoCapacity ?? 0) - totalCargo(transferAmount);
		transferAmount.ironium = clamp(
			(transferAmount.ironium ?? 0) + amount,
			-(cargo.ironium ?? 0),
			(transferAmount.ironium ?? 0) + available
		);
	};

	const transferBoranium = (amount: number) => {
		const available = (src?.spec.cargoCapacity ?? 0) - totalCargo(transferAmount);
		transferAmount.boranium = clamp(
			(transferAmount.boranium ?? 0) + amount,
			-(cargo.boranium ?? 0),
			(transferAmount.boranium ?? 0) + available
		);
	};

	const transferGermanium = (amount: number) => {
		const available = (src?.spec.cargoCapacity ?? 0) - totalCargo(transferAmount);
		transferAmount.germanium = clamp(
			(transferAmount.germanium ?? 0) + amount,
			-(cargo.germanium ?? 0),
			(transferAmount.germanium ?? 0) + available
		);
	};

	const transferColonists = (amount: number) => {
		const available = (src?.spec.cargoCapacity ?? 0) - totalCargo(transferAmount);
		transferAmount.colonists = clamp(
			(transferAmount.colonists ?? 0) + amount,
			-(cargo.colonists ?? 0),
			(transferAmount.colonists ?? 0) + available
		);
	};

	const dispatch = createEventDispatcher();

	hotkeys('Esc', () => cancel());
	hotkeys('Enter', () => {
		ok();
	});

	const planetService = new PlanetService();

	$: src && (cargo = src.cargo);
	$: src && (fuel = src.fuel);

	function minus(cargo: Cargo, transferAmount: Cargo): Cargo {
		throw new Error('Function not implemented.');
	}
</script>

{#if src}
	<div
		class="flex h-full bg-base-200 shadow-xl max-h-fit min-h-fit rounded-sm border-2 border-base-300"
	>
		<div class="flex-col h-full w-full">
			<div class="flex flex-col h-full w-full">
				<div class="flex flex-row h-full w-full grid-cols-3">
					<div class="flex-1 h-full bg-base-100 py-1 px-1">
						<h1 class="text-xl text-center font-semibold">{src.name}</h1>
						<FleetTransfer
							{transferAmount}
							{fuelTransferAmount}
							{cargo}
							{fuel}
							cargoCapacity={src.spec.cargoCapacity}
							fuelCapacity={src.spec.fuelCapacity}
						/>
					</div>
					<div class="flex-none h-full mx-0.5 w-20 px-1 mt-8">
						{#if dest?.type == MapObjectType.Fleet}
							<TransferButtons class="my-2" />
						{:else}
							<div class="h-6" />
						{/if}
						<div class="mt-16">
							<TransferButtons
								on:transfer-to-source={() => transferIronium(getQuantityModifier())}
								on:transfer-to-dest={() => transferIronium(-getQuantityModifier())}
								class="my-2"
							/>
							<TransferButtons
								on:transfer-to-source={() => transferBoranium(getQuantityModifier())}
								on:transfer-to-dest={() => transferBoranium(-getQuantityModifier())}
								class="my-2"
							/>
							<TransferButtons
								on:transfer-to-source={() => transferGermanium(getQuantityModifier())}
								on:transfer-to-dest={() => transferGermanium(-getQuantityModifier())}
								class="my-2"
							/>
							<TransferButtons
								on:transfer-to-source={() => transferColonists(getQuantityModifier())}
								on:transfer-to-dest={() => transferColonists(-getQuantityModifier())}
								class="my-2"
							/>
						</div>
					</div>
					<div class="flex-1 h-full bg-base-100 py-1 px-1">
						<div class="flex flex-col h-full">
							<h1 class="text-xl text-center font-semibold">
								{#if dest}
									{dest.name}
								{:else}
									Deep Space
								{/if}
							</h1>

							{#if dest?.type == MapObjectType.Planet}
								<PlanetTransfer cargo={dest.cargo} transferAmount={negativeCargo(transferAmount)} />
							{:else}
								Deep Space
							{/if}
						</div>
					</div>
				</div>
				<div class="flex justify-end pt-2">
					<button on:click={ok} class="btn btn-primary">Ok</button>
				</div>
			</div>
		</div>
	</div>
{/if}
