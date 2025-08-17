<script lang="ts">
	import MineralMini from '$lib/components/game/MineralMini.svelte';
	import type { OnCancel, OnOk } from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import type { AnyFleet } from '$lib/services/Universe';
	import { population } from '$lib/types/Cargo';
	import { ReportAgeUnexplored } from '$lib/types/Consts';
	import type { MysteryTraderIntel, PlanetIntel } from '$lib/types/cs-proto';
	import { getMapObjectName, key, owned, ownedBy, type MapObjectLike } from '$lib/types/MapObject';
	import { None } from '$lib/types/Consts';
	import { onMount } from 'svelte';

	const { player, universe, settings } = getGameContext();

	type Props = {
		maxPlanetResults?: number;
		maxFleetResults?: number;
		maxMiscResults?: number;
		onOk?: OnOk<MapObjectLike | undefined>;
		onCancel?: OnCancel;
	};

	let {
		maxPlanetResults = 10,
		maxFleetResults = 10,
		maxMiscResults = 10,
		onOk,
		onCancel
	}: Props = $props();

	type Results = {
		planets: PlanetIntel[];
		fleets: AnyFleet[];
		mysteryTraders: MysteryTraderIntel[];
	};

	function getResults(search: string): Results {
		if (search == '') {
			return {
				planets: [],
				fleets: [],
				mysteryTraders: []
			};
		}
		const terms = search.split(' ');

		const planets = $universe.getPlanets($settings.sortPlanetsKey, $settings.sortPlanetsDescending);
		const fleets = $universe.getAllFleets($settings.sortFleetsKey, $settings.sortFleetsDescending);
		const mysteryTraders = $universe.mysteryTraderIntels;

		// return true if a mapboject name or player matches a search term
		const termSearch = (term: string, mo: MapObjectLike): boolean =>
			mo.mapObject?.name.toLowerCase().indexOf(term.toLowerCase()) != -1 ||
			(mo.mapObject?.playerNum != None &&
				$universe
					.getPlayerPluralName(mo.mapObject?.playerNum)
					.toLowerCase()
					.indexOf(term.toLowerCase()) != -1);

		return {
			planets:
				planets
					.filter((i) => terms.every((term) => termSearch(term, i)))
					.slice(0, maxPlanetResults) ?? [],
			fleets:
				fleets
					.filter((i) => terms.every((term) => termSearch(term, i)))
					.slice(0, maxFleetResults) ?? [],

			mysteryTraders:
				mysteryTraders
					.filter((i) => terms.every((term) => termSearch(term, i)))
					.slice(0, maxMiscResults) ?? []
		};
	}

	function ok() {
		onOk?.(selectedItem);
	}

	function selectPrevious() {
		selectedItemIndex = Math.max(0, selectedItemIndex - 1);
	}

	function selectNext() {
		selectedItemIndex = Math.min(
			results.planets.length + results.fleets.length + results.mysteryTraders.length,
			selectedItemIndex + 1
		);
	}

	function onSearchKeyDown(event: KeyboardEvent) {
		switch (event.key) {
			case 'ArrowDown':
				// Do something for "down arrow" key press.
				selectNext();
				event.preventDefault();
				break;
			case 'ArrowUp':
				// Do something for "up arrow" key press.
				selectPrevious();
				event.preventDefault();
				break;
			case 'Enter':
				onOk?.(selectedItem);
				event.preventDefault();
				break;
			case 'Escape':
				if ($settings.searchQuery != '') {
					$settings.searchQuery = '';
				} else {
					onCancel?.();
					event.preventDefault();
				}
				break;
		}
	}

	let searchInput: HTMLInputElement | undefined = $state();
	onMount(() => {
		searchInput?.focus();
		selectedItemIndex = 0;
	});
	// the currently selected item
	let selectedItemIndex = $state(0);

	// when search chnages, update our search results
	let results = $derived(getResults($settings.searchQuery));
	let selectedItem = $derived(
		selectedItemIndex < results.planets.length
			? results.planets[selectedItemIndex]
			: selectedItemIndex < results.planets.length + results.fleets.length
				? results.fleets[selectedItemIndex - results.planets.length]
				: selectedItemIndex <
					  results.planets.length + results.fleets.length + results.mysteryTraders.length
					? results.mysteryTraders[
							selectedItemIndex - results.planets.length + results.fleets.length
						]
					: undefined
	);
</script>

