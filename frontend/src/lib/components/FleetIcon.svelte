<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { Fleet, ShipToken } from '$lib/types/Fleet';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import { kebabCase } from 'lodash-es';
	import { onShipDesignTooltip } from './game/tooltips/ShipDesignTooltip.svelte';
	import { None } from '$lib/types/MapObject';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { NoSymbol, QuestionMarkCircle } from '@steeze-ui/heroicons';

	const { game, player, universe } = getGameContext();

	export let fleet: Fleet;
	export let tokens: ShipToken[] = fleet.tokens ?? [];

	let icon = '';
	let design: ShipDesign | undefined;

	$: {
		icon = '';
		if (tokens && tokens.length > 0) {
			const designNum = tokens.find((token) => token.quantity > 0)?.designNum ?? None;
			const design = $universe.getDesign(fleet.playerNum, designNum);
			if (design) {
				icon = `hull-${kebabCase(design.hull)}-${design.hullSetNumber ?? 0}`;
			}
		}
	}
</script>

<div class="avatar mr-2">
	<div
		class="border-2 border-neutral p-2 bg-black"
		style={`border-color: ${$universe.getPlayerColor(fleet.playerNum)};`}
	>
		{#if tokens && tokens.reduce((count, t) => count + t.quantity, 0) > 1}
			<div class="absolute -right-2 -top-1 text-xl w-6 h-6">+</div>
		{/if}

		<div class="fleet-avatar {icon} bg-black">
			{#if !icon}
				<Icon src={NoSymbol} size="64" />
			{/if}
			<button
				type="button"
				class="w-full h-full cursor-help"
				on:pointerdown|preventDefault={(e) => onShipDesignTooltip(e, design)}
			/>
		</div>
	</div>
</div>
