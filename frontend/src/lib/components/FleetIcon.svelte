<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { getHullIcon } from '$lib/techicon';
	import type { ShipToken } from '$lib/types/cs';
	import type { Fleet } from '$lib/types/cs';
	import type { ShipDesign } from '$lib/types/cs';
	import { NoSymbol } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { onShipDesignTooltip } from './game/tooltips/ShipDesignTooltip.svelte';

	const { universe } = getGameContext();

	type Props = {
		fleet: Fleet;
		tokens?: ShipToken[];
	};

	let { fleet, tokens = fleet.tokens ?? [] }: Props = $props();

	const design: ShipDesign | undefined = $derived.by(() => {
		if (fleet.tokens && fleet.tokens.length > 0) {
			const designNum = fleet.tokens[0].designNum;
			return $universe.getDesign(fleet.playerNum, designNum);
		}
	});
</script>

<div class="avatar mr-2">
	<div
		class="border-2 border-neutral p-2 bg-black"
		style={`border-color: ${$universe.getPlayerColor(fleet.playerNum)};`}
	>
		{#if tokens && tokens.reduce((count, t) => count + t.quantity, 0) > 1}
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
