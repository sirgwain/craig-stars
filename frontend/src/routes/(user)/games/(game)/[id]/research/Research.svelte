<script lang="ts">
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import SectionHeader from '$lib/components/SectionHeader.svelte';
	import { CommandedPlayer, TechFields } from '$lib/types/Player';
	import { Beaker } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import SpinnerNumberText from '$lib/components/SpinnerNumberText.svelte';
	import Factory from '$lib/components/icons/Factory.svelte';
	import Microscope from '$lib/components/icons/Microscope.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { enumToString } from '$lib/types/Enums';
	import { get } from '$lib/types/TechLevel';
	import {
		NextResearchField,
		TechField,
		type TechLevelJson,
		PlayerResearchSpecSchema
	} from '$lib/types/cs-proto';
	import FutureTechs from './FutureTechs.svelte';
	import { create } from '@bufbuild/protobuf';

	const { cs, player } = getGameContext();

	type Props = {
		onUpdatePlayer?: () => Promise<void>;
	};
	let { onUpdatePlayer }: Props = $props();

	let spec = $state(create(PlayerResearchSpecSchema));

	$effect(() => {
		$player.getResearchSpec(cs).then((s) => (spec = s));
	});

	const getLevel = (player: CommandedPlayer, field: TechField): number => {
		return get(player.techLevels, field);
	};

	let field: keyof TechLevelJson = $derived(
		`${$player.playerOrders.researching}`.toLowerCase() as keyof TechLevelJson
	);
	let spent = $derived($player.techLevelsSpent[field] ?? 0);

	let leftToSpend = $derived((spec.currentResearchCost ?? 0) - spent);
	let yearsLeft = $derived(Math.ceil(leftToSpend / (spec.resourcesPerYearResearchEstimated ?? 0)));
</script>

<ItemTitle>Research</ItemTitle>
<div class="stats stats-vertical sm:stats-horizontal sm:flex shadow border border-base-200 w-full">
	<div class="stat place-items-center sm:grow">
		<div class="stat-title">Researching</div>
		<div class="stat-figure"><Icon class="w-8 h-8" src={Beaker} /></div>
		<div class="stat-value">
			{enumToString(TechField, $player.playerOrders.researching)}
			{getLevel($player, $player.playerOrders.researching) + 1}
		</div>
		<div class="stat-desc pt-1">
			{spent ?? 0}/{spec.currentResearchCost} resources
			{#if yearsLeft < 100}
				, {yearsLeft.toFixed()}
				{Math.floor(yearsLeft) > 1 ? 'years' : 'year'}
			{/if}
		</div>
	</div>
	<div class="stat place-items-center sm:grow">
		<div class="stat-title">Resources Available</div>
		<div class="stat-figure"><Factory class="w-8 h-8 fill-primary" /></div>
		<div class="stat-value">
			{spec.resourcesPerYear ?? 0}
		</div>
	</div>
</div>
<div class="stats stats-vertical sm:stats-horizontal sm:flex shadow border border-base-200 w-full">
	<div class="stat place-items-center sm:grow">
		<div class="stat-title">Spent Last Year</div>
		<div class="stat-figure"><Microscope class="w-8 h-8 fill-primary" /></div>
		<div class="stat-value">
			{$player.researchSpentLastYear ?? 0}
		</div>
	</div>
	<div class="grow stat place-items-center">
		<div class="stat-title">Estimated Spending Next Year</div>
		<div class="stat-figure"><Microscope class="w-8 h-8 fill-warning" /></div>
		<div class="stat-value">
			{spec.resourcesPerYearResearchEstimated ?? 0}
		</div>
	</div>
</div>
<div class="grid grid-cols-1 md:grid-cols-2 justify-center gap-2">
	<div class="w-full">
		<SectionHeader>Tech Levels</SectionHeader>

		<SpinnerNumberText
			class="flex flex-row gap-1 place-content-center text-2xl mb-2"
			bind:value={$player.playerOrders.researchAmount}
			min={0}
			max={100}
			step={1}
			unit="%"
			onChange={onUpdatePlayer}
		>
			{#snippet begin()}
				Research Budget
			{/snippet}
			{#snippet end()}{/snippet}
		</SpinnerNumberText>

		<div class="grid grid-cols-2">
			<div class="text-center">
				Field of Study <div class="divider secondary w-[90%]"></div>
			</div>
			<div class="text-center">
				Current Level <div class="divider secondary w-[90%]"></div>
			</div>
			{#each TechFields as field (field)}
				<div class="form-control">
					<label class="label cursor-pointer">
						<span class="label-text">{enumToString(TechField, field)}</span>
						<input
							type="radio"
							name="researching"
							value={field}
							class="radio radio-sm checked:bg-primary"
							bind:group={$player.playerOrders.researching}
							onchange={onUpdatePlayer}
						/>
					</label>
				</div>
				<div class="text-center">
					{getLevel($player, field)}
				</div>
			{/each}
		</div>
		<EnumSelect
			name="nextResearchField"
			enumType={NextResearchField}
			bind:value={$player.playerOrders.nextResearchField}
			onchange={onUpdatePlayer}
		/>
	</div>

	<div class="w-full">
		<SectionHeader>Expected Research Benefits</SectionHeader>
		<FutureTechs field={$player.playerOrders.researching} />
	</div>
</div>
