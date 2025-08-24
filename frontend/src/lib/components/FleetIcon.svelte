<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { getHullIcon } from '$lib/techicon';
	import type { Fleet, ShipDesign, ShipToken } from '$lib/types/cs-proto';
	import { NoSymbol } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { onShipDesignTooltip } from './game/tooltips/ShipDesignTooltip.svelte';

	const { universe } = getGameContext();

	type Props = {
		fleet: Fleet;
		tokens?: ShipToken[];
	};

	let { fleet, tokens = fleet.tokens }: Props = $props();

	const design: ShipDesign | undefined = $derived.by(() => {
		const token = tokens.find((t) => t.quantity > 0);
		if (token) {
			return $universe.getDesign(fleet.mapObject?.playerNum, token.designNum);
		}
	});
</script>

<div class="avatar mr-2">
	<div
		class="border-2 border-neutral p-2 bg-black"
		style={`border-color: ${$universe.getPlayerColor(fleet.mapObject?.playerNum)};`}
	>
		{#if tokens.reduce((count, t) => count + t.quantity, 0) > 1}
			<div class="absolute -right-2 -top-1 text-xl w-6 h-6">+</div>
		{/if}

		<div class="fleet-avatar {getHullIcon(design)} bg-black">
			{#if !design}
				<Icon src={NoSymbol} size="64" />
			{/if}
			<button
				type="button"
				class="w-full h-full cursor-help"
				aria-label="Opens ship design tooltip"
				onpointerdown={(e) => onShipDesignTooltip(e, design)}
			></button>
		</div>
	</div>
</div>
