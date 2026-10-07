import type { Universe } from '$lib/services/Universe';
import { FleetSchema, MapObjectType, PlanetSchema, Prt, WormholeSchema } from '$lib/types/cs-proto';
import type { CommandedPlayer } from '$lib/types/Player';
import { PlanetViewState, PlayerSettings } from '$lib/types/PlayerSettings';
import { create } from '@bufbuild/protobuf';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { readScannerColors } from './colors';
import { buildFleets, buildPlanets, buildWaypointPaths, packetColor } from './scene';

vi.mock('../Scanner', () => ({ getEnemiesAndFriends: () => ({ enemies: false, friends: false }) }));

afterEach(() => {
	vi.restoreAllMocks();
	vi.unstubAllGlobals();
});

function palette(overrides: Record<string, string>) {
	vi.stubGlobal('getComputedStyle', () => ({
		getPropertyValue: (name: string) => overrides[name] ?? ''
	}));
	vi.spyOn(console, 'warn').mockImplementation(() => {});
	return readScannerColors({} as HTMLCanvasElement);
}

function buildOne(
	planetData: Parameters<typeof create<typeof PlanetSchema>>[1],
	viewState: PlanetViewState,
	playerOverrides: Partial<CommandedPlayer> = {},
	colorOverrides: Record<string, string> = {}
) {
	const colors = palette(colorOverrides);
	const planet = create(PlanetSchema, planetData);
	const universe = {
		planetIntels: [planet],
		getMapObjectsByPosition: () => [],
		getPlanet: () => planet
	} as unknown as Universe;
	const player = {
		num: 1,
		race: { prt: Prt.JOAT },
		isFriend: () => false,
		...playerOverrides
	} as unknown as CommandedPlayer;
	const settings = new PlayerSettings();
	settings.planetViewState = viewState;
	const [draw] = buildPlanets(universe, player, settings, undefined, undefined, colors);
	return { draw, colors };
}

describe('population view', () => {
	// expected radii come from vrgPopRad in Stars!, which counts hundreds of colonists
	it.each([
		{ colonists: 0, radius: 2 },
		{ colonists: 24, radius: 2 },
		{ colonists: 25, radius: 3 },
		{ colonists: 2500, radius: 11 },
		{ colonists: 25000, radius: 21 },
		{ colonists: 100000, radius: 21 }
	])('sizes a planet with $colonists kT of colonists like Stars!', ({ colonists, radius }) => {
		const { draw } = buildOne(
			{ mapObject: { num: 1, playerNum: 1 }, cargo: { colonists } },
			PlanetViewState.Population
		);
		expect(draw.pop?.radius).toBe(radius);
	});

	it('colors planets by owner using CSS colors', () => {
		const overrides = {
			'--scanner-planet-pop-owned': '#000001',
			'--scanner-planet-pop-owned-ring': '#000002',
			'--scanner-planet-pop-friend': '#000003',
			'--scanner-planet-pop-friend-ring': '#000004',
			'--scanner-planet-pop-enemy': '#000005',
			'--scanner-planet-pop-enemy-ring': '#000006'
		};
		const mine = buildOne(
			{ mapObject: { num: 1, playerNum: 1 }, cargo: { colonists: 100 } },
			PlanetViewState.Population,
			{},
			overrides
		).draw;
		const friend = buildOne(
			{ mapObject: { num: 1, playerNum: 2 }, cargo: { colonists: 100 } },
			PlanetViewState.Population,
			{ isFriend: () => true },
			overrides
		).draw;
		const enemy = buildOne(
			{ mapObject: { num: 1, playerNum: 3 }, cargo: { colonists: 100 } },
			PlanetViewState.Population,
			{},
			overrides
		).draw;

		expect(mine.pop).toMatchObject({ fill: '#000001', ring: '#000002' });
		expect(friend.pop).toMatchObject({ fill: '#000003', ring: '#000004' });
		expect(enemy.pop).toMatchObject({ fill: '#000005', ring: '#000006' });
	});

	it('draws unowned planets normally', () => {
		const { draw } = buildOne(
			{ mapObject: { num: 1 }, cargo: { colonists: 100 } },
			PlanetViewState.Population
		);
		expect(draw.pop).toBeUndefined();
		expect(draw.normal).toBeDefined();
	});

	it('uses the CSS name color', () => {
		const { draw } = buildOne(
			{ mapObject: { num: 1, playerNum: 1 } },
			PlanetViewState.Population,
			{},
			{ '--scanner-name': '#abcdef' }
		);
		expect(draw.nameColor).toBe('#abcdef');
	});
});

