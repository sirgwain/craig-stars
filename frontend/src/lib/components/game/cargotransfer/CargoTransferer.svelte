<script lang="ts">
	import QuantityModifierButtons from '$lib/components/QuantityModifierButtons.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { clamp } from '$lib/services/Math';
	import type { AnyFleet } from '$lib/services/Universe';
	import { add, negativeCargo, totalCargo } from '$lib/types/Cargo';
	import type { CargoDest } from '$lib/types/CargoTransferRequest.svelte';
	import { CargoTransferRequest, negative } from '$lib/types/CargoTransferRequest.svelte';
	import {
		Colonists,
		Fuel,
		MapObjectTypeFleet,
		MapObjectTypeMineralPacket,
		MapObjectTypePlanet,
		MapObjectTypeSalvage
	} from '$lib/types/cs';
	import { canTransferCargoType, type CommandedFleet } from '$lib/types/Fleet';
	import FleetTransfer from './FleetTransfer.svelte';
	import MineralPacketTransfer from './MineralPacketTransfer.svelte';
	import PlanetTransfer from './PlanetTransfer.svelte';
	import SalvageTransfer from './SalvageTransfer.svelte';
	import TransferButtons from './TransferButtons.svelte';

	const { player } = getGameContext();

	type Props = {
		src: CommandedFleet;
		dest: CargoDest;
		transferAmount?: CargoTransferRequest;
		showHeader?: boolean;
		srcCargoCapacity?: number;
		srcFuelCapacity?: number;
		destCargoCapacity?: number;
		destFuelCapacity?: number;
		quantityModifier?: number;
	};

	let {
		src,
		dest,
		transferAmount = $bindable(new CargoTransferRequest()),
		showHeader = true,
		srcCargoCapacity = src.spec.cargoCapacity ?? 0,
		srcFuelCapacity = src.spec.fuelCapacity ?? 0,
		destCargoCapacity = getCargoCapacity(dest),
		destFuelCapacity = getFuelCapacity(dest),
		quantityModifier = $bindable(1)
	}: Props = $props();

	let srcCargo = $derived(new CargoTransferRequest(src.cargo, src.fuel));
	let destCargo = $derived(
		new CargoTransferRequest(
			dest ? dest.cargo : $player.getByHandTransfer({ targetPosition: src.position }), // we are either tranfering to a location, or jettisoning
			dest && 'fuel' in dest ? dest.fuel : 0
		)
	);

	let destFleet = $derived(dest?.type === MapObjectTypeFleet ? (dest as AnyFleet) : undefined);

	function getCargoCapacity(dest: CargoDest): number {
		if (dest && 'spec' in dest && dest.spec && 'cargoCapacity' in dest.spec) {
			return dest.spec.cargoCapacity ?? 0;
		}
		if (dest?.type === MapObjectTypeFleet) {
			// if a fleet doesn't have a capcity, it's 0
			return 0;
		}

		// no limit on how much cargo we can dump here
		return Number.MAX_SAFE_INTEGER;
	}

	function getFuelCapacity(dest: CargoDest): number {
		if (dest && 'spec' in dest && dest.spec && 'fuelCapacity' in dest.spec) {
			return dest.spec.fuelCapacity ?? 0;
		}
		// can't dump fuel here
		return 0;
	}

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
		const updatedSourceCargo = add(srcCargo, transferAmount);
		const sourceRemainingCapacity = srcCargoCapacity - totalCargo(updatedSourceCargo);
		const updatedDestCargo = add(destCargo, negativeCargo(transferAmount));

		const destRemainingCapacity = destFleet
			? destCargoCapacity - totalCargo(updatedDestCargo)
			: Number.MAX_SAFE_INTEGER;

		// if we have 30 ironium and they have 25 capacity, we can transfer -25 to them
		// if we have 30 ironium and they have 100 capacity, we can transfer -100 to them
		const canGiveAmount = Math.min(srcAmountAvailable, destRemainingCapacity);

		// if we have 30kT ironium available on the dest, but only 5kT of space, take only 5kT
		// if we have 30kT ironium available on the dest, but 35kT of space, take up to 30kT
		const canTakeAmount = Math.min(destAmountAvailable, sourceRemainingCapacity);

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
			const updatedSourceFuel = srcCargo.fuel + transferAmount.fuel;
			const sourceRemainingCapacity = (srcFuelCapacity ?? 0) - updatedSourceFuel;
			const updatedDestFuel = destCargo.fuel - transferAmount.fuel;
			const destRemainingCapacity = destFleet
				? destFuelCapacity - updatedDestFuel
				: Number.MAX_SAFE_INTEGER;

			// if we have 30 fuel and they have 25 capacity, we can transfer -25 to them
			// if we have 30 fuel and they have 100 capacity, we can transfer -100 to them
			const canGiveAmount = Math.min(srcAmountAvailable, destRemainingCapacity);

			// if we have 30 fuel available on the dest, but only 5 of space, take only 5
			// if we have 30 fuel available on the dest, but 35 of space, take up to 30
			const canTakeAmount = Math.min(destAmountAvailable, sourceRemainingCapacity);

			return clamp(amountRequested, -canGiveAmount, canTakeAmount);
		} else {
			console.log("can't give/take any fuel");
			return 0;
		}
	}

	function transferFuel(amount: number): { src: number; dest: number } {
		// console.log('amount', amount, 'fuel', fuel, 'fuelTransferAmount', fuelTransferAmount);
		if (dest && 'fuel' in dest) {
			transferAmount.fuel =
				transferAmount.fuel +
				getFuelTransferAmount(
					amount,
					srcCargo.fuel + transferAmount.fuel,
					dest.fuel - transferAmount.fuel
				);
		}
		return { src: srcCargo.fuel + transferAmount.fuel, dest: destCargo.fuel - transferAmount.fuel };
	}

	function transferIronium(amount: number): { src: number; dest: number } {
		// update the amount we are transfering
		transferAmount.ironium =
			transferAmount.ironium +
			getTransferAmount(
				amount,
				(srcCargo.ironium ?? 0) + transferAmount.ironium,
				(destCargo?.ironium ?? 0) - transferAmount.ironium
			);
		return {
			src: (srcCargo.ironium ?? 0) + transferAmount.ironium,
			dest: (destCargo?.ironium ?? 0) - transferAmount.ironium
		};
	}

	function transferBoranium(amount: number): { src: number; dest: number } {
		// update the amount we are transfering
		transferAmount.boranium =
			transferAmount.boranium +
			getTransferAmount(
				amount,
				(srcCargo.boranium ?? 0) + transferAmount.boranium,
				(destCargo?.boranium ?? 0) - transferAmount.boranium
			);
		return {
			src: (srcCargo.boranium ?? 0) + transferAmount.boranium,
			dest: (destCargo?.boranium ?? 0) - transferAmount.boranium
		};
	}

	function transferGermanium(amount: number): { src: number; dest: number } {
		// update the amount we are transfering
		transferAmount.germanium =
			transferAmount.germanium +
			getTransferAmount(
				amount,
				(srcCargo.germanium ?? 0) + transferAmount.germanium,
				(destCargo?.germanium ?? 0) - transferAmount.germanium
			);
		return {
			src: (srcCargo.germanium ?? 0) + transferAmount.germanium,
			dest: (destCargo?.germanium ?? 0) - transferAmount.germanium
		};
	}

	function transferColonists(amount: number): { src: number; dest: number } {
		// update the amount we are transfering
		transferAmount.colonists =
			transferAmount.colonists +
			getTransferAmount(
				amount,
				(srcCargo.colonists ?? 0) + transferAmount.colonists,
				(destCargo?.colonists ?? 0) - transferAmount.colonists
			);
		return {
			src: (srcCargo.colonists ?? 0) + transferAmount.colonists,
			dest: (destCargo?.colonists ?? 0) - transferAmount.colonists
		};
	}
