<script lang="ts">
	import { onMount, onDestroy, getContext } from 'svelte';
	import { get } from 'svelte/store';
	import { clone, create } from '@bufbuild/protobuf';
	import { page } from '$app/state';
	import { replaceState } from '$app/navigation';
	import { gameKey, type GameContext } from '#lib/services/GameContext.js';
	import { techs } from '#lib/services/Stores.js';
	import { canLearnTech, type CommandedPlayer } from '#lib/types/Player.js';
	import { hasRequiredLevels } from '#lib/types/TechLevel.js';
	import { loadWasm, type CS } from '#lib/wasm.js';
	import {
		BattlePlanSchema,
		BattleSimulationSideSchema,
		ShipDesignSchema,
		type ShipDesign,
		type Fleet
	} from '#lib/types/cs-proto.js';
	import ShipDesigner from '#lib/components/game/design/ShipDesigner.svelte';
	import Tooltip from '#lib/components/game/tooltips/Tooltip.svelte';
	import Popup from '#lib/components/game/tooltips/Popup.svelte';
	import SideEditor from './SideEditor.svelte';
	import SimulationResults from './SimulationResults.svelte';
	import {
		newSides,
		newStack,
		owner,
		encodeScenario,
		decodeScenario,
		validateSides,
		defaultPlan,
		type SimulationRun
	} from './scenario';

	const context = getContext<GameContext | undefined>(gameKey);
	let sides = $state(newSides());
	let cs: CS | undefined = $state();
	let ready = $state(false);
	let freshSides = newSides();
	let storageKey = 'craig-stars:battle-simulator:v1';
	let storageMessage = $state('');
	let error = $state('');
	let activeSide = $state(0);
	let tab = $state<'setup' | 'results'>('setup');
	let iterations = $state(1);
	let seed = $state('');
	let running = $state(false);
	let completed = $state(0);
	let cancelled = false;
	const wasCancelled = () => cancelled;
	let runs: SimulationRun[] = $state([]);
	let availableDesigns: ShipDesign[] = $state([]);
	let fleets: Fleet[] = $state([]);
	let importSide: number | undefined = $state();
	let importNum = $state(0);
	let editing: { side: number; stack?: number; design?: ShipDesign } | undefined = $state();
	let selectedHull = $state('Destroyer');
	let editorError = $state('');
	let editorPlayer = $derived(editing ? owner(sides[editing.side]) : undefined);
	let hullChoices = $derived(editorPlayer ? availableHulls(editorPlayer) : []);
	let gameLabel = $state('Standalone · Custom race and tech for each side');

	onMount(() => {
		let disposed = false;
		async function init() {
			try {
				cs = context?.cs ?? (await loadWasm());
				if (disposed) return;
				if (context) {
					const player = get(context.player),
						universe = get(context.universe),
						game = get(context.game);
					storageKey = `craig-stars:battle-simulator:v1:${game.id}:${player.num}`;
					availableDesigns = universe.designs.filter((d) => d.playerNum === player.num);
					fleets = universe.fleets.filter((f) => f.mapObject?.playerNum === player.num);
					sides[0] = clone(
						BattleSimulationSideSchema,
						create(BattleSimulationSideSchema, {
							race: player.race,
							techLevels: player.techLevels,
							acquiredTechs: player.acquiredTechs,
							stacks: availableDesigns[0] ? [newStack(availableDesigns[0])] : []
						})
					);
					gameLabel = `${game.name} · Current game race and technology`;
				}
				freshSides = sides.map((s) => clone(BattleSimulationSideSchema, s));
				try {
					const saved = localStorage.getItem(storageKey);
					if (saved) sides = decodeScenario(saved);
				} catch {
					storageMessage = 'The saved scenario could not be loaded. Starting a new scenario.';
				}
				if (context) {
					const fleetNum = Number(page.url.searchParams.get('fleet'));
					const designNum = Number(page.url.searchParams.get('design'));
					if (fleetNum) {
						const fleet = fleets.find((f) => f.mapObject?.num === fleetNum);
						if (fleet) importFleet(0, fleet);
					}
					if (designNum) {
						const design = availableDesigns.find((d) => d.num === designNum);
						if (design) sides[0].stacks = [newStack(design)];
					}
				}
				if (page.url.searchParams.has('fleet') || page.url.searchParams.has('design'))
					replaceState(page.url.pathname, {});
				ready = true;
			} catch (e) {
				error = `Unable to load simulator: ${e}`;
			}
		}
		void init();
		return () => {
			disposed = true;
		};
	});
	onDestroy(() => {
		cancelled = true;
	});
	$effect(() => {
		if (!ready) return;
		try {
			const json = encodeScenario(sides);
			localStorage.setItem(storageKey, json);
			storageMessage = 'Scenario saved on this device';
		} catch {
			storageMessage = 'Unable to save this scenario on this device.';
		}
	});

	function useGameSettings(side: number) {
		if (!context) return;
		const player = get(context.player);
		const settings = clone(
			BattleSimulationSideSchema,
			create(BattleSimulationSideSchema, {
				race: player.race,
				techLevels: player.techLevels,
				acquiredTechs: player.acquiredTechs
			})
		);
		sides[side].race = settings.race;
		sides[side].techLevels = settings.techLevels;
		sides[side].acquiredTechs = settings.acquiredTechs;
	}
	function importFleet(side: number, fleet: Fleet) {
		if (!context) return;
		const universe = get(context.universe),
			player = get(context.player);
		sides[side].stacks = fleet.tokens.flatMap((token) => {
			const design = universe.getDesign(fleet.mapObject?.playerNum, token.designNum);
			if (!design) throw new Error('A fleet design is unavailable.');
			const stack = newStack(design);
			stack.quantity = token.quantity;
			stack.quantityDamaged = token.quantityDamaged;
			stack.damage = token.damage;
			stack.plan = clone(
				BattlePlanSchema,
				player.getBattlePlan(fleet.fleetOrders?.battlePlanNum ?? 0) ?? defaultPlan()
			);
			return [stack];
		});
		useGameSettings(side);
		importSide = undefined;
	}
	function availableHulls(player: CommandedPlayer) {
		return $techs.hulls.filter(
			(h) =>
				h.tech &&
				canLearnTech(player, h) &&
				hasRequiredLevels(player.techLevels, h.tech.requirements?.techLevel)
		);
	}
	function startDesign(side: number, stack?: number) {
		editorError = '';
		const allowedHulls = availableHulls(owner(sides[side]));
		if (!allowedHulls.some((h) => h.tech?.name === selectedHull))
			selectedHull = allowedHulls[0]?.tech?.name ?? '';
		editing = {
			side,
			stack,
			design:
				stack === undefined ? undefined : clone(ShipDesignSchema, sides[side].stacks[stack].design!)
		};
	}
	function chooseHull() {
		if (editing)
			editing = {
				...editing,
				design: create(ShipDesignSchema, { name: selectedHull, hull: selectedHull })
			};
	}
	async function saveDesign() {
		if (!editing?.design || !cs) return;
		editorError = '';
		const draft = editing;
		try {
			const { spec } = await cs.wasmService.computeShipDesignSpec({
				design: editing.design,
				player: editorPlayer,
				validate: true
			});
			if (editing !== draft) return;
			editing.design.spec = spec;
			const stack = newStack(editing.design);
			if (editing.stack === undefined) sides[editing.side].stacks.push(stack);
			else
				sides[editing.side].stacks[editing.stack].design = clone(ShipDesignSchema, editing.design);
			editing = undefined;
		} catch (e) {
			editorError = `${e}`;
		}
	}
	async function simulate() {
		if (!cs || running) return;
		error = '';
		try {
			validateSides(sides);
			if (seed && !/^-?\d+$/.test(seed))
				throw new Error('Seed must be a whole number, or blank for random.');
			const firstSeed = seed ? BigInt(seed) : BigInt(Math.floor(Math.random() * 2 ** 32));
			const snapshot = sides.map((s) => clone(BattleSimulationSideSchema, s));
			running = true;
			completed = 0;
			cancelled = false;
			const results: SimulationRun[] = [];
			for (let i = 0; i < iterations && !wasCancelled(); i++) {
				// Yield between battles so mobile browsers can paint progress and accept Cancel.
				await new Promise((resolve) => setTimeout(resolve, 0));
				if (wasCancelled()) break;
				const runSeed = firstSeed + BigInt(i);
				const response = await cs.wasmService.simulateBattle({ sides: snapshot, seed: runSeed });
				results.push({ response, sides: snapshot, seed: runSeed });
				completed = results.length;
			}
			if (results.length) {
				runs = results;
				tab = 'results';
			}
		} catch (e) {
			error = `${e}`;
		} finally {
			running = false;
		}
	}
	function reset() {
		sides = freshSides.map((s) => clone(BattleSimulationSideSchema, s));
		runs = [];
		error = '';
		tab = 'setup';
	}
