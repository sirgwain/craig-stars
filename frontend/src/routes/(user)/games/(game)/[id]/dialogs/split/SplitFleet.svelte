<script lang="ts">
	import FleetIcon from '$lib/components/FleetIcon.svelte';
	import CargoTransferer from '$lib/components/game/cargotransfer/CargoTransferer.svelte';
	import type { OnCancel, OnOk, SplitFleetEvent } from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import { clamp } from '$lib/services/Math';
	import { emptyCargoJson, totalCargo } from '$lib/types/Cargo';
	import { absoluteCargoSize, CargoTransferRequest } from '$lib/types/CargoTransferRequest.svelte';
	import {
		CargoSchema,
		FleetSchema,
		FleetSpecSchema,
		GameDBObjectSchema,
		MapObjectSchema,
		ShipDesignSpecSchema,
		type CargoJson,
		type Fleet,
		type ShipToken
	} from '$lib/types/cs-proto';
	import { CommandedFleet, moveDamagedTokens } from '$lib/types/Fleet';
	import { clone, create } from '@bufbuild/protobuf';
	import { ArrowLongLeft, ArrowLongRight } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import hotkeys from 'hotkeys-js';
	import { cloneDeep } from 'lodash-es';
	import { onMount } from 'svelte';

	const { universe } = getGameContext();

	type Props = {
		src: CommandedFleet;
		dest?: Fleet | undefined;
		onOk?: OnOk<SplitFleetEvent>;
		onCancel?: OnCancel;
	};

	let { src, dest: destFleetProp, onOk, onCancel }: Props = $props();

	let transferAmount = $state(new CargoTransferRequest());
	let srcTokens: ShipToken[] = $state([]);
	let destTokens: ShipToken[] = $state([]);
	let dest = $state<Fleet>(
		destFleetProp
			? clone(FleetSchema, { ...destFleetProp, gameDbObject: create(GameDBObjectSchema) })
			: newEmptyDestFleet(src)
	);
	let srcFuelCapacity: number = $state(src.spec.shipDesignSpec?.fuelCapacity ?? 0);
	let destFuelCapacity: number = $state(destFleetProp?.spec?.shipDesignSpec?.fuelCapacity ?? 0);
	let srcCargoCapacity: number = $state(src.spec.shipDesignSpec?.cargoCapacity ?? 0);
	let destCargoCapacity: number = $state(destFleetProp?.spec?.shipDesignSpec?.cargoCapacity ?? 0);
	let quantityModifier = $state(1);

	const totalFuel = src.fuel + (destFleetProp?.fuel ?? 0);

	function split() {
		onOk?.({ src, dest, srcTokens, destTokens, transferAmount });
	}

	function cancel() {
		onCancel?.();
	}

	// move some number of tokens from the source to the destination
	// if quantity is positive, this means moving source -> dest
	// if quantity is negative, this means moving dest -> source
	function moveToken(quantity: number, token: ShipToken, index: number) {
		const design = $universe.getMyDesign(token.designNum);
		if (!dest.mapObject || !design || !quantity) {
			return;
		}

		const designFuelCapacity = design.spec?.fuelCapacity ?? 0;
		const designCargoCapacity = design.spec?.cargoCapacity ?? 0;

		// determine what percent of the total fleet's fuel belongs to these tokens
		const fuelPercent = (designFuelCapacity * quantity) / (srcFuelCapacity + destFuelCapacity);
		transferAmount.fuel -= Math.sign(fuelPercent) * Math.floor(Math.abs(totalFuel * fuelPercent));

		const srcToken = srcTokens[index];
		const destToken = destTokens[index];

		const destShipQuantity = destTokens.reduce((count, token) => count + token.quantity, 0);

		srcToken.quantity -= quantity;
		destToken.quantity += quantity;

		if (quantity > 0) {
			// update the dest fleet name to be the name of the first token moved
			if (destShipQuantity == 0) {
				dest.baseName = design.name;
				dest.mapObject.name = dest.baseName;
			}
			// move from left to right
			moveDamagedTokens(srcToken, destToken, quantity);
		} else if (quantity < 0) {
			// move from right to left
			moveDamagedTokens(destToken, srcToken, -quantity);
		}

		srcTokens[index] = srcToken;
		destTokens[index] = destToken;

		srcFuelCapacity -= designFuelCapacity * quantity;
		srcCargoCapacity -= designCargoCapacity * quantity;

		destFuelCapacity += designFuelCapacity * quantity;
		destCargoCapacity += designCargoCapacity * quantity;

		// if we have more cargo on the source than space available, move some out
		if (totalCargo(src.cargo) - absoluteCargoSize(transferAmount) > srcCargoCapacity) {
			let overload = totalCargo(src.cargo) - absoluteCargoSize(transferAmount) - srcCargoCapacity;

			let key: keyof CargoJson;
			for (key in emptyCargoJson()) {
				// the value of the source including what we've already transferred in/out
				const value = src.cargo[key] + transferAmount[key];
				if (value > 0) {
					// we have some left to transfer
					transferAmount[key] -= Math.min(value, overload);
					overload -= Math.min(value, overload);
				}
			}
		} else if (
			dest.cargo &&
			totalCargo(dest.cargo) + absoluteCargoSize(transferAmount) > destCargoCapacity
		) {
			let overload = totalCargo(dest.cargo) + absoluteCargoSize(transferAmount) - destCargoCapacity;

			let key: keyof CargoJson;
			for (key in emptyCargoJson()) {
				// the value of the dest including what we've already transferred in/out
				const value = dest.cargo[key] - transferAmount[key];
				if (value > 0) {
					transferAmount[key] += Math.min(value, overload);
					overload -= Math.min(value, overload);
				}
			}
		}
		srcTokens = srcTokens;
		destTokens = destTokens;
	}

	function newEmptyDestFleet(src: CommandedFleet): Fleet {
		const fleet: Fleet = clone(FleetSchema, { ...src, gameDbObject: create(GameDBObjectSchema) });
		fleet.mapObject = fleet.mapObject ?? create(MapObjectSchema);
		fleet.mapObject.num = 0;
		fleet.spec = clone(FleetSpecSchema, src.spec);
		fleet.mapObject.name = `${fleet.baseName}`;
		fleet.tokens = src.tokens.map((t) =>
			Object.assign({}, t, { quantity: 0, quantityDamaged: 0, damage: 0 })
		);
		fleet.spec.shipDesignSpec = create(ShipDesignSpecSchema, {
			fuelCapacity: 0,
			cargoCapacity: 0
		});

		fleet.fuel = 0;
		fleet.cargo = create(CargoSchema);

		return fleet;
	}

	onMount(() => {
		const originalScope = hotkeys.getScope();
		const scope = 'cargoTransfer';
		hotkeys('Esc', scope, cancel);
		hotkeys('Enter', scope, split);
		hotkeys.setScope(scope);

		if (!destFleetProp) {
			srcTokens = cloneDeep(src.tokens);
			destTokens = cloneDeep(dest.tokens);
		} else {
			// we have a source and a dest, make the srcTokens and destTokens match up
			srcTokens = cloneDeep(src.tokens);
			destTokens = cloneDeep(
				src.tokens.map((t) => Object.assign({}, t, { quantity: 0, quantityDamaged: 0, damage: 0 }))
			);

			dest.tokens.forEach((token) => {
				const tokenWithDesignInSrc = srcTokens.find((t) => t.designNum === token.designNum);
				if (!tokenWithDesignInSrc) {
					// this token only exists in the destination, so add a 0 quantity copy to the src
					srcTokens.push(Object.assign({}, token, { quantity: 0, quantityDamaged: 0, damage: 0 }));
					destTokens.push(Object.assign({}, token));
				} else {
					// this token exists in the src, so update the quantity in the destination
					const tokenWithDesignInDest = destTokens.find((t) => t.designNum === token.designNum);
					if (tokenWithDesignInDest) {
						tokenWithDesignInDest.quantity = token.quantity;
						tokenWithDesignInDest.quantityDamaged = token.quantityDamaged;
						tokenWithDesignInDest.damage = token.damage;
					}
				}
			});
		}

		return () => {
			hotkeys.unbind('Esc', scope, cancel);
			hotkeys.unbind('Enter', scope, split);
			hotkeys.deleteScope(scope);
			hotkeys.setScope(originalScope);
		};
	});
