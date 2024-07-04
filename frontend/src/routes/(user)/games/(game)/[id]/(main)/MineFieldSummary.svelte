<script lang="ts">
	import TextTooltip, {
		type TextTooltipProps
	} from '$lib/components/game/tooltips/TextTooltip.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { showTooltip } from '$lib/services/Stores';
	import { ownedBy } from '$lib/types/MapObject';
	import { MineFieldTypes, type MineField } from '$lib/types/MineField';
	import { QuestionMarkCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	const { game, player, universe } = getGameContext();

	export let mineField: MineField;

	$: stats = $game.rules.mineFieldStatsByType[mineField.mineFieldType];

	function onTooltip(e: PointerEvent) {
		showTooltip<TextTooltipProps>(e.x, e.y, TextTooltip, {
			text: 'Numbers in parenthesis are for fleets containing a ship with ram scoop engines. Note that the chance of hitting a mine goes up the % listed for EACH warp you exceed the safe speed.'
		});
	}
</script>

<div class="flex flex-row min-h-[11rem]">
	<div class="flex flex-col">
		<div class="avatar ">
			<div class="border-2 border-neutral mr-2 p-2 bg-black">
				<div
					class:standard-mine-field={mineField.mineFieldType === MineFieldTypes.Standard}
					class:heavy-mine-field={mineField.mineFieldType === MineFieldTypes.Heavy}
					class:speed-bump-mine-field={mineField.mineFieldType === MineFieldTypes.SpeedBump}
					class="mapobject-avatar bg-black"
				/>
			</div>
		</div>
		<div class="text-center">{$universe.getPlayerName(mineField.playerNum)}</div>
	</div>

	<div class="flex flex-col grow">
		<div class="flex flex-row">
			<div class="w-40">Location:</div>
			<div>
				({mineField.position.x}, {mineField.position.y})
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-40">Field Type:</div>
			<div>
				{mineField.mineFieldType}
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-40">Field Radius:</div>
			<div>
				{mineField.spec.radius.toFixed()} l.y. ({mineField.numMines} mines)
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
				<span class="cursor-help" on:pointerdown|preventDefault={(e) => onTooltip(e)}>
					<Icon src={QuestionMarkCircle} size="16" class=" cursor-help inline-block" />
				</span>
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-40">Min damage done to fleet:</div>
			<div>
				{stats.minDamagePerFleet} ({stats.minDamagePerFleetRS})
				<span class="cursor-help" on:pointerdown|preventDefault={(e) => onTooltip(e)}>
					<Icon src={QuestionMarkCircle} size="16" class=" cursor-help inline-block" />
				</span>
			</div>
		</div>
		{#if ownedBy(mineField, $player.num)}
			<div class="flex flex-row">
				<div class="w-40">Decay Rate:</div>
				<div>
					{mineField.spec.decayRate} / year
				</div>
			</div>
		{/if}
	</div>
</div>
