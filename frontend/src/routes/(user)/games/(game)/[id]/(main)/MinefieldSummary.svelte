<script lang="ts">
	import TextTooltip, {
		type TextTooltipProps
	} from '$lib/components/game/tooltips/TextTooltip.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { showTooltip } from '$lib/services/Stores';
	import type { AnyMinefield } from '$lib/services/Universe';
	import { ownedBy } from '$lib/types/MapObject';
	import { MinefieldTypeHeavy, MinefieldTypeSpeedBump, MinefieldTypeStandard } from '$lib/types/cs';
	import { QuestionMarkCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import type { ChangeEventHandler } from 'svelte/elements';

	const { game, player, universe, updateMinefieldOrders } = getGameContext();

	type Props = {
		minefield: AnyMinefield;
	};

	let { minefield = $bindable() }: Props = $props();

	let stats = $derived($game.rules.minefieldStatsByType[minefield.minefieldType]);

	function onTooltip(e: PointerEvent) {
		showTooltip<TextTooltipProps>(e.x, e.y, TextTooltip, {
			text: 'Numbers in parenthesis are for fleets containing a ship with ram scoop engines. Note that the chance of hitting a mine goes up the % listed for EACH warp you exceed the safe speed.'
		});
	}

	// update the minefield to detonate on the server
	const minefieldDetonateChecked: ChangeEventHandler<HTMLInputElement> = async (e) => {
		if ('detonate' in minefield) {
			minefield.detonate = e.currentTarget.checked;
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
					class:standard-minefield={minefield.minefieldType === MinefieldTypeStandard}
					class:heavy-minefield={minefield.minefieldType === MinefieldTypeHeavy}
					class:speed-bump-minefield={minefield.minefieldType === MinefieldTypeSpeedBump}
					class="mapobject-avatar"
				></div>
			</div>
		</div>
		<div class="text-center">{$universe.getPlayerPluralName(minefield.playerNum)}</div>
	</div>

	<div class="flex flex-col grow">
		<div class="flex flex-row">
			<div class="w-40">Location:</div>
			<div>
				({minefield.position.x}, {minefield.position.y})
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-40">Field Type:</div>
			<div>
				{minefield.minefieldType}
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-40">Field Radius:</div>
			<div>
				{minefield.spec.radius.toFixed()} l.y. ({minefield.numMines} mines)
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
				{stats.damagePerEngine} ({stats.damagePerEngineRS}) / engine
				<span class="cursor-help" onpointerdown={(e) => onTooltip(e)}>
					<Icon src={QuestionMarkCircle} size="16" class=" cursor-help inline-block" />
				</span>
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-40">Min damage done to fleet:</div>
			<div>
				{stats.minDamagePerFleet} ({stats.minDamagePerFleetRS})
				<span class="cursor-help" onpointerdown={(e) => onTooltip(e)}>
					<Icon src={QuestionMarkCircle} size="16" class=" cursor-help inline-block" />
				</span>
			</div>
		</div>
		{#if ownedBy(minefield, $player.num)}
			<div class="flex flex-row">
				<div class="w-40">Decay Rate:</div>
				<div>
					{minefield.spec.decayRate} / year
				</div>
			</div>
			{#if 'detonate' in minefield && minefield.spec.canDetonate}
				<div class="flex flex-row mt-2">
					<label>
						<input
							checked={minefield.detonate}
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