</script>

{#if src?.spec}
	<div class="flex flex-row h-full w-full grid-cols-3">
		<div class="flex-1 h-full bg-base-100 py-1 px-1">
			<h1 class="text-xl text-center font-semibold h-[2rem]">
				<span class:hidden={!showHeader}>
					{src.name}
				</span>
			</h1>
			<FleetTransfer
				{transferAmount}
				cargo={srcCargo}
				cargoCapacity={srcCargoCapacity}
				fuelCapacity={srcFuelCapacity}
				allowFuelTransfers={dest && 'fuel' in dest && canTransferCargoType(src, dest, Fuel)}
				allowColonistTransfers={canTransferCargoType(src, dest, Colonists)}
				onTransferFuel={(amount) => transferFuel(amount).src}
				onTransferIronium={(amount) => transferIronium(amount).src}
				onTransferBoranium={(amount) => transferBoranium(amount).src}
				onTransferGermanium={(amount) => transferGermanium(amount).src}
				onTransferColonists={(amount) => transferColonists(amount).src}
			/>
		</div>
		<div class="flex-none flex flex-col mx-0.5 w-20 px-1 mt-8">
			{#if dest?.type == MapObjectTypeFleet}
				<TransferButtons
					data-id="fuel"
					disabled={!canTransferCargoType(src, dest, Fuel)}
					onTransferToSource={() => transferFuel(quantityModifier)}
					onTransferToDest={() => transferFuel(-quantityModifier)}
					class="mt-8 sm:mt-2"
				/>
			{:else}
				<div class="h-8"></div>
			{/if}
			<div class="mt-28 h-40 sm:mt-16 sm:h-28 flex flex-col justify-between">
				<TransferButtons
					data-id="ironium"
					onTransferToSource={() => transferIronium(quantityModifier)}
					onTransferToDest={() => transferIronium(-quantityModifier)}
				/>
				<TransferButtons
					data-id="boranium"
					onTransferToSource={() => transferBoranium(quantityModifier)}
					onTransferToDest={() => transferBoranium(-quantityModifier)}
				/>
				<TransferButtons
					data-id="germanium"
					onTransferToSource={() => transferGermanium(quantityModifier)}
					onTransferToDest={() => transferGermanium(-quantityModifier)}
				/>
				<TransferButtons
					data-id="colonists"
					disabled={!canTransferCargoType(src, dest, Colonists)}
					onTransferToSource={() => transferColonists(quantityModifier)}
					onTransferToDest={() => transferColonists(-quantityModifier)}
				/>
			</div>
			<div class="flex flex-col justify-between mt-2 gap-1 mx-1">
				<QuantityModifierButtons bind:modifier={quantityModifier} />
			</div>
		</div>
		<div class="flex-1 h-full bg-base-100 py-1 px-1">
			<div class="flex flex-col h-full">
				<h1 class="text-xl text-center font-semibold h-[2rem]">
					<span class:hidden={!showHeader}>
						{#if dest && dest.name}
							{dest.name}
						{:else}
							Deep Space
						{/if}
					</span>
				</h1>

				{#if dest?.type == MapObjectTypePlanet}
					<PlanetTransfer cargo={destCargo} transferAmount={negativeCargo(transferAmount)} />
				{:else if !dest || dest?.type == MapObjectTypeSalvage}
					<SalvageTransfer cargo={destCargo} transferAmount={negative(transferAmount)} />
				{:else if !dest || dest?.type == MapObjectTypeMineralPacket}
					<MineralPacketTransfer cargo={destCargo} transferAmount={negative(transferAmount)} />
				{:else if destFleet}
					<FleetTransfer
						cargo={destCargo}
						transferAmount={negative(transferAmount)}
						cargoCapacity={destCargoCapacity}
						fuelCapacity={destFuelCapacity}
						allowFuelTransfers={canTransferCargoType(src, dest, Fuel)}
						allowColonistTransfers={canTransferCargoType(src, dest, Colonists)}
						onTransferFuel={(amount) => transferFuel(-amount).dest}
						onTransferIronium={(amount) => transferIronium(-amount).dest}
						onTransferBoranium={(amount) => transferBoranium(-amount).dest}
						onTransferGermanium={(amount) => transferGermanium(-amount).dest}
						onTransferColonists={(amount) => transferColonists(-amount).dest}
					/>
				{:else}
					Deep Space
				{/if}
			</div>
		</div>
	</div>
{/if}