describe('mineral views', () => {
	it('graphs concentration at 1 pixel per 5%, up to 20', () => {
		const { draw } = buildOne(
			{ mapObject: { num: 1 }, mineralConcentration: { ironium: 100, boranium: 54, germanium: 4 } },
			PlanetViewState.MineralConcentration
		);
		expect(draw.minerals).toEqual([20, 10, 0]);
	});

	it('graphs surface minerals with the mineral scale as 20 pixels, rounded', () => {
		// with the default 5000kT scale each pixel is 250kT, rounded to the nearest pixel
		const { draw } = buildOne(
			{
				mapObject: { num: 1, playerNum: 1 },
				cargo: { ironium: 5000, boranium: 374, germanium: 375 }
			},
			PlanetViewState.SurfaceMinerals
		);
		expect(draw.minerals).toEqual([20, 1, 2]);
	});

	it('graphs our own planets even with no surface minerals', () => {
		const { draw } = buildOne(
			{ mapObject: { num: 1, playerNum: 1 } },
			PlanetViewState.SurfaceMinerals
		);
		expect(draw.minerals).toEqual([0, 0, 0]);
	});

	it("doesn't graph other planets with no known minerals", () => {
		const { draw } = buildOne(
			{ mapObject: { num: 1, playerNum: 2 } },
			PlanetViewState.SurfaceMinerals
		);
		expect(draw.minerals).toBeUndefined();
	});
});

it('uses the CSS color for owned packets while retaining foreign player colors', () => {
	const colors = palette({ '--scanner-owned-packet': '#123456' });
	const player = { num: 1 } as CommandedPlayer;
	const universe = { getPlayerColor: () => '#abcdef' } as unknown as Universe;
	const settings = new PlayerSettings();
	settings.showPlayerColors = true;

	expect(packetColor(1, player, universe, settings, colors)).toBe('#123456');
	expect(packetColor(2, player, universe, settings, colors)).toBe('#abcdef');
});

describe('planet value view', () => {
	function valueDraw(habitability: number, terraformedHabitability: number, prt = Prt.JOAT) {
		const colors = palette({});
		const planet = create(PlanetSchema, {
			mapObject: { num: 1 },
			spec: { habitability, terraformedHabitability }
		});
		const universe = {
			planetIntels: [planet],
			getMapObjectsByPosition: () => []
		} as unknown as Universe;
		const player = { num: 1, race: { prt } } as unknown as CommandedPlayer;
		const settings = new PlayerSettings();
		settings.planetViewState = PlanetViewState.Percent;

		const [draw] = buildPlanets(universe, player, settings, undefined, undefined, colors);
		return { value: draw.value, colors };
	}

	// expected radii come from the scanViewPlanetValue case in Stars! DrawScanner
	it.each([
		{ habitability: 100, outer: 10, inner: 8 },
		{ habitability: 50, outer: 6, inner: 4 },
		{ habitability: 10, outer: 2, inner: 1 },
		{ habitability: 0, outer: 2, inner: 1 },
		{ habitability: 33, outer: 5, inner: 3 }
	])('sizes a $habitability% planet like Stars!', ({ habitability, outer, inner }) => {
		const { value, colors } = valueDraw(habitability, habitability);
		expect(value?.outerRadius).toBe(outer);
		expect(value?.innerRadius).toBe(inner);
		expect(value?.inner).toBe(colors.planet.habitable);
		expect(value?.outer).toBe(colors.planet.habitableRing);
	});

	it('shows the terraformed value in yellow when terraforming makes a planet habitable', () => {
		const { value, colors } = valueDraw(-10, 30);
		expect(value?.outerRadius).toBe(4);
		expect(value?.innerRadius).toBe(3);
		expect(value?.inner).toBe(colors.planet.terraformable);
		expect(value?.outer).toBe(colors.planet.terraformableRing);
	});

	it('sizes uninhabitable planets by their terraformed value, in red', () => {
		const { value, colors } = valueDraw(-45, -20);
		expect(value?.outerRadius).toBe(6);
		expect(value?.innerRadius).toBe(4);
		expect(value?.inner).toBe(colors.planet.uninhabitable);
		expect(value?.outer).toBe(colors.planet.uninhabitableRing);
	});

	it('shows CA races the terraformed value in green', () => {
		const { value, colors } = valueDraw(20, 80, Prt.CA);
		expect(value?.outerRadius).toBe(9);
		expect(value?.inner).toBe(colors.planet.habitable);
	});
});

