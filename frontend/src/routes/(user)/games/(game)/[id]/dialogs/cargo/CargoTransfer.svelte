<script lang="ts">
	import { quantityModifier } from '$lib/quantityModifier';
	import { getGameContext } from '$lib/services/Contexts';
	import { clamp } from '$lib/services/Math';
	import { add, negativeCargo, totalCargo, CargoTransferRequest } from '$lib/types/Cargo';
	import type { CommandedFleet, Fleet } from '$lib/types/Fleet';
	import { MapObjectType } from '$lib/types/MapObject';
	import type { Planet } from '$lib/types/Planet';
	import type { Salvage } from '$lib/types/Salvage';
	import hotkeys from 'hotkeys-js';
	import { createEventDispatcher, onMount } from 'svelte';
	import type { CargoTransferEvent } from './CargoTranfserDialog.svelte';
	import FleetTransfer from './FleetTransfer.svelte';
	import PlanetTransfer from './PlanetTransfer.svelte';
	import SalvageTransfer from './SalvageTransfer.svelte';
	import TransferButtons from './TransferButtons.svelte';

	const dispatch = createEventDispatcher<CargoTransferEvent>();

	const { game, player, universe } = getGameContext();

	export let src: CommandedFleet;
	export let dest: Fleet | Planet | Salvage | undefined;

	let transferAmount = new CargoTransferRequest();
	let cargo = new CargoTransferRequest();

	function reset() {
		transferAmount = new CargoTransferRequest();
		cargo = new CargoTransferRequest();
	}

	function ok() {
		dispatch('transfer-cargo', { src, dest, transferAmount });
		reset();
	}

	function cancel() {
		reset();
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

	$: src && src.cargo && (cargo = new CargoTransferRequest(src.cargo, src.fuel));
	$: destFleet = dest?.type === MapObjectType.Fleet ? (dest as Fleet) : undefined;

	// get the amount we actually transfer based on how much we have available to give/take
	// only give as much as we have up to the amount the dest can hold
	// only take as much as we can hold, up to the amount the dest has available
	function getTransferAmount(
		amountRequested: number,
		srcAmountAvailable: number, // the amount of i.e. ironium available to give
		destAmountAvailable: number // the amount of i.e. ironium available to take
	): number {
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
	}

	// get the amount we actually transfer based on how much we have available to give/take
	// only give as much as we have up to the amount the dest can hold
	// only take as much as we can hold, up to the amount the dest has available
	function getFuelTransferAmount(
		amountRequested: number,
		srcAmountAvailable: number, // the amount of i.e. fuel available to give
		destAmountAvailable: number // the amount of i.e. fuel available to take
	): number {
		if (dest && 'fuel' in dest) {
			// given the current transferAmount, figure out the current state of the source
			// and destination cargos
			const sourceFuel = src.fuel + transferAmount.fuel;
			const sourceCapacity = (src.spec.fuelCapacity ?? 0) - sourceFuel;
			const destFuel = (dest.fuel ?? 0) - transferAmount.fuel;
			const destCapacity = destFleet
				? (destFleet.spec?.fuelCapacity ?? 0) - destFuel
				: Number.MAX_SAFE_INTEGER;

			// if we have 30 fuel and they have 25 capacity, we can transfer -25 to them
			// if we have 30 fuel and they have 100 capacity, we can transfer -100 to them
			const canGiveAmount = Math.min(srcAmountAvailable, destCapacity);

			// if we have 30 fuel available on the dest, but only 5 of space, take only 5
			// if we have 30 fuel available on the dest, but 35 of space, take up to 30
			const canTakeAmount = Math.min(destAmountAvailable, sourceCapacity);

			return clamp(amountRequested, -canGiveAmount, canTakeAmount);
		} else {
			console.log("can't give/take any fuel");
			return 0;
		}
	}

	function transferFuel(amount: number) {
		// console.log('amount', amount, 'fuel', fuel, 'fuelTransferAmount', fuelTransferAmount);
		if (dest && 'fuel' in dest) {
			transferAmount.fuel =
				transferAmount.fuel +
				getFuelTransferAmount(
					amount,
					(cargo.fuel ?? 0) + transferAmount.fuel,
					(dest.fuel ?? 0) - transferAmount.fuel
				);
		}
	}

	function transferIronium(amount: number) {
		// update the amount we are transfering
		transferAmount.ironium =
			transferAmount.ironium +
			getTransferAmount(
				amount,
				(cargo.ironium ?? 0) + transferAmount.ironium,
				(dest?.cargo?.ironium ?? 0) - transferAmount.ironium
			);
	}

	function transferBoranium(amount: number) {
		// update the amount we are transfering
		transferAmount.boranium =
			transferAmount.boranium +
			getTransferAmount(
				amount,
				(cargo.boranium ?? 0) + transferAmount.boranium,
				(dest?.cargo?.boranium ?? 0) - transferAmount.boranium
			);
	}

	function transferGermanium(amount: number) {
		// update the amount we are transfering
		transferAmount.germanium =
			transferAmount.germanium +
			getTransferAmount(
				amount,
				(cargo.germanium ?? 0) + transferAmount.germanium,
				(dest?.cargo?.germanium ?? 0) - transferAmount.germanium
			);
	}

	function transferColonists(amount: number) {
		// update the amount we are transfering
		transferAmount.colonists =
			transferAmount.colonists +
			getTransferAmount(
				amount,
				(cargo.colonists ?? 0) + transferAmount.colonists,
				(dest?.cargo?.colonists ?? 0) - transferAmount.colonists
			);
	}
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
							{cargo}
							cargoCapacity={src.spec.cargoCapacity ?? 0}
							fuelCapacity={src.spec.fuelCapacity ?? 0}
							allowFuelTransfers={dest && 'fuel' in dest}
							on:transfer-fuel={(e) => transferFuel(e.detail)}
							on:transfer-ironium={(e) => transferIronium(e.detail)}
							on:transfer-boranium={(e) => transferBoranium(e.detail)}
							on:transfer-germanium={(e) => transferGermanium(e.detail)}
							on:transfer-colonists={(e) => transferColonists(e.detail)}
						/>
					</div>
					<div class="flex-none h-full mx-0.5 w-20 px-1 mt-8">
						{#if dest?.type == MapObjectType.Fleet}
							<TransferButtons
								on:transfer-to-source={(e) => transferFuel(quantityModifier(e.detail))}
								on:transfer-to-dest={(e) => transferFuel(-quantityModifier(e.detail))}
								class="my-2"
							/>
						{:else}
							<div class="h-6" />
						{/if}
						<div class="mt-16">
							<TransferButtons
								on:transfer-to-source={(e) => transferIronium(quantityModifier(e.detail))}
								on:transfer-to-dest={(e) => transferIronium(-quantityModifier(e.detail))}
								class="my-2"
							/>
							<TransferButtons
								on:transfer-to-source={(e) => transferBoranium(quantityModifier(e.detail))}
								on:transfer-to-dest={(e) => transferBoranium(-quantityModifier(e.detail))}
								class="my-2"
							/>
							<TransferButtons
								on:transfer-to-source={(e) => transferGermanium(quantityModifier(e.detail))}
								on:transfer-to-dest={(e) => transferGermanium(-quantityModifier(e.detail))}
								class="my-2"
							/>
							<TransferButtons
								on:transfer-to-source={(e) => transferColonists(quantityModifier(e.detail))}
								on:transfer-to-dest={(e) => transferColonists(-quantityModifier(e.detail))}
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
								<SalvageTransfer cargo={dest?.cargo} transferAmount={transferAmount.negative()} />
							{:else if destFleet}
								<FleetTransfer
									cargo={new CargoTransferRequest(destFleet.cargo, destFleet.fuel)}
									transferAmount={transferAmount.negative()}
									cargoCapacity={destFleet.spec?.cargoCapacity ?? 0}
									fuelCapacity={destFleet.spec?.fuelCapacity ?? 0}
									on:transfer-fuel={(e) => transferFuel(-e.detail)}
									on:transfer-ironium={(e) => transferIronium(-e.detail)}
									on:transfer-boranium={(e) => transferBoranium(-e.detail)}
									on:transfer-germanium={(e) => transferGermanium(-e.detail)}
									on:transfer-colonists={(e) => transferColonists(-e.detail)}
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
