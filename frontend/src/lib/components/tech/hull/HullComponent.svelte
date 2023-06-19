<script lang="ts">
	import { HullSlotType, type TechHullComponent } from '$lib/types/Tech';
	import { kebabCase } from 'lodash-es';
	import { $enum as eu } from 'ts-enum-util';

	export let type: HullSlotType = HullSlotType.General;
	export let capacity: number = 1;
	export let required = false;
	export let component: TechHullComponent | undefined = undefined;
	export let quantity: number | undefined = undefined;

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

	const icon = () => {
		return kebabCase(component?.name.replace("'", '').replace(' ', '').replace('±', ''));
	};
</script>

<div
	class="flex flex-col justify-between border border-slate-900 bg-base-300 tech-avatar text-sm avatar {icon()}"
>
	{#if component}
		<span class="mt-auto text-center font-bold">{quantity} of {capacity}</span>
	{:else}
		<div class="flex-grow whitespace-pre-wrap text-center">{typeDescription()}</div>
		{#if required}
			<div class="text-center text-red-500 font-bold">needs {capacity}</div>
		{:else}
			<span class="mt-auto text-center font-bold">Up to {capacity}</span>
		{/if}
	{/if}
</div>
