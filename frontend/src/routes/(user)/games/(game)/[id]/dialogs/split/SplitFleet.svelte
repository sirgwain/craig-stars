<script lang="ts" context="module">
	import { getGameContext } from '$lib/services/Contexts';
	import { CommandedFleet, type Fleet, type ShipToken } from '$lib/types/Fleet';
	import { ArrowLongLeft, ArrowLongRight } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	export type SplitFleetEventDetails = {
		src: CommandedFleet;
		dest: Fleet | undefined;
		srcTokens: ShipToken[];
		destTokens: ShipToken[];
		transferAmount: CargoTransferRequest;
	};

	export type SplitFleetEvent = {
		'split-fleet': SplitFleetEventDetails;
		'split-all': CommandedFleet;
		cancel: void;
	};
</script>

<script lang="ts">
	import FleetIcon from '$lib/components/FleetIcon.svelte';
	import CargoTransferer from '$lib/components/game/cargotransfer/CargoTransferer.svelte';
	import { quantityModifier } from '$lib/quantityModifier';
	import { CargoTransferRequest, emptyCargo, type Cargo } from '$lib/types/Cargo';
	import { total } from '$lib/types/Cost';
	import hotkeys from 'hotkeys-js';
	import { cloneDeep, maxBy } from 'lodash-es';
	import { createEventDispatcher, onMount } from 'svelte';

	const dispatch = createEventDispatcher<SplitFleetEvent>();
	const { game, player, universe } = getGameContext();

	export let src: CommandedFleet;
	export let dest: Fleet | undefined = undefined;

	let transferAmount = new CargoTransferRequest();
	let srcTokens: ShipToken[] = [];
	let destTokens: ShipToken[] = [];
	let srcFuelCapacity: number = src.spec.fuelCapacity ?? 0;
	let destFuelCapacity: number = dest?.spec?.fuelCapacity ?? 0;
	let srcCargoCapacity: number = src.spec.cargoCapacity ?? 0;
	let destCargoCapacity: number = dest?.spec?.cargoCapacity ?? 0;
	const totalFuel = src.fuel + (dest?.fuel ?? 0);

	function ok() {
		dispatch('split-fleet', { src, dest, srcTokens, destTokens, transferAmount });
	}

	function cancel() {
		dispatch('cancel');
	}

	// move some number of tokens from the source to the destination
	// if quantity is positive, this means moving source -> dest
	// if quantity is negative, this means moving dest -> source
	function moveToken(quantity: number, token: ShipToken, index: number) {
		const design = $universe.getMyDesign(token.designNum);
		if (!dest || !destTokens || !design || !quantity) {
			return;
		}

		const designFuelCapacity = design.spec.fuelCapacity ?? 0;
		const designCargoCapacity = design.spec.cargoCapacity ?? 0;

		// determine what percent of the total fleet's fuel belongs to these tokens
		const fuelPercent = (designFuelCapacity * quantity) / (srcFuelCapacity + destFuelCapacity);
		transferAmount.fuel -= Math.sign(fuelPercent) * Math.floor(Math.abs(totalFuel * fuelPercent));

		srcTokens[index].quantity -= quantity;
		destTokens[index].quantity += quantity;

		srcFuelCapacity -= designFuelCapacity * quantity;
		srcCargoCapacity -= designCargoCapacity * quantity;

		destFuelCapacity += designFuelCapacity * quantity;
		destCargoCapacity += designCargoCapacity * quantity;

		// if we have more cargo on the source than space available, move some out
		if (total(src.cargo) - transferAmount.absoluteCargoSize() > srcCargoCapacity) {
			let overload = total(src.cargo) - transferAmount.absoluteCargoSize() - srcCargoCapacity;

			let key: keyof Cargo;
			for (key in emptyCargo()) {
				// move over as much cargo as necessary
				const value = (src.cargo[key] ?? 0) + transferAmount[key];
				if ((value ?? 0) + transferAmount[key] > 0) {
					transferAmount[key] -= Math.min(value ?? 0, overload);
					overload -= Math.min(value ?? 0, overload);
				}
			}
		} else if (
			dest &&
			dest.cargo &&
			total(dest.cargo) + transferAmount.absoluteCargoSize() > destCargoCapacity
		) {
			let overload = total(dest.cargo) + transferAmount.absoluteCargoSize() - destCargoCapacity;

			let key: keyof Cargo;
			for (key in emptyCargo()) {
				// move over as much cargo as necessary
				const value = (dest.cargo[key] ?? 0) - transferAmount[key];
				if ((value ?? 0) - transferAmount[key] > 0) {
					transferAmount[key] += Math.min(value ?? 0, overload);
					overload -= Math.min(value ?? 0, overload);
				}
			}
		}
	}

	onMount(() => {
		const originalScope = hotkeys.getScope();
		const scope = 'cargoTransfer';
		hotkeys('Esc', scope, cancel);
		hotkeys('Enter', scope, ok);
		hotkeys.setScope(scope);

		if (!dest) {
			dest = new CommandedFleet(src);
			dest.num = 0;
			dest.spec = cloneDeep(src.spec);
			dest.name = `${dest.baseName}`;
			dest.tokens = src.tokens.map((t) => Object.assign({}, t, { quantity: 0 }));
			dest.spec.fuelCapacity = 0;
			dest.spec.cargoCapacity = 0;
			dest.fuel = 0;
			dest.cargo = {};

			srcTokens = cloneDeep(src.tokens);
			destTokens = cloneDeep(dest.tokens ?? []);
		} else {
			// we have a source and a dest, make the srcTokens and destTokens match up
			srcTokens = cloneDeep(src.tokens);
			destTokens = cloneDeep(src.tokens.map((t) => Object.assign({}, t, { quantity: 0 })));

			dest.tokens?.forEach((token) => {
				const tokenWithDesignInSrc = srcTokens.find((t) => t.designNum === token.designNum);
				if (!tokenWithDesignInSrc) {
					// this token only exists in the destination, so add a 0 quantity copy to the src
					srcTokens.push(Object.assign({}, token, { quantity: 0 }));
					destTokens.push(Object.assign({}, token));
				} else {
					// this token exists in the src, so update the quantity in the destination
					const tokenWithDesignInDest = destTokens.find((t) => t.designNum === token.designNum);
					if (tokenWithDesignInDest) {
						tokenWithDesignInDest.quantity = token.quantity;
					}
				}
			});
		}

		return () => {
			hotkeys.unbind('Esc', scope, cancel);
			hotkeys.unbind('Enter', scope, ok);
			hotkeys.deleteScope(scope);
			hotkeys.setScope(originalScope);
		};
	});