<div class="flex flex-col gap-1 h-full pb-2">
	<input
		type="search"
		name="search"
		placeholder="search"
		class="input input-bordered input-sm sm:w-auto mt-1"
		autocomplete="off"
		autocorrect="off"
		autocapitalize="off"
		spellcheck="false"
		bind:this={searchInput}
		bind:value={$settings.searchQuery}
		onkeydown={onSearchKeyDown}
		onfocus={() => searchInput?.select()}
	/>
	<div class="h-full">
		<div class="mt-2 w-full h-full bg-base-200 border-2 border-base-300 overflow-y-auto pl-2">
			{#if results.planets.length > 0}
				<h3 class="text-2xl font-bold mb-1">Planets</h3>
				<ul class="mx-1">
					{#each results.planets as planet, index (planet.mapObject?.num)}
						<!-- svelte-ignore a11y_mouse_events_have_key_events -->
						<li
							class="rounded-lg px-2"
							class:bg-primary={selectedItemIndex == index}
							onmouseover={() => (selectedItemIndex = index)}
						>
							<button class="text-xl text-left w-full" onclick={ok}>
								<div class="flex flex-row gap-1">
									{#if planet.mapObject?.playerNum != None}
										<span style={`color: ${$universe.getPlayerColor(planet.mapObject?.playerNum)}`}
											>{$universe.getPlayerPluralName(planet.mapObject?.playerNum)}</span
										>
										{planet.mapObject?.name}
									{:else}
										{planet.mapObject?.name}
									{/if}
									{#if 'reportAge' in planet && planet.reportAge !== ReportAgeUnexplored}
										{#if owned(planet)}
											<div>-</div>
											<div class="text-base my-auto">
												{population(planet.cargo)
													? population(planet.cargo).toLocaleString() + ' pop'
													: ''}
											</div>
										{/if}
										<div>-</div>
										<div class="text-base my-auto">
											{#if planet.spec?.canTerraform}
												<span
													class:text-habitable={(planet.spec?.habitability ?? 0) > 0}
													class:text-uninhabitable={(planet.spec?.habitability ?? 0) < 0}
													>{planet.spec?.habitability ?? 0}%</span
												>
												/
												<span class="text-terraformable"
													>{planet.spec?.terraformedHabitability ?? 0}%</span
												>
											{:else}
												<span
													class:text-habitable={(planet.spec?.habitability ?? 0) > 0}
													class:text-uninhabitable={(planet.spec?.habitability ?? 0) < 0}
												>
													{planet.spec?.habitability ?? 0}%</span
												>
											{/if}
										</div>
										{#if ownedBy(planet, $player.num)}
											<div>-</div>
											<div class="text-base my-auto">
												{planet.spec?.resourcesPerYear
													? planet.spec?.resourcesPerYear.toLocaleString() + ' res'
													: ''}
											</div>
											<div>-</div>
											<div class="text-base my-auto">
												<MineralMini mineral={planet.cargo} />
											</div>
										{/if}
									{/if}
								</div>
							</button>
						</li>
					{/each}
				</ul>
			{/if}
			{#if results.fleets.length > 0}
				<h3 class="text-2xl font-bold mb-1">Fleets</h3>
				<ul class="mx-1">
					{#each results.fleets as fleet, index (key(fleet))}
						<!-- svelte-ignore a11y_mouse_events_have_key_events -->
						<li
							class="rounded-lg px-2"
							class:bg-primary={selectedItemIndex == results.planets.length + index}
							onmouseover={() => (selectedItemIndex = results.planets.length + index)}
						>
							<button class="text-xl text-left w-full" onclick={ok}>
								<span style={`color: ${$universe.getPlayerColor(fleet.mapObject?.playerNum)}`}
									>{$universe.getPlayerPluralName(fleet.mapObject?.playerNum)}</span
								>
								{getMapObjectName(fleet)}
							</button>
						</li>
					{/each}
				</ul>
			{/if}
			{#if results.mysteryTraders.length > 0}
				<h3 class="text-2xl font-bold mb-1">Mystery Traders</h3>
				<ul class="mx-1">
					{#each results.mysteryTraders as mysterytrader, index (mysterytrader.mapObject?.num)}
						<!-- svelte-ignore a11y_mouse_events_have_key_events -->
						<li
							class="rounded-lg px-2"
							class:bg-primary={selectedItemIndex ==
								results.planets.length + results.fleets.length + index}
							onmouseover={() =>
								(selectedItemIndex = results.planets.length + results.fleets.length + index)}
						>
							<button class="text-xl text-left w-full" onclick={ok}>
								<span class="text-mystery-trader"> {mysterytrader.mapObject?.name}</span></button
							>
						</li>
					{/each}
				</ul>
			{/if}
		</div>
	</div>
</div>
