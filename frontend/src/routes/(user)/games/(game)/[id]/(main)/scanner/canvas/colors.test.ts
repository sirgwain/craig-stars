import { afterEach, describe, expect, it, vi } from 'vitest';
import { readScannerColors } from './colors';

vi.mock('$app/env', () => ({ dev: true }));

afterEach(() => {
	vi.restoreAllMocks();
	vi.unstubAllGlobals();
});

describe('readScannerColors', () => {
	it('uses CSS overrides without warnings', () => {
		vi.stubGlobal('getComputedStyle', () => ({
			getPropertyValue: () => '  #123456  '
		}));
		const warn = vi.spyOn(console, 'warn').mockImplementation(() => {});

		const colors = readScannerColors({} as HTMLCanvasElement);

		expect(colors.scanner).toBe('#123456');
		expect(colors.starbase.stroke).toBe('#123456');
		expect(colors.orbit.ring.both).toBe('#123456');
		expect(colors.background).toBe('#123456');
		expect(colors.salvage).toBe('#123456');
		expect(colors.myPacket).toBe('#123456');
		expect(colors.commandedFleet).toBe('#123456');
		expect(colors.name).toBe('#123456');
		expect(colors.selection).toBe('#123456');
		expect(colors.mineralAxis).toBe('#123456');
		expect(colors.warpLine).toBe('#123456');
		expect(Object.values(colors.planet).every((value) => value === '#123456')).toBe(true);
		expect(colors.waypointLine).toBe('#0000FF');
		expect(colors.waypointLineCommanded).toBe('#00FFFF');
		expect(colors.yearTick).toBe('#00FF00');
		expect(warn).not.toHaveBeenCalled();
	});

	it('preserves the remaining classic colors when the stylesheet is missing', () => {
		vi.stubGlobal('getComputedStyle', () => ({ getPropertyValue: () => '' }));
		vi.spyOn(console, 'warn').mockImplementation(() => {});

		const colors = readScannerColors({} as HTMLCanvasElement);

		expect(colors).toMatchObject({
			background: '#000000',
			salvage: '#ffff00',
			myPacket: '#0900ff',
			commandedFleet: '#ffff00',
			name: '#ffffff',
			selection: '#ffffff',
			mineralAxis: '#c0c0c0',
			warpLine: '#ffffff',
			planet: {
				owned: '#00ff00',
				unexplored: '#999999',
				explored: '#ffffff',
				outline: '#999999',
				habitable: '#00ff00',
				terraformable: '#ffff00',
				uninhabitable: '#ff0000',
				habitableRing: '#007f00',
				terraformableRing: '#7f7f00',
				uninhabitableRing: '#7f0000',
				popOwned: '#00ff00',
				popOwnedRing: '#007f00',
				popFriend: '#ffff00',
				popFriendRing: '#7f7f00',
				popEnemy: '#ff0000',
				popEnemyRing: '#ff0000'
			}
		});
	});

	it.each(['', ' \t '])('falls back and warns for missing or blank variables (%j)', (value) => {
		vi.stubGlobal('getComputedStyle', () => ({
			getPropertyValue: (name: string) =>
				['--scanner-range', '--starbase-stroke', '--orbit-mixed-count'].includes(name)
					? value
					: '#123456'
		}));
		const warn = vi.spyOn(console, 'warn').mockImplementation(() => {});

		const colors = readScannerColors({} as HTMLCanvasElement);

		expect(colors.scanner).toBe('#8b0000');
		expect(colors.starbase.stroke).toBe('#7f7f00');
		expect(colors.orbit.text.both).toBe('#eb2eeb');
		expect(colors.scannerPen).toBe('#123456');
		expect(warn).toHaveBeenCalledTimes(3);
		expect(warn).toHaveBeenCalledWith(
			'Missing scanner color variable --scanner-range; using #8b0000.'
		);
	});
});