</script>

<div class="mx-auto w-full max-w-7xl min-w-0 p-2 md:p-4">
	<div class="mb-3 flex flex-wrap items-center justify-between gap-2">
		<div class="min-w-0 flex-1">
			<h1 class="text-2xl font-semibold">Battle Simulator</h1>
			<p class="mt-1 truncate text-secondary" title={gameLabel}>{gameLabel}</p>
		</div>
		<button class="btn btn-sm btn-outline btn-secondary" disabled={running} onclick={reset}
			>New scenario</button
		>
	</div>
	<div
		role="tablist"
		aria-label="Simulator views"
		class="tabs tabs-border mb-3 border-b border-base-300"
	>
		<button
			role="tab"
			aria-selected={tab === 'setup'}
			class="tab"
			class:tab-active={tab === 'setup'}
			onclick={() => (tab = 'setup')}>Setup</button
		><button
			role="tab"
			aria-selected={tab === 'results'}
			class="tab"
			class:tab-active={tab === 'results'}
			disabled={!runs.length}
			onclick={() => (tab = 'results')}>Results</button
		>
	</div>
	{#if error}<div role="alert" class="alert alert-error mb-3">{error}</div>{/if}
	{#if !ready}<div class="p-4" role="status">Loading simulator…</div>
	{:else if tab === 'setup'}
		<fieldset disabled={running}>
			<div class="join mb-3 flex md:hidden" role="group" aria-label="Choose side">
				{#each [0, 1] as i (i)}<button
						class="join-item btn btn-sm grow"
						class:btn-secondary={activeSide === i}
						onclick={() => (activeSide = i)}>{i === 0 ? 'A · Your forces' : 'B · Opponent'}</button
					>{/each}
			</div>
			<div class="grid gap-3 md:grid-cols-2">
				{#each sides as _side, i (i)}<div class:hidden={activeSide !== i} class="min-w-0 md:block">
						<SideEditor
							bind:side={sides[i]}
							index={i}
							designs={availableDesigns}
							hasGame={!!context}
							onEdit={(stack) => startDesign(i, stack)}
							onNew={() => startDesign(i)}
							onUseGame={() => useGameSettings(i)}
							onImport={() => {
								importSide = i;
								importNum = 0;
							}}
						/>
					</div>{/each}
			</div>
		</fieldset>
	{:else if runs.length}{#key runs}<SimulationResults {runs} />{/key}{/if}
	{#if ready}
		<div
			class:sticky={tab === 'setup'}
			class="bottom-0 z-10 mt-4 md:static flex flex-wrap items-end gap-3 border-t border-secondary bg-base-100 py-3"
		>
			<label class="flex flex-col gap-1"
				>Battles<select
					aria-label="Number of battles"
					class="select select-sm"
					bind:value={iterations}
					disabled={running}
					><option value={1}>1 battle</option><option value={10}>10 battles</option><option
						value={100}>100 battles</option
					></select
				></label
			>
			<label class="flex min-w-0 flex-col gap-1"
				>Seed (optional)<input
					aria-label="Simulation seed"
					class="input input-sm w-32"
					placeholder="Random"
					bind:value={seed}
					disabled={running}
				/></label
			>
			{#if running}<span role="status">Simulating {completed} / {iterations}</span><button
					class="btn btn-sm btn-outline btn-secondary"
					onclick={() => (cancelled = true)}>Cancel</button
				>{:else}<button class="btn btn-sm btn-secondary" onclick={simulate}
					>{tab === 'results' ? 'Run again' : 'Simulate'}</button
				>{/if}
			<span class="text-secondary">{storageMessage}</span>
		</div>
	{/if}
</div>
{#if importSide !== undefined}
	<div class="modal modal-open" role="dialog" aria-modal="true" aria-label="Import fleet">
		<div class="modal-box">
			<h2 class="text-xl font-semibold">Import fleet as stacks</h2>
			<p class="my-2">
				Replaces this side’s stacks and copies the fleet’s battle plan, race and tech. Game orders
				stay unchanged.
			</p>
			<select aria-label="Fleet to import" class="select w-full" bind:value={importNum}
				>{#each fleets as fleet, i (fleet.mapObject?.num)}<option value={i}
						>{fleet.mapObject?.name || fleet.baseName} #{fleet.mapObject?.num}</option
					>{/each}</select
			>
			<div class="modal-action">
				<button
					class="btn btn-sm btn-outline btn-secondary"
					onclick={() => (importSide = undefined)}>Cancel</button
				><button
					class="btn btn-sm btn-secondary"
					disabled={!fleets[importNum]}
					onclick={() => {
						try {
							importFleet(importSide!, fleets[importNum]);
						} catch (e) {
							error = `${e}`;
							importSide = undefined;
						}
					}}>Import stacks</button
				>
			</div>
		</div>
	</div>
{/if}
{#if editing && cs}
	<div
		class="modal modal-open"
		role="dialog"
		aria-modal="true"
		aria-label="Simulation ship designer"
	>
		<div class="modal-box w-full max-w-5xl">
			<div class="mb-3 flex items-center justify-between gap-2">
				<h2 class="text-xl font-semibold">
					{editing.stack === undefined ? 'New simulation design' : 'Edit simulation design'}
				</h2>
				<button class="btn btn-sm btn-outline btn-secondary" onclick={() => (editing = undefined)}
					>Cancel</button
				>
			</div>
			{#if editorError}<div class="alert alert-error mb-3" role="alert">{editorError}</div>{/if}
			<p class="mb-3 text-secondary">
				Uses side {editing.side === 0 ? 'A' : 'B'}’s race and tech. Saved only in this scenario.
			</p>
			{#if !editing.design}<label class="flex flex-col gap-1"
					>Hull<select aria-label="Simulation hull" class="select w-full" bind:value={selectedHull}
						>{#each hullChoices as hull (hull.tech?.name)}<option value={hull.tech?.name}
								>{hull.tech?.name}</option
							>{/each}</select
					></label
				><button
					class="btn btn-secondary mt-3"
					disabled={!hullChoices.some((h) => h.tech?.name === selectedHull)}
					onclick={chooseHull}>Design ship</button
				>
			{:else if $techs.getHull(editing.design.hull)}<ShipDesigner
					bind:design={editing.design}
					hull={$techs.getHull(editing.design.hull)!}
					client={cs}
					playerOverride={editorPlayer}
					error=""
					onSave={saveDesign}
				/>
				<div class="mt-3 flex justify-end">
					<button class="btn btn-secondary" onclick={saveDesign}>Save simulation design</button>
				</div>{/if}
		</div>
	</div>
{/if}
{#if !context}<Tooltip /><Popup />{/if}