</script>

{#if destTokens}
	<div class="flex flex-col px-1 w-full h-full">
		<div class="flex flex-col grow">
			<div class="text-xl font-semibold w-full text-center">Split Fleet</div>

			<!-- Token split page -->
			<div class="flex flex-row justify-around font-semibold mb-2">
				<!-- Source fleet -->
				<div class="grow flex flex-col">
					<div class="flex flex-col place-items-center">
						<FleetIcon fleet={src} tokens={srcTokens} />
						<div class="h-[2rem] font-semibold text-xl">{src.mapObject.name}</div>
					</div>
					<div class="border border-secondary p-2">
						{#each srcTokens as token (token)}
							{@const design = $universe.getMyDesign(token.designNum)}
							<div class="flex flex-row gap-1 min-h-8">
								<div class="grow flex flex-row">
									<div class="text-right pr-1 w-full">{design?.name}</div>
									<div
										class="w-16 sm:w-20 h-8 my-auto ml-auto px-1 text-right border border-secondary"
									>
										{token.quantity}
									</div>
								</div>
							</div>
						{/each}
					</div>
				</div>
				<!-- buttons -->
				<div class="flex-none flex flex-col">
					<!-- Keep a 2rem empty header so the buttons line up -->
					<div class="h-[120px]"></div>
					<div class="grow p-2 flex flex-col justify-between">
						{#each srcTokens as token, index (token)}
							<div class="flex flex-row h-full">
								<button
									type="button"
									onclick={() => {
										moveToken(
											-clamp(quantityModifier, 0, destTokens[index].quantity),
											token,
											index
										);
									}}
									class="btn btn-outline btn-xs normal-case btn-secondary inline-block p-1"
									><Icon src={ArrowLongLeft} size="16" class="hover:stroke-accent inline" />
								</button>
								<button
									type="button"
									onclick={() => {
										moveToken(clamp(quantityModifier, 0, srcTokens[index].quantity), token, index);
									}}
									class="btn btn-outline btn-xs normal-case btn-secondary inline-block p-1"
									><Icon
										src={ArrowLongRight}
										size="16"
										class="hover:stroke-accent inline"
									/></button
								>
							</div>
						{/each}
					</div>
				</div>
				<!-- Dest fleet -->
				<div class="grow flex flex-col">
					<div class="flex flex-col place-items-center">
						<FleetIcon fleet={dest} tokens={destTokens} />
						<div class="h-[2rem] font-semibold text-xl">{dest.mapObject?.name}</div>
					</div>
					<div class="border border-secondary p-2">
						{#each destTokens as token (token)}
							{@const design = $universe.getMyDesign(token.designNum)}
							<div class="flex flex-row gap-1 min-h-8">
								<div class="grow flex flex-row">
									<div class="text-right pr-1 w-full">{design?.name}</div>
									<div
										class="w-16 sm:w-20 h-8 my-auto ml-auto px-1 text-right border border-secondary"
									>
										{token.quantity}
									</div>
								</div>
							</div>
						{/each}
					</div>
				</div>
			</div>

			<CargoTransferer
				{src}
				{dest}
				showHeader={false}
				{srcCargoCapacity}
				{srcFuelCapacity}
				{destCargoCapacity}
				{destFuelCapacity}
				bind:transferAmount
				bind:quantityModifier
			/>
		</div>
		<div class="flex flex-none justify-end pt-2 my-auto">
			<button onclick={split} class="btn btn-primary">Ok</button>
			<button onclick={onCancel} class="btn btn-secondary">Cancel</button>
		</div>
	</div>
{/if}