describe('fleets in deep space', () => {
	function fleetsAt(...fleets: { num: number; x: number; ships: number }[]) {
		return fleets.map((f) =>
			create(FleetSchema, {
				mapObject: { num: f.num, playerNum: 1, position: { x: f.x, y: 0 } },
				tokens: [{ quantity: f.ships }]
			})
		);
	}

	function build(fleets: ReturnType<typeof fleetsAt>, commandedNum?: number) {
		const colors = palette({});
		const universe = { getAllFleets: () => fleets } as unknown as Universe;
		const player = { num: 1, isFriend: () => false } as unknown as CommandedPlayer;
		const settings = new PlayerSettings();
		settings.showFleetTokenCounts = true;
		const commanded = fleets.find((f) => f.mapObject?.num === commandedNum);
		return buildFleets(
			universe,
			player,
			settings,
			commanded as Parameters<typeof buildFleets>[3],
			colors
		);
	}

	it('shows one count for all the fleets at a location, like Stars!', () => {
		const draws = build(
			fleetsAt(
				{ num: 1, x: 10, ships: 2 },
				{ num: 2, x: 10, ships: 3 },
				{ num: 3, x: 50, ships: 1 }
			)
		);
		const counts = draws
			.filter((d) => d.count)
			.map((d) => [d.fleet.mapObject?.num, d.count?.value]);
		expect(counts).toEqual([
			[2, 5],
			[3, 1]
		]);
	});

	it('draws the commanded fleet last and puts the count on it', () => {
		const draws = build(fleetsAt({ num: 1, x: 10, ships: 2 }, { num: 2, x: 10, ships: 3 }), 1);
		expect(draws.map((d) => d.fleet.mapObject?.num)).toEqual([2, 1]);
		expect(draws[1].commanded).toBe(true);
		expect(draws[1].count?.value).toBe(5);
		expect(draws[0].count).toBeUndefined();
	});
});

describe('waypoint paths through wormholes', () => {
	function pathFor(destinationNum: number) {
		const wormholes = [
			create(WormholeSchema, {
				mapObject: { num: 1, position: { x: 10, y: 10 } },
				destinationNum
			}),
			create(WormholeSchema, { mapObject: { num: 2, position: { x: 60, y: 60 } } })
		];
		const fleet = create(FleetSchema, {
			mapObject: { num: 1, playerNum: 1 },
			fleetOrders: {
				waypoints: [
					{ position: { x: 0, y: 0 } },
					{
						position: { x: 10, y: 10 },
						mapObjectTarget: { targetType: MapObjectType.WORMHOLE, targetNum: 1 }
					},
					{ position: { x: 100, y: 60 } }
				]
			}
		});
		const universe = {
			fleets: [fleet],
			getWormhole: (num: number) => wormholes.find((w) => w.mapObject?.num === num)
		} as unknown as Universe;
		const player = { num: 1 } as CommandedPlayer;
		const [path] = buildWaypointPaths(universe, player, new PlayerSettings(), undefined);
		return path;
	}

	it('continues from the exit of a wormhole we know the destination of', () => {
		const { exits } = pathFor(2);
		expect(exits[0]).toBeUndefined();
		expect(exits[1]).toMatchObject({ x: 60, y: 60 });
		expect(exits[2]).toBeUndefined();
	});

	it("continues from the entrance when we don't know where a wormhole goes", () => {
		expect(pathFor(0).exits).toEqual([undefined, undefined, undefined]);
	});
});
