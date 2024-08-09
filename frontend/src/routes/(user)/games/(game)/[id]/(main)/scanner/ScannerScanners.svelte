<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';

	import { positionKey } from '$lib/types/MapObject';
	import { NoScanner } from '$lib/types/Tech';
	import type { Vector } from '$lib/types/Vector';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';

	const { game, player, universe, settings } = getGameContext();
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');

	type Scanner = {
		position: Vector;
		scanRange: number;
		scanRangePen: number;
	};

	let scanners: Scanner[] = [];
	$: {
		if ($data) {
			const scannersByPosition = new Map<string, Scanner>();

			$universe.planets
				.filter((p) => p.playerNum == $player.num && p.spec?.scanner)
				.forEach((planet) =>
					scannersByPosition.set(positionKey(planet), {
						position: planet.position,
						scanRange: planet.spec?.scanRange ?? 0,
						scanRangePen: planet.spec?.scanRangePen ?? 0
					})
				);

			$universe.fleets
				.filter(
					(fleet) =>
						fleet.playerNum == $player.num &&
						((fleet.spec?.scanRange ?? 0) > 0 || (fleet.spec?.scanRangePen ?? 0) > 0)
				)
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

			$universe.mineralPackets
				.filter(
					(packet) =>
						packet.playerNum == $player.num &&
						(packet.scanRange != NoScanner || packet.scanRangePen != NoScanner)
				)
				.forEach((packet) => {
					const key = positionKey(packet);
					const scanner = {
						position: packet.position,
						scanRange: packet.scanRange,
						scanRangePen: packet.scanRangePen
					};
					const existing = scannersByPosition.get(key);
					if (existing) {
						existing.scanRange = Math.max(existing.scanRange, scanner.scanRange);
						existing.scanRangePen = Math.max(existing.scanRangePen, scanner.scanRangePen);
					} else {
						scannersByPosition.set(key, scanner);
					}
				});

			if ($settings.showAllyScanners) {
				$universe.planets
					.filter((p) => $player.isSharingMap(p.playerNum) && p.spec?.scanner)
					.forEach((planet) =>
						scannersByPosition.set(positionKey(planet), {
							position: planet.position,
							scanRange: planet.spec?.scanRange ?? 0,
							scanRangePen: planet.spec?.scanRangePen ?? 0
						})
					);

				// find ally's scanners
				$universe.fleets
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

				$universe.mineralPackets
					.filter(
						(packet) =>
							$player.isSharingMap(packet.playerNum) &&
							(packet.scanRange != NoScanner || packet.scanRangePen != NoScanner)
					)
					.forEach((packet) => {
						const key = positionKey(packet);
						const scanner = {
							position: packet.position,
							scanRange: packet.scanRange,
							scanRangePen: packet.scanRangePen
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

			scanners = [];
			scanners = Array.from(scannersByPosition.values());
		}
	}

	$: scannerScale = $settings.scannerPercent / 100.0;
</script>

{#if $settings.showScanners}
	{#each scanners as scanner}
		<circle
			cx={$xGet(scanner)}
			cy={$yGet(scanner)}
			r={$xScale(scanner.scanRange * scannerScale)}
			class="scanner"
		/>
	{/each}
	{#each scanners as scanner}
		{#if scanner.scanRangePen > 0}
			<circle
				cx={$xGet(scanner)}
				cy={$yGet(scanner)}
				r={$xScale(scanner.scanRangePen * scannerScale)}
				class="scanner-pen"
			/>
		{/if}
	{/each}
{/if}
