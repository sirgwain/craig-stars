<script lang="ts">
	import { getQuantityModifier } from '$lib/quantityModifier';
	import { getGameContext } from '$lib/services/Contexts';
	import { clamp } from '$lib/services/Math';
	import { emptyCargo, negativeCargo, totalCargo, type Cargo, add } from '$lib/types/Cargo';
	import type { CommandedFleet, Fleet } from '$lib/types/Fleet';
	import { MapObjectType } from '$lib/types/MapObject';
	import type { Planet } from '$lib/types/Planet';
	import hotkeys from 'hotkeys-js';
	import { createEventDispatcher } from 'svelte';
	import type { CargoTransferEvent } from './CargoTranfserDialog.svelte';
	import FleetTransfer from './FleetTransfer.svelte';
	import PlanetTransfer from './PlanetTransfer.svelte';
	import TransferButtons from './TransferButtons.svelte';
	import SalvageTransfer from './SalvageTransfer.svelte';
	import type { Salvage } from '$lib/types/Salvage';

	const dispatch = createEventDispatcher<CargoTransferEvent>();

	const { game, player, universe } = getGameContext();

	export let src: CommandedFleet;
	export let dest: Fleet | Planet | Salvage | undefined;

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
		dispatch('transfer-cargo', { src, dest, transferAmount });
		reset();
	};

	const cancel = () => {
		reset();
		dispatch('cancel');
	};

	hotkeys('Esc', () => cancel());
	hotkeys('Enter', () => {
		ok();
	});

	$: src && src.cargo && (cargo = src.cargo);
	$: src && src.fuel && (fuel = src.fuel);

	$: destFleet = dest?.type === MapObjectType.Fleet ? (dest as Fleet) : undefined;

	// get the amount we actually transfer based on how much we have available to give/take
	// only give as much as we have up to the amount the dest can hold
	// only take as much as we can hold, up to the amount the dest has available
	const getTransferAmount = (
		amountRequested: number,
		srcAmountAvailable: number, // the amount of i.e. ironium available to give
		destAmountAvailable: number // the amount of i.e. ironium available to take
	): number => {
		// given the current transferAmount, figure out the current state of the source
		// and destination cargos
		const sourceCargo = add(src.cargo, transferAmount);
		const sourceCapacity = (src.spec.cargoCapacity ?? 0) - totalCargo(sourceCargo);
		const destCargo = add(dest?.cargo ?? {}, negativeCargo(transferAmount));
		const destCapacity = destFleet
			? (destFleet.spec?.cargoCapacity ?? 0) - totalCargo(destCargo)
			: Number.MAX_SAFE_INTEGER;

		// if we have 30 ironium and they have 25 capacity, we can transfer -25 to them
		// if we have 30 ironium and they have 100 capacity, we can transfer -100 to them
		const canGiveAmount = Math.min(srcAmountAvailable, destCapacity);

		// if we have 30kT ironium available on the dest, but only 5kT of space, take only 5kT
		// if we have 30kT ironium available on the dest, but 35kT of space, take up to 30kT
		const canTakeAmount = Math.min(destAmountAvailable, sourceCapacity);

		return clamp(amountRequested, -canGiveAmount, canTakeAmount);
	};

	const transferIronium = (amount: number) => {
		// update the amount we are transfering
		transferAmount.ironium =
			(transferAmount.ironium ?? 0) +
			getTransferAmount(
				amount,
				(cargo.ironium ?? 0) + (transferAmount.ironium ?? 0),
				(dest?.cargo?.ironium ?? 0) - (transferAmount.ironium ?? 0)
			);
	};

	const transferBoranium = (amount: number) => {
		// update the amount we are transfering
		transferAmount.boranium =
			(transferAmount.boranium ?? 0) +
			getTransferAmount(
				amount,
				(cargo.boranium ?? 0) + (transferAmount.boranium ?? 0),
				(dest?.cargo?.boranium ?? 0) - (transferAmount.boranium ?? 0)
			);
	};

	const transferGermanium = (amount: number) => {
		// update the amount we are transfering
		transferAmount.germanium =
			(transferAmount.germanium ?? 0) +
			getTransferAmount(
				amount,
				(cargo.germanium ?? 0) + (transferAmount.germanium ?? 0),
				(dest?.cargo?.germanium ?? 0) - (transferAmount.germanium ?? 0)
			);
	};

	const transferColonists = (amount: number) => {
		// update the amount we are transfering
		transferAmount.colonists =
			(transferAmount.colonists ?? 0) +
			getTransferAmount(
				amount,
				(cargo.colonists ?? 0) + (transferAmount.colonists ?? 0),
				(dest?.cargo?.colonists ?? 0) - (transferAmount.colonists ?? 0)
			);
	};
</script>

{#if src?.spec}
	<div
		class="flex h-full bg-base-200 shadow max-h-fit min-h-fit rounded-sm border-2 border-base-300"
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
							on:transferironium={(e) => transferIronium(e.detail)}
							on:transferboranium={(e) => transferBoranium(e.detail)}
							on:transfergermanium={(e) => transferGermanium(e.detail)}
							on:transfercolonists={(e) => transferColonists(e.detail)}
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
								{#if dest && dest.name}
									{dest.name}
								{:else}
									Deep Space
								{/if}
							</h1>

							{#if dest?.type == MapObjectType.Planet}
								<PlanetTransfer cargo={dest.cargo} transferAmount={negativeCargo(transferAmount)} />
							{:else if !dest || dest?.type == MapObjectType.Salvage}
								<SalvageTransfer
									cargo={dest?.cargo}
									transferAmount={negativeCargo(transferAmount)}
								/>
							{:else if destFleet}
								<FleetTransfer
									cargo={destFleet.cargo}
									transferAmount={negativeCargo(transferAmount)}
									fuelTransferAmount={-fuelTransferAmount}
									fuel={destFleet.fuel}
									cargoCapacity={destFleet.spec?.cargoCapacity ?? 0}
									fuelCapacity={destFleet.spec?.fuelCapacity ?? 0}
									on:transferironium={(e) => transferIronium(-e.detail)}
									on:transferboranium={(e) => transferBoranium(-e.detail)}
									on:transfergermanium={(e) => transferGermanium(-e.detail)}
									on:transfercolonists={(e) => transferColonists(-e.detail)}
								/>
							{:else}
								Deep Space
							{/if}
						</div>
					</div>
				</div>
				<div class="flex justify-end pt-2">
					<button on:click={ok} class="btn btn-primary">Ok</button>
					<button on:click={cancel} class="btn btn-secondary">Cancel</button>
				</div>
			</div>
		</div>
	</div>
{/if}