</script>

{#if dest && destTokens}
	<div class="flex flex-col px-1 w-full h-full">
		<div class="flex flex-col grow">
			<div class="text-xl font-semibold w-full text-center">Split Fleet</div>

			<!-- Token split page -->
			<div class="flex flex-row justify-around font-semibold mb-2">
				<!-- Source fleet -->
				<div class="grow flex flex-col">
					<div class="flex flex-col place-items-center">
						<FleetIcon fleet={src} tokens={srcTokens} />
						<div class="h-[2rem] font-semibold text-xl">{src.name}</div>
					</div>
					<div class="border border-secondary p-2">
						{#each srcTokens as token}
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
					<div class="h-[120px]" />
					<div class="grow p-2 flex flex-col justify-between">
						{#each srcTokens as token, index}
							<div class="flex flex-row h-full">
								<button
									on:click={(e) => {
										moveToken(-quantityModifier(e, 0, destTokens[index].quantity), token, index);
									}}
									class="btn btn-outline btn-xs normal-case btn-secondary inline-block p-1"
									><Icon src={ArrowLongLeft} size="16" class="hover:stroke-accent inline" />
								</button>
								<button
									on:click={(e) => {
										moveToken(quantityModifier(e, 0, srcTokens[index].quantity), token, index);
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
						<div class="h-[2rem] font-semibold text-xl">{dest.name}</div>
					</div>
					<div class="border border-secondary p-2">
						{#each destTokens as token}
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
				bind:transferAmount
				{srcCargoCapacity}
				{srcFuelCapacity}
				{destCargoCapacity}
				{destFuelCapacity}
			/>
		</div>
		<div class="flex flex-none justify-end pt-2 my-auto">
			<button on:click={ok} class="btn btn-primary">Ok</button>
			<button on:click={cancel} class="btn btn-secondary">Cancel</button>
		</div>
	</div>
{/if}
