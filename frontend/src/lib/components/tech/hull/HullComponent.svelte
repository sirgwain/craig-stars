<script lang="ts">
	import { techs, showPopupTech } from '$lib/services/Context';
	import type { ShipDesignSlot } from '$lib/types/ShipDesign';
	import { HullSlotType } from '$lib/types/Tech';
	import { Minus, Plus, Trash } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { kebabCase } from 'lodash-es';
	import { createEventDispatcher } from 'svelte';
	import { $enum as eu } from 'ts-enum-util';

	const dispatch = createEventDispatcher();

	export let type: HullSlotType = HullSlotType.General;
	export let capacity: number = 1;
	export let required = false;
	export let shipDesignSlot: ShipDesignSlot | undefined = undefined;
	export let highlighted = false;
	export let highlightedClass = 'border-accent';

	function typeDescription() {
		switch (type) {
			case HullSlotType.MineLayer:
				return 'Mine\nLayer';
			case HullSlotType.Mechanical:
				return 'Mech';
			case HullSlotType.SpaceDock:
				return 'Space Dock';
			case HullSlotType.ShieldArmor:
				return 'Shield\nor\nArmor';
			case HullSlotType.ShieldElectricalMechanical:
				return 'Shield\nElect\nMech';
			case HullSlotType.OrbitalElectrical:
				return 'Orbital\nor\nElectrical';
			case HullSlotType.WeaponShield:
				return 'Weapon\nor\nShield';
			case HullSlotType.ScannerElectricalMechanical:
				return 'Scanner\nElec\nMech';
			case HullSlotType.ArmorScannerElectricalMechanical:
				return 'Armor\nScanner\nElec/Mech';
			case HullSlotType.MineElectricalMechanical:
				return 'Mine\nElec\nMech';
			case HullSlotType.General:
				return 'General\nPurpose';
			default:
				return eu(HullSlotType).getKeyOrDefault(type, 'General');
		}
	}

	const icon = (c: string | undefined) => {
		return kebabCase(c?.replace("'", '').replace(' ', '').replace('±', ''));
	};
</script>

<div
	class={`flex bg-base-300 tech-avatar text-sm avatar ${icon(shipDesignSlot?.hullComponent)} ${
		highlighted ? highlightedClass : ''
	}`}
	class:border={!shipDesignSlot}
	class:border-2={highlighted}
	class:border-slate-900={!shipDesignSlot && !highlighted}
	class:z-20={highlighted}
>
	<button
		type="button"
		on:click={() => {
			dispatch('clicked');
		}}
		on:pointerdown|preventDefault={(e) =>
			shipDesignSlot?.hullComponent && !highlighted &&
			showPopupTech($techs.getHullComponent(shipDesignSlot?.hullComponent), e.x, e.y)}
		class="w-full h-full"
	>
		<div class="flex flex-col justify-between w-full h-full">
			{#if shipDesignSlot}
				<div class="grow">&nbsp;</div>
				<span class="h-[1rem] mt-auto text-center font-bold text-black"
					>{shipDesignSlot.quantity ?? 0} of {capacity}</span
				>
			{:else}
				<div class="grow whitespace-pre-wrap text-center">{typeDescription()}</div>
				{#if required}
					<div class="h-[1rem] mt-auto text-center text-red-500 font-bold">needs {capacity}</div>
				{:else}
					<span class="h-[1rem] mt-auto text-center font-bold">Up to {capacity}</span>
				{/if}
			{/if}
		</div>
	</button>
</div>
<div class="flex flex-row -ml-5 mt-1 gap-1" class:hidden={!highlighted || !shipDesignSlot}>
	<button
		type="button"
		class="btn btn-sm px-1 z-30"
		disabled={capacity === shipDesignSlot?.quantity}
		on:click={() => shipDesignSlot?.quantity && shipDesignSlot.quantity++}
	>
		<Icon src={Plus} size="24" class="hover:stroke-accent" />
	</button>
	<button
		type="button"
		class="btn btn-sm px-1 z-30"
		on:click={() => {
			if (shipDesignSlot?.quantity != undefined) {
				shipDesignSlot.quantity--;
				if (shipDesignSlot.quantity === 0) {
					dispatch('deleted');
				} else {
					dispatch('updated');
				}
			}
		}}
	>
		<Icon src={Minus} size="24" class="hover:stroke-accent" />
	</button>
	<button
		type="button"
		class="btn btn-sm px-1 z-30"
		on:click={() => {
			dispatch('deleted');
		}}
	>
		<Icon src={Trash} size="24" class="hover:stroke-accent" />
	</button>
</div>
