<script lang="ts">
	
	import { cargoPercent, totalCargo, type Cargo } from '$lib/types/Cargo';

	type Props = {
		value?: Cargo;
		capacity?: number | undefined;
		canTransferCargo?: boolean;
		onpointerdown?: (e: PointerEvent) => void | undefined;
	};

	let {
		value = {
			ironium: 0,
			boranium: 0,
			germanium: 0,
			colonists: 0
		},
		capacity = 0,
		canTransferCargo = false,
		onpointerdown
	}: Props = $props();

	let percent: Cargo = $derived(cargoPercent(value, capacity));
</script>

<div
	onpointerdown={(e) => (canTransferCargo && onpointerdown ? onpointerdown(e) : undefined)}
	class="border border-secondary h-[1rem] text-[0rem] relative bg-gauge select-none"
	class:cursor-pointer={canTransferCargo}
>
	<div
		class="font-semibold text-sm text-center align-middle text-white mix-blend-difference w-full bg-blend-difference absolute"
	>
		{totalCargo(value)} of {capacity ?? 0}kT
	</div>
	<div
		style={`left: 0%; width: ${percent.ironium?.toFixed()}%`}
		class="ironium-bar h-full inline-block"
	></div>
	<div
		style={`width: ${percent.boranium?.toFixed()}%`}
		class="boranium-bar h-full inline-block"
	></div>
	<div
		style={`width: ${percent.germanium?.toFixed()}%`}
		class="germanium-bar h-full inline-block"
	></div>
	<div
		style={`width: ${percent.colonists?.toFixed()}%`}
		class="colonists-bar h-full inline-block"
	></div>
</div>
