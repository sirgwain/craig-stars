<script lang="ts">
	import TextTooltip, {
		type TextTooltipProps
	} from '$lib/components/game/tooltips/TextTooltip.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { showTooltip } from '$lib/services/Stores';
	import { type Minefield, type MinefieldSpec, MinefieldType } from '$lib/types/cs-proto';
	import { enumToString } from '$lib/types/Enums';
	import { ownedBy } from '$lib/types/MapObject';
	import { QuestionMarkCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import type { ChangeEventHandler } from 'svelte/elements';

	const { cs, game, player, universe, updateMinefieldOrders } = getGameContext();

	type Props = {
		minefield: Minefield;
	};

	let { minefield = $bindable() }: Props = $props();

	let spec: MinefieldSpec | undefined = $state();

	let stats = $derived($game.rules.minefieldStatsByType[minefield.minefieldType]);

	$effect(() => {
		cs.wasmService.computeMinefieldSpec({ minefield: minefield as Minefield }).then((resp) => {
			spec = resp.spec;
		});
	});

	function onTooltip(e: PointerEvent) {
		showTooltip<TextTooltipProps>(e.x, e.y, TextTooltip, {
			text: 'Numbers in parenthesis are for fleets containing a ship with ram scoop engines. Note that the chance of hitting a mine goes up the % listed for EACH warp you exceed the safe speed.'
		});
	}

	// update the minefield to detonate on the server
	const minefieldDetonateChecked: ChangeEventHandler<HTMLInputElement> = async (e) => {
		if (minefield.minefieldOrders) {
			minefield.minefieldOrders.detonate = e.currentTarget.checked;
			await updateMinefieldOrders(minefield);
		} else {
			console.error("can't detonate minefield not owned by player");
		}
	};
</script>

<div class="flex flex-row md:min-h-[11rem]">
	<div class="flex flex-col">
		<div class="avatar">
			<div class="mapobject-avatar-wrapper">
				<div
					class:standard-minefield={minefield.minefieldType === MinefieldType.STANDARD}
					class:heavy-minefield={minefield.minefieldType === MinefieldType.HEAVY}
					class:speed-bump-minefield={minefield.minefieldType === MinefieldType.SPEED_BUMP}
					class="mapobject-avatar"
				></div>
			</div>
		</div>
		<div class="text-center">{$universe.getPlayerPluralName(minefield.mapObject?.playerNum)}</div>
	</div>

	<div class="flex flex-col grow">
		<div class="flex flex-row">
			<div class="w-40">Location:</div>
			<div>
				({minefield.mapObject?.position?.x ?? 0}, {minefield.mapObject?.position?.y ?? 0})
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-40">Field Type:</div>
			<div>
				{enumToString(MinefieldType, minefield.minefieldType)}
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-40">Field Radius:</div>
			<div>
				{Math.sqrt(minefield.numMines).toFixed()} l.y. ({minefield.numMines} mines)
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-40">Maximum Safe Speed:</div>
			<div>
				Warp {stats.maxSpeed}
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-40">Chance/l.y. of a Hit:</div>
			<div>
				{(stats.chanceOfHit * 100).toFixed(2)}%
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-40">Dmg done to each ship:</div>
			<div>
				{stats.damagePerEngine} ({stats.damagePerEngineRs}) / engine
				<span class="cursor-help" onpointerdown={(e) => onTooltip(e)}>
					<Icon src={QuestionMarkCircle} size="16" class=" cursor-help inline-block" />
				</span>
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-40">Min damage done to fleet:</div>
			<div>
				{stats.minDamagePerFleet} ({stats.minDamagePerFleetRs})
				<span class="cursor-help" onpointerdown={(e) => onTooltip(e)}>
					<Icon src={QuestionMarkCircle} size="16" class=" cursor-help inline-block" />
				</span>
			</div>
		</div>
		{#if spec && ownedBy(minefield, $player.num)}
			<div class="flex flex-row">
				<div class="w-40">Decay Rate:</div>
				<div>
					{spec.decayRate} / year
				</div>
			</div>
			{#if spec.canDetonate}
				<div class="flex flex-row mt-2">
					<label>
						<input
							checked={minefield.minefieldOrders?.detonate}
							onchange={minefieldDetonateChecked}
							class="checkbox checkbox-xs"
							type="checkbox"
						/> Detonate
					</label>
				</div>
			{/if}
		{/if}
	</div>
</div>
