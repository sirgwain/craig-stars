<script lang="ts">
	import type { FullGame } from '$lib/services/FullGame';
	import { settings } from '$lib/services/Settings';

	import { positionKey } from '$lib/types/MapObject';
	import { NoScanner } from '$lib/types/Tech';
	import type { Vector } from '$lib/types/Vector';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';

	const game = getContext<FullGame>('game');
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

			game.universe.planets.forEach((planet) =>
				scannersByPosition.set(positionKey(planet), {
					position: planet.position,
					scanRange: planet.spec?.scanRange ?? 0,
					scanRangePen: planet.spec?.scanRangePen ?? 0
				})
			);

			game.universe.fleets
				.filter((fleet) => fleet.spec?.scanner)
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

			game.universe.mineralPackets
				.filter((packet) => packet.scanRange != NoScanner || packet.scanRangePen != NoScanner)
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

			scanners = [];
			scanners = Array.from(scannersByPosition.values());
		}
	}
</script>

{#if $settings.showScanners}
	{#each scanners as scanner}
		<circle
			cx={$xGet({ position: scanner.position })}
			cy={$yGet({ position: scanner.position })}
			r={$xScale(scanner.scanRange)}
			class="scanner"
		/>
	{/each}
	{#each scanners as scanner}
		{#if scanner.scanRangePen > 0}
			<circle
				cx={$xGet({ position: scanner.position })}
				cy={$yGet({ position: scanner.position })}
				r={$xScale(scanner.scanRangePen)}
				class="scanner-pen"
			/>
		{/if}
	{/each}
{/if}
