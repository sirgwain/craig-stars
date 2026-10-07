import { dev } from '$app/env';

export const fixedScannerColors = {
	// These path and tick colors stay fixed because Stars! uses XOR to combine them.
	// Stars! draws all fleet paths with a pure blue pen, then XORs the selected fleet's path
	// with green, which makes it cyan
	waypointLine: '#0000FF',
	waypointLineCommanded: '#00FFFF',
	// year ticks are XORed onto the scanner with this color
	yearTick: '#00FF00'
};

export function readScannerColors(canvas: HTMLCanvasElement) {
	const style = getComputedStyle(canvas);
	// Preserve the classic palette if a CSS variable is missing or the stylesheet isn't loaded.
	const color = (name: string, fallback: string) => {
		const value = style.getPropertyValue(`--${name}`).trim();
		if (!value && dev) {
			console.warn(`Missing scanner color variable --${name}; using ${fallback}.`);
		}
		return value || fallback;
	};
	const marker = (name: string, fill: string, stroke: string) => ({
		fill: color(`${name}-fill`, fill),
		stroke: color(`${name}-stroke`, stroke)
	});
	const orbit = {
		none: color('orbit-neutral', '#ffffff'),
		friends: color('orbit-friends', '#ffff00'),
		enemies: color('orbit-enemies', '#ff0000')
	};

	return {
		...fixedScannerColors,
		background: color('scanner-background', '#000000'),
		salvage: color('scanner-salvage', '#ffff00'),
		myPacket: color('scanner-owned-packet', '#0900ff'),
		commandedFleet: color('scanner-commanded-fleet', '#ffff00'),
		name: color('scanner-name', '#ffffff'),
		selection: color('scanner-selection', '#ffffff'),
		mineralAxis: color('scanner-mineral-axis', '#c0c0c0'),
		warpLine: color('scanner-warp-line', '#ffffff'),
		planet: {
			owned: color('scanner-planet-owned', '#00ff00'),
			unexplored: color('scanner-planet-unexplored', '#999999'),
			explored: color('scanner-planet-explored', '#ffffff'),
			outline: color('scanner-planet-outline', '#999999'),
			habitable: color('scanner-planet-habitable', '#00ff00'),
			terraformable: color('scanner-planet-terraformable', '#ffff00'),
			uninhabitable: color('scanner-planet-uninhabitable', '#ff0000'),
			habitableRing: color('scanner-planet-habitable-ring', '#007f00'),
			terraformableRing: color('scanner-planet-terraformable-ring', '#7f7f00'),
			uninhabitableRing: color('scanner-planet-uninhabitable-ring', '#7f0000'),
			popOwned: color('scanner-planet-pop-owned', '#00ff00'),
			popOwnedRing: color('scanner-planet-pop-owned-ring', '#007f00'),
			popFriend: color('scanner-planet-pop-friend', '#ffff00'),
			popFriendRing: color('scanner-planet-pop-friend-ring', '#7f7f00'),
			popEnemy: color('scanner-planet-pop-enemy', '#ff0000'),
			popEnemyRing: color('scanner-planet-pop-enemy-ring', '#ff0000')
		},
		scanner: color('scanner-range', '#8b0000'),
		scannerPen: color('scanner-range-pen', '#7f7f00'),
		packetDestLine: color('scanner-packet-dest', '#6b0065'),
		routeLine: color('scanner-route', '#fdfd00'),
		wormhole: color('scanner-wormhole', '#8a2be2'),
		wormholeLink: color('scanner-wormhole', '#8a2be2'),
		mysteryTrader: color('mystery-trader', '#00ffff'),
		ironium: color('scanner-ironium', '#0000ff'),
		boranium: color('scanner-boranium', '#007f00'),
		germanium: color('scanner-germanium', '#ffff00'),
		starbase: marker('starbase', '#fdfd00', '#7f7f00'),
		starbaseFort: marker('starbase-fort', '#00007f', '#00002f'),
		stargate: marker('stargate', '#008100', '#005c00'),
		massDriver: marker('massdriver', '#8d0085', '#6b0065'),
		locator: color('scanner-locator', '#ffffff'),
		orbit: {
			ring: { ...orbit, both: color('orbit-mixed-ring', '#ff00ff') },
			text: { ...orbit, both: color('orbit-mixed-count', '#eb2eeb') }
		}
	};
}

export type ScannerColors = ReturnType<typeof readScannerColors>;
