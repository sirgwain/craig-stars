<script lang="ts">
	import { onTechTooltip } from '$lib/components/game/tooltips/TechTooltip.svelte';
	import { techs } from '$lib/services/Stores';
	import {
		HullSlotTypeArmor,
		HullSlotTypeArmorScannerElectricalMechanical,
		HullSlotTypeBomb,
		HullSlotTypeCargo,
		HullSlotTypeElectrical,
		HullSlotTypeElectricalMechanical,
		HullSlotTypeEngine,
		HullSlotTypeGeneral,
		HullSlotTypeMechanical,
		HullSlotTypeMineElectricalMechanical,
		HullSlotTypeMineLayer,
		HullSlotTypeMining,
		HullSlotTypeNone,
		HullSlotTypeOrbital,
		HullSlotTypeOrbitalElectrical,
		HullSlotTypeScanner,
		HullSlotTypeScannerElectricalMechanical,
		HullSlotTypeShield,
		HullSlotTypeShieldArmor,
		HullSlotTypeShieldElectricalMechanical,
		HullSlotTypeSpaceDock,
		HullSlotTypeWeapon,
		HullSlotTypeWeaponShield,
		type HullSlotType
	} from '$lib/types/Consts';
	import type { ShipDesignSlot } from '$lib/types/cs-proto';
	import { Minus, Plus, Trash } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { kebabCase } from 'lodash-es';

	type Props = {
		type?: HullSlotType;
		capacity?: number;
		required?: boolean;
		shipDesignSlot?: ShipDesignSlot | undefined;
		highlighted?: boolean;
		highlightedClass?: string;
		showTooltips?: boolean;
		onClick?: () => void;
		onDelete?: () => void;
		onUpdate?: () => void;
	};

	let {
		type = HullSlotTypeGeneral,
		capacity = 1,
		required = false,
		shipDesignSlot = $bindable(),
		highlighted = false,
		highlightedClass = 'border-accent',
		showTooltips = false,
		onClick: onclick,
		onDelete: onDelete,
		onUpdate: onUpdate
	}: Props = $props();

	function typeDescription() {
		switch (type) {
			case HullSlotTypeNone:
				return 'None';
			case HullSlotTypeEngine:
				return 'Engine';
			case HullSlotTypeScanner:
				return 'Scanner';
			case HullSlotTypeBomb:
				return 'Bomb';
			case HullSlotTypeMining:
				return 'Mining';
			case HullSlotTypeElectrical:
				return 'Electrical';
			case HullSlotTypeShield:
				return 'Shield';
			case HullSlotTypeArmor:
				return 'Armor';
			case HullSlotTypeCargo:
				return 'Cargo';
			case HullSlotTypeWeapon:
				return 'Weapon';
			case HullSlotTypeOrbital:
				return 'Orbital';
			case HullSlotTypeMineLayer:
				return 'Mine\nLayer';
			case HullSlotTypeMechanical:
				return 'Mech';
			case HullSlotTypeSpaceDock:
				return 'Space Dock';
			case HullSlotTypeShieldArmor:
				return 'Shield\nor\nArmor';
			case HullSlotTypeShieldElectricalMechanical:
				return 'Shield\nElect\nMech';
			case HullSlotTypeOrbitalElectrical:
				return 'Orbital\nor\nElectrical';
			case HullSlotTypeWeaponShield:
				return 'Weapon\nor\nShield';
			case HullSlotTypeScannerElectricalMechanical:
				return 'Scanner\nElec\nMech';
			case HullSlotTypeArmorScannerElectricalMechanical:
				return 'Armor\nScanner\nElec/Mech';
			case HullSlotTypeElectricalMechanical:
				return 'Elec\nor\nMech';
			case HullSlotTypeMineElectricalMechanical:
				return 'Mine\nElec\nMech';
			case HullSlotTypeGeneral:
				return 'General\nPurpose';
			default:
				return 'Unknown';
		}
	}

	const icon = (c: string | undefined) => {
		return kebabCase(c?.replace("'", '').replace(' ', '').replace('±', ''));
	};
</script>

<div
	class={`flex bg-base-300 dark:bg-base-200 tech-avatar text-sm avatar ${icon(shipDesignSlot?.hullComponent)} ${
		highlighted ? highlightedClass : ''
	}`}
	class:border={!shipDesignSlot}
	class:border-2={highlighted}
	class:border-slate-900={!shipDesignSlot && !highlighted}
	class:z-20={highlighted}
>
	<button
		type="button"
		{onclick}
		onpointerdown={(e) => {
			if (highlighted && shipDesignSlot?.hullComponent && showTooltips) {
				onTechTooltip(e, $techs.getHullComponent(shipDesignSlot.hullComponent));
			}
		}}
		class="w-full h-full"
		data-type="hull-component-button"
		data-id={shipDesignSlot?.hullComponent}
	>
		<div class="flex flex-col justify-between w-full h-full">
			{#if shipDesignSlot}
				<div class="grow">&nbsp;</div>
				<span class="h-[1rem] mt-auto text-center font-bold text-black"
					>{shipDesignSlot.quantity} of {capacity}</span
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
		onclick={() => shipDesignSlot?.quantity && shipDesignSlot.quantity++}
		data-type="increase-hull-component-button"
		data-id={shipDesignSlot?.hullComponent}
	>
		<Icon src={Plus} size="24" class="hover:stroke-accent" />
	</button>
	<button
		type="button"
		class="btn btn-sm px-1 z-30"
		onclick={() => {
			if (shipDesignSlot?.quantity != undefined) {
				shipDesignSlot.quantity--;
				if (shipDesignSlot.quantity === 0) {
					onDelete?.();
				} else {
					onUpdate?.();
				}
			}
		}}
		data-type="decrease-hull-component-button"
		data-id={shipDesignSlot?.hullComponent}
	>
		<Icon src={Minus} size="24" class="hover:stroke-accent" />
	</button>
	<button
		type="button"
		class="btn btn-sm px-1 z-30"
		onclick={onDelete}
		data-type="delete-hull-component-button"
		data-id={shipDesignSlot?.hullComponent}
	>
		<Icon src={Trash} size="24" class="hover:stroke-accent" />
	</button>
</div>
