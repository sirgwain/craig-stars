<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';

	import { NoScanner } from '$lib/types/cs';
	import { positionKey } from '$lib/types/MapObject';
	import type { Vector } from '$lib/types/cs';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import { SvelteMap } from 'svelte/reactivity';

	const { player, universe, settings } = getGameContext();
	const { xGet, yGet, xScale } = getContext<LayerCake>('LayerCake');

	type Scanner = {
		position: Vector;
		scanRange: number;
		scanRangePen: number;
	};

	let scanners = $derived.by(() => {
		const scannersByPosition = new SvelteMap<string, Scanner>();

		if ($settings.showScanners) {
			$universe.planets
				.filter((p) => p.spec?.scanner)
				.forEach((planet) =>
					scannersByPosition.set(positionKey(planet), {
						position: planet.position,
						scanRange: planet.spec?.scanRange ?? 0,
						scanRangePen: planet.spec?.scanRangePen ?? 0
					})
				);

			$universe.fleets
				.filter((fleet) => (fleet.spec?.scanRange ?? 0) > 0 || (fleet.spec?.scanRangePen ?? 0) > 0)
				.forEach((fleet) => {
					const key = positionKey(fleet);
					const scanner = {
						position: fleet.position,
						scanRange: fleet.spec?.scanRange ?? 0,
						scanRangePen: fleet.spec?.scanRangePen ?? 0
					};
					const existing = scannersByPosition.get(key);
					if (existing) {
						existing.scanRange = Math.max(existing.scanRange, scanner.scanRange);
						existing.scanRangePen = Math.max(existing.scanRangePen, scanner.scanRangePen);
					} else {
						scannersByPosition.set(key, scanner);
					}
				});

			$universe.mineralPacketIntels
				.filter(
					(packet) =>
						packet.playerNum == $player.num &&
						(packet.scanRange != NoScanner || packet.scanRangePen != NoScanner)
				)
				.forEach((packet) => {
					const key = positionKey(packet);
					const scanner = {
						position: packet.position,
						scanRange: packet.scanRange ?? 0,
						scanRangePen: packet.scanRangePen ?? 0
					};
					const existing = scannersByPosition.get(key);
					if (existing) {
						existing.scanRange = Math.max(existing.scanRange, scanner.scanRange);
						existing.scanRangePen = Math.max(existing.scanRangePen, scanner.scanRangePen);
					} else {
						scannersByPosition.set(key, scanner);
					}
				});
		}
		if ($settings.showAllyScanners) {
			$universe.planetIntels
				.filter((p) => $player.isSharingMap(p.playerNum) && p.spec?.scanner)
				.forEach((planet) =>
					scannersByPosition.set(positionKey(planet), {
						position: planet.position,
						scanRange: planet.spec?.scanRange ?? 0,
						scanRangePen: planet.spec?.scanRangePen ?? 0
					})
				);

			// find ally's scanners
			$universe.fleetIntels
				.filter(
					(fleet) =>
						$player.isSharingMap(fleet.playerNum) &&
						((fleet.scanRange ?? 0) > 0 || (fleet.scanRangePen ?? 0) > 0)
				)
				.forEach((fleet) => {
					const key = positionKey(fleet);
					const scanner = {
						position: fleet.position,
						scanRange: fleet.scanRange ?? 0,
						scanRangePen: fleet.scanRangePen ?? 0
					};
					const existing = scannersByPosition.get(key);
					if (existing) {
						existing.scanRange = Math.max(existing.scanRange, scanner.scanRange);
						existing.scanRangePen = Math.max(existing.scanRangePen, scanner.scanRangePen);
					} else {
						scannersByPosition.set(key, scanner);
					}
				});

			$universe.mineralPacketIntels
				.filter(
					(packet) =>
						$player.isSharingMap(packet.playerNum) &&
						(packet.scanRange != NoScanner || packet.scanRangePen != NoScanner)
				)
				.forEach((packet) => {
					const key = positionKey(packet);
					const scanner = {
						position: packet.position,
						scanRange: packet.scanRange ?? 0,
						scanRangePen: packet.scanRangePen ?? 0
					};
					const existing = scannersByPosition.get(key);
					if (existing) {
						existing.scanRange = Math.max(existing.scanRange, scanner.scanRange);
						existing.scanRangePen = Math.max(existing.scanRangePen, scanner.scanRangePen);
					} else {
						scannersByPosition.set(key, scanner);
					}
				});
		}
		return Array.from(scannersByPosition.values());
	});

	let scannerScale = $derived($settings.scannerPercent / 100.0);
</script>

{#each scanners as scanner (scanner)}
	<circle
		cx={$xGet(scanner)}
		cy={$yGet(scanner)}
		r={$xScale(scanner.scanRange * scannerScale)}
		class="scanner"
	/>
{/each}
{#each scanners as scanner (scanner)}
	{#if scanner.scanRangePen > 0}
		<circle
			cx={$xGet(scanner)}
			cy={$yGet(scanner)}
			r={$xScale(scanner.scanRangePen * scannerScale)}
			class="scanner-pen"
		/>
	{/if}
{/each}
