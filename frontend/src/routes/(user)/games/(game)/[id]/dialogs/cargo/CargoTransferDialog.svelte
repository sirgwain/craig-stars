<script lang="ts">
	import { getQuantityModifier } from '$lib/quantityModifier';
	import { FleetService } from '$lib/services/FleetService';
	import { clamp } from '$lib/services/Math';
	import { PlanetService } from '$lib/services/PlanetService';
	import { emptyCargo, negativeCargo, subtract, totalCargo, type Cargo } from '$lib/types/Cargo';
	import type { CommandedFleet, Fleet } from '$lib/types/Fleet';
	import { MapObjectType } from '$lib/types/MapObject';
	import type { Planet } from '$lib/types/Planet';
	import hotkeys from 'hotkeys-js';
	import { createEventDispatcher } from 'svelte';
	import FleetTransfer from './FleetTransfer.svelte';
	import PlanetTransfer from './PlanetTransfer.svelte';
	import TransferButtons from './TransferButtons.svelte';
	import { game } from '$lib/services/Context';

	export let src: CommandedFleet | undefined;
	export let dest: Fleet | Planet | undefined;

	let transferAmount: Cargo = emptyCargo();
	let cargo: Cargo = emptyCargo();
	let fuelTransferAmount: number = 0;
	let fuel: number = 0;

	const reset = () => {
		transferAmount = emptyCargo();
		fuelTransferAmount = 0;
		cargo = emptyCargo();
		fuel = 0;
	};

	const ok = async () => {
		if (src && $game) {
			const fleetService = new FleetService();
			try {
				const result = await FleetService.transferCargo($game.id, src, dest, transferAmount);
				// TODO: should the parent do this? or do we update the src/dest here?
				src.cargo = result.cargo;
				if (dest?.cargo) {
					dest.cargo = subtract(dest.cargo, transferAmount);
				}
				reset();
				dispatch('ok', result);
			} catch (e) {
				console.error(e);
			}
		}
	};

	const cancel = () => {
		reset();
		dispatch('cancel');
	};

	// get the available space to transfer into this fleet based on the current transferAmount set in the UI, the cargo capacity,
	// and the total cargo already in the hold when we opened the dialog
	const availableSpace = () =>
		(src?.spec?.cargoCapacity ?? 0) - totalCargo(transferAmount) - totalCargo(src?.cargo ?? {});

	const transferIronium = (amount: number) => {
		transferAmount.ironium = clamp(
			(transferAmount.ironium ?? 0) + amount,
			-(cargo.ironium ?? 0),
			(transferAmount.ironium ?? 0) + availableSpace()
		);
	};

	const transferBoranium = (amount: number) => {
		transferAmount.boranium = clamp(
			(transferAmount.boranium ?? 0) + amount,
			-(cargo.boranium ?? 0),
			(transferAmount.boranium ?? 0) + availableSpace()
		);
	};

	const transferGermanium = (amount: number) => {
		transferAmount.germanium = clamp(
			(transferAmount.germanium ?? 0) + amount,
			-(cargo.germanium ?? 0),
			(transferAmount.germanium ?? 0) + availableSpace()
		);
	};

	const transferColonists = (amount: number) => {
		transferAmount.colonists = clamp(
			(transferAmount.colonists ?? 0) + amount,
			-(cargo.colonists ?? 0),
			(transferAmount.colonists ?? 0) + availableSpace()
		);
	};

	const dispatch = createEventDispatcher();

	hotkeys('Esc', () => cancel());
	hotkeys('Enter', () => {
		ok();
	});

	const planetService = new PlanetService();

	$: src && src.cargo && (cargo = src.cargo);
	$: src && src.fuel && (fuel = src.fuel);
</script>

{#if src?.spec}
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
							cargoCapacity={src.spec.cargoCapacity ?? 0}
							fuelCapacity={src.spec.fuelCapacity ?? 0}
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
