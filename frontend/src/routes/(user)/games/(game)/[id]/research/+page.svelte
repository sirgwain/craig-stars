<script lang="ts">
	import { Icon } from '@steeze-ui/svelte-icon';
	import { Beaker } from '@steeze-ui/heroicons';
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import SectionHeader from '$lib/components/SectionHeader.svelte';
	import { player } from '$lib/services/Context';
	import { PlayerService } from '$lib/services/PlayerService';
	import { NextResearchField, TechField, type TechLevel } from '$lib/types/Player';

	import { startCase } from 'lodash-es';
	import { $enum as eu } from 'ts-enum-util';
	import NumberInput from '$lib/components/NumberInput.svelte';
	import Factory from '$lib/components/icons/Factory.svelte';
	import Microscope from '$lib/components/icons/Microscope.svelte';

	const getLevel = (field: TechField | string): number => {
		const f: keyof TechLevel = `${field}`.toLowerCase() as keyof TechLevel;
		return $player?.techLevels[f] ?? 0;
	};

	const updatePlayerOrders = async () => {
		if ($player) {
			const result = await PlayerService.updateOrders($player);
			Object.assign($player, result?.player);

			if (result?.planets) {
				$player.planets = result?.planets;
			}
		}
	};

	let spent = 0;
	$: {
		if ($player) {
			const field: keyof TechLevel = `${$player.researching}`.toLowerCase() as keyof TechLevel;
			spent = $player.techLevelsSpent[field] ?? 0;
		}
	}

	$: leftToSpend = ($player?.spec.currentResearchCost ?? 0) - spent;
	$: yearsLeft = leftToSpend / ($player?.spec.resourcesPerYearResearch ?? 0);
</script>

{#if $player}
	<div class="w-full mx-auto md:max-w-2xl">
		<ItemTitle>Research</ItemTitle>
		<div
			class="stats stats-vertical sm:stats-horizontal sm:flex shadow border border-base-200 w-full"
		>
			<div class="stat place-items-center sm:grow">
				<div class="stat-title">Researching</div>
				<div class="stat-figure"><Icon class="w-8 h-8" src={Beaker} /></div>
				<div class="stat-value">
					{$player.researching}
					{getLevel($player.researching) + 1}
				</div>
				<div class="stat-desc pt-1">
					{spent ?? 0}/{$player?.spec.currentResearchCost} resources, {yearsLeft.toFixed()} years
				</div>
			</div>
			<div class="stat place-items-center sm:grow">
				<div class="stat-title">Resources Available</div>
				<div class="stat-figure"><Factory class="w-8 h-8" /></div>
				<div class="stat-value">
					{$player.spec.resourcesPerYear}
				</div>
			</div>
		</div>
		<div
			class="stats stats-vertical sm:stats-horizontal sm:flex shadow border border-base-200 w-full"
		>
			<div class="stat place-items-center sm:grow">
				<div class="stat-title">Spent Last Year</div>
				<div class="stat-figure"><Microscope class="w-8 h-8" /></div>
				<div class="stat-value">
					{$player.researchSpentLastYear ?? 0}
				</div>
			</div>
			<div class="grow stat place-items-center">
				<div class="stat-title">Spending Next Year</div>
				<div class="stat-figure"><Microscope class="w-8 h-8" /></div>
				<div class="stat-value">
					{$player.spec.resourcesPerYearResearch ?? 0}
				</div>
			</div>
		</div>
		<div class="grid grid-cols-1 md:grid-cols-2 justify-center gap-2">
			<div class="w-full">
				<SectionHeader>Tech Levels</SectionHeader>

				<NumberInput
					name="researchAmount"
					title="Research Budget"
					bind:value={$player.researchAmount}
					min={0}
					max={100}
					step={1}
					unit="%"
					on:change={updatePlayerOrders}
				/>

				<div class="grid grid-cols-2">
					<div class="text-center">
						Field of Study <div class="divider secondary w-[90%]" />
					</div>
					<div class="text-center">
						Current Level <div class="divider secondary w-[90%]" />
					</div>
					{#each eu(TechField).getKeys() as field}
						<div class="form-control">
							<label class="label cursor-pointer">
								<span class="label-text">{startCase(field.toString())}</span>
								<input
									type="radio"
									name="researching"
									value={field}
									class="radio radio-sm checked:bg-primary"
									bind:group={$player.researching}
									on:change={updatePlayerOrders}
								/>
							</label>
						</div>
						<div class="text-center">
							{getLevel(field)}
						</div>
					{/each}
				</div>
				<EnumSelect
					name="nextResearchField"
					enumType={NextResearchField}
					bind:value={$player.nextResearchField}
					on:change={updatePlayerOrders}
				/>
			</div>

			<div class="w-full">
				<SectionHeader>Expected Research Benefits</SectionHeader>
				<div class="grid grid-cols-2" />
			</div>
		</div>
	</div>
{/if}
