/**
 * Scene builders turn the universe into flat lists of things to draw. These only rebuild when
 * the universe, settings or commanded objects change, not on every pan/zoom frame.
 */
import type { Universe } from '#lib/services/Universe.js';
import { population } from '#lib/types/Cargo.js';
import { NoScanner, None, ReportAgeUnexplored } from '#lib/types/Consts.js';
import {
	MapObjectTargetSchema,
	MapObjectType,
	Prt,
	type Fleet,
	type Minefield,
	type Planet,
	type Waypoint
} from '#lib/types/cs-proto.js';
import type { CommandedFleet } from '#lib/types/Fleet.js';
import { filterFleet } from '#lib/types/Filter.js';
import {
	equal,
	owned,
	positionKey,
	type MapObjectLike,
	type Position
} from '#lib/types/MapObject.js';
import type { CommandedPlanet } from '#lib/types/Planet.js';
import type { CommandedPlayer } from '#lib/types/Player.js';
import { PlanetViewState, type PlayerSettings } from '#lib/types/PlayerSettings.js';
import { emptyVector } from '#lib/types/Vector.js';
import { getDisplayColor } from '#lib/utils/colorUtils.js';
import { create } from '@bufbuild/protobuf';
import { getEnemiesAndFriends } from '../Scanner';
import type { ScannerColors } from './colors';

type OrbitKind = keyof ScannerColors['orbit']['ring'];

function orbitKind(enemies: boolean, friends: boolean): OrbitKind {
	if (friends && enemies) return 'both';
	if (friends) return 'friends';
	if (enemies) return 'enemies';
	return 'none';
}

export type ScannerCircle = { position: Position; scanRange: number; scanRangePen: number };

export function buildScanners(
	universe: Universe,
	player: CommandedPlayer,
	settings: PlayerSettings
): ScannerCircle[] {
	const scannersByPosition = new Map<string, ScannerCircle>();

	const add = (position: Position | undefined, scanRange: number, scanRangePen: number) => {
		const key = positionKey(position);
		const existing = scannersByPosition.get(key);
		if (existing) {
			existing.scanRange = Math.max(existing.scanRange, scanRange);
			existing.scanRangePen = Math.max(existing.scanRangePen, scanRangePen);
		} else {
			scannersByPosition.set(key, { position: position ?? emptyVector(), scanRange, scanRangePen });
		}
	};

	const addFleet = (fleet: Fleet) => {
		const spec = fleet.spec?.shipDesignSpec;
		if ((spec?.scanRange ?? 0) > 0 || (spec?.scanRangePen ?? 0) > 0) {
			add(fleet.mapObject?.position, spec?.scanRange ?? 0, spec?.scanRangePen ?? 0);
		}
	};

	if (settings.showScanners) {
		universe.planets
			.filter((p) => p.spec?.scanner)
			.forEach((p) =>
				add(p.mapObject?.position, p.spec?.scanRange ?? 0, p.spec?.scanRangePen ?? 0)
			);
		universe.fleets.forEach(addFleet);
		universe.mineralPackets
			.filter(
				(p) =>
					p.mapObject?.playerNum === player.num &&
					(p.scanRange != NoScanner || p.scanRangePen != NoScanner)
			)
			.forEach((p) => add(p.mapObject?.position, p.scanRange, p.scanRangePen));
	}

	if (settings.showAllyScanners) {
		universe.planetIntels
			.filter((p) => player.isSharingMap(p.mapObject?.playerNum) && p.spec?.scanner)
			.forEach((p) =>
				add(p.mapObject?.position, p.spec?.scanRange ?? 0, p.spec?.scanRangePen ?? 0)
			);
		universe.fleetIntels
			.filter((f) => player.isSharingMap(f.mapObject?.playerNum))
			.forEach(addFleet);
		universe.mineralPackets
			.filter(
				(p) =>
					p.mapObject?.playerNum !== player.num &&
					player.isSharingMap(p.mapObject?.playerNum) &&
					(p.scanRange != NoScanner || p.scanRangePen != NoScanner)
			)
			.forEach((p) => add(p.mapObject?.position, p.scanRange, p.scanRangePen));
	}

	return Array.from(scannersByPosition.values());
}

export type MinefieldDraw = { minefield: Minefield; radius: number; color: string };

export function buildMinefields(
	universe: Universe,
	player: CommandedPlayer,
	settings: PlayerSettings
): MinefieldDraw[] {
	return universe.allMinefields.map((minefield) => ({
		minefield,
		radius: Math.sqrt(minefield.numMines),
		color: getDisplayColor(minefield.mapObject?.playerNum, player, universe, settings)
	}));
}

export type DestLine = {
	from: Position;
	to: Position;
	width: number;
	// light years travelled per year, used to dash the line
	yearDist: number;
};

export function buildPacketDests(
	universe: Universe,
	commandedPlanet?: CommandedPlanet
): DestLine[] {
	return universe.planets
		.filter((p) => p.planetOrders?.packetTargetNum && p.planetOrders.packetTargetNum != None)
		.map((planet) => {
			const target = universe.getPlanet(planet.planetOrders?.packetTargetNum ?? None);
			const speed = planet.planetOrders?.packetSpeed ?? 0;
			return {
				from: planet.mapObject?.position ?? emptyVector(),
				to: target?.mapObject?.position ?? planet.mapObject?.position ?? emptyVector(),
				width: planet.mapObject?.num === commandedPlanet?.mapObject.num ? 1.5 : 1,
				yearDist: speed * speed
			};
		});
}

export function buildRouteDests(universe: Universe, commandedPlanet?: CommandedPlanet): DestLine[] {
	return universe.planets
		.filter((p) => (p.planetOrders?.routeTargetNum ?? None) != None)
		.map((planet) => {
			const target = universe.getMapObject(
				create(MapObjectTargetSchema, {
					targetPosition: emptyVector(),
					targetType: planet.planetOrders?.routeTargetType ?? MapObjectType.UNSPECIFIED,
					targetNum: planet.planetOrders?.routeTargetNum ?? 0,
					targetPlayerNum: planet.planetOrders?.routeTargetPlayerNum ?? 0
				})
			);
			const speed = planet.planetOrders?.packetSpeed ?? 0;
			return {
				from: planet.mapObject?.position ?? emptyVector(),
				to: target?.mapObject?.position ?? planet.mapObject?.position ?? emptyVector(),
				width: planet.mapObject?.num === commandedPlanet?.mapObject.num ? 1.5 : 1,
				yearDist: speed * speed
			};
		});
}

export type WaypointPath = {
	waypoints: Waypoint[];
	// for waypoints at a wormhole we know the destination of, where the fleet comes out
	exits: (Position | undefined)[];
	commanded: boolean;
};

// Fleets that reach a wormhole jump to its destination and continue from there, so a path
// through a wormhole we know the destination of continues from the exit.
function wormholeExits(universe: Universe, waypoints: Waypoint[]): (Position | undefined)[] {
	return waypoints.map((wp) => {
		const target = wp.mapObjectTarget;
		if (target?.targetType !== MapObjectType.WORMHOLE || !target.targetNum) {
			return undefined;
		}
		const destinationNum = universe.getWormhole(target.targetNum)?.destinationNum;
		return destinationNum ? universe.getWormhole(destinationNum)?.mapObject?.position : undefined;
	});
}

export function buildWaypointPaths(
	universe: Universe,
	player: CommandedPlayer,
	settings: PlayerSettings,
	commandedFleet: CommandedFleet | undefined
): WaypointPath[] {
	const paths: WaypointPath[] = universe.fleets
		.filter(
			(f) =>
				(equal(commandedFleet, f) || filterFleet(player, f, settings)) &&
				f.mapObject?.num !== commandedFleet?.mapObject.num &&
				(f.fleetOrders?.waypoints?.length ?? 0) > 1
		)
		.map((f) => {
			const waypoints = f.fleetOrders?.waypoints ?? [];
			return { waypoints, exits: wormholeExits(universe, waypoints), commanded: false };
		});

	// draw the commanded fleet's path last so it's on top
	if (commandedFleet && commandedFleet.fleetOrders.waypoints.length > 1) {
		const waypoints = commandedFleet.fleetOrders.waypoints;
		paths.push({ waypoints, exits: wormholeExits(universe, waypoints), commanded: true });
	}
	return paths;
}

/**
 * Everything needed to draw a planet. All offsets are in icon units and are multiplied by the
 * view's iconScale when drawn. top/bottom are the vertical extents of the icon so names, counts
 * and the selection arrow can be placed just outside of it.
 */
export type PlanetDraw = {
	planet: Planet;
	top: number;
	bottom: number;
	nameColor: string;
	// the classic dot with orbit ring and starbase markers
	normal?: {
		radius: number;
		fill: string;
		strokeWidth: number;
		ring?: { radius: number; width: number; color: string; dashed: boolean };
		starbaseWidth: number;
		starbaseXOffset: number;
		starbaseYOffset: number;
		starbase?: { fill: string; stroke: string };
		stargate: boolean;
		massDriver: boolean;
	};
	// the planet value view: a bright disc inside a darker ring, sized by the planet's value
	value?: {
		outerRadius: number;
		innerRadius: number;
		outer: string;
		inner: string;
		flag?: string;
	};
	// the population view: a circle sized by the population (radius before zooming out, in pixels)
	pop?: { radius: number; fill: string; ring: string };
	// the mineral views: ironium, boranium and germanium bar heights, 0-20 pixels
	minerals?: [number, number, number];
	count?: { value: number; color: string };
};

// the planet value view's flag pole height and maximum disc radius, in pixels
const valueFlagHeight = 20;
const valueMaxRadius = 10;
// Population view circles grow a pixel each time the population passes one of these, in
// hundreds of colonists, from vrgPopRad in Stars!
const popRadiusThresholds = [
	25, 50, 100, 200, 400, 800, 1000, 1500, 2250, 3000, 4000, 5000, 6000, 7500, 9000, 11000, 14000,
	18000, 25000
];
// mineral bars are at most this many pixels tall
export const mineralBarMax = 20;

export function buildPlanets(
	universe: Universe,
	player: CommandedPlayer,
	settings: PlayerSettings,
	commandedMapObject: MapObjectLike | undefined,
	commandedPlanet: CommandedPlanet | undefined,
	colors: ScannerColors
): PlanetDraw[] {
	return universe.planetIntels.map((planet) => {
		const position = planet.mapObject?.position ?? emptyVector();
		const orbiting = universe
			.getMapObjectsByPosition(position)
			.filter((mo) => mo.mapObject?.type === MapObjectType.FLEET);
		const orbitingVisible = orbiting.filter((f) => filterFleet(player, f as Fleet, settings));

		const commanded =
			(commandedMapObject?.mapObject?.type === MapObjectType.FLEET &&
				(commandedMapObject as Fleet).orbitingPlanetNum === planet.mapObject?.num) ||
			commandedPlanet?.mapObject.num === planet.mapObject?.num;

		const draw: PlanetDraw = {
			planet,
			top: 0,
			bottom: 0,
			nameColor:
				settings.showPlayerColors &&
				planet.mapObject?.playerNum &&
				planet.mapObject.playerNum !== player.num
					? getDisplayColor(planet.mapObject.playerNum, player, universe, settings)
					: colors.name
		};

		const explored = planet.mapObject?.reportAge !== ReportAgeUnexplored;
		switch (settings.planetViewState) {
			case PlanetViewState.Percent:
				if (explored) {
					setPlanetValue(draw, planet, player, universe, settings, colors);
				} else {
					setNormal(draw, planet, player, universe, settings, orbitingVisible, false, colors);
				}
				break;
			case PlanetViewState.Population:
				// like Stars!, only inhabited planets get a population circle
				if (explored && planet.mapObject?.playerNum) {
					setPopulation(draw, planet, player, universe, settings, colors);
				} else {
					setNormal(draw, planet, player, universe, settings, orbitingVisible, false, colors);
				}
				break;
			case PlanetViewState.MineralConcentration: {
				setNormal(draw, planet, player, universe, settings, orbitingVisible, false, colors);
				if (explored) {
					// 1 pixel per 5% concentration
					const conc = planet.mineralConcentration;
					draw.minerals = [conc?.ironium ?? 0, conc?.boranium ?? 0, conc?.germanium ?? 0].map((c) =>
						Math.min(Math.trunc(c / 5), mineralBarMax)
					) as [number, number, number];
				}
				break;
			}
			case PlanetViewState.SurfaceMinerals: {
				setNormal(draw, planet, player, universe, settings, orbitingVisible, false, colors);
				// we always know our own surface minerals. Planet intel doesn't say whether minerals
				// were discovered, so only graph other planets with minerals we know about
				const cargo = universe.getPlanet(planet.mapObject?.num)?.cargo ?? planet.cargo;
				const owned = planet.mapObject?.playerNum === player.num;
				if (owned || cargo?.ironium || cargo?.boranium || cargo?.germanium) {
					// the mineral scale is 20 pixels, rounded like Stars!
					const max = settings.mineralScale;
					const unit = Math.max(Math.trunc(max / 20), 1);
					draw.minerals = [cargo?.ironium ?? 0, cargo?.boranium ?? 0, cargo?.germanium ?? 0].map(
						(amount) => Math.min(Math.trunc((Math.trunc(max / 40) + amount) / unit), mineralBarMax)
					) as [number, number, number];
				}
				break;
			}
			default:
				setNormal(draw, planet, player, universe, settings, orbitingVisible, commanded, colors);
		}

		if (settings.showFleetTokenCounts) {
			const value = orbitingVisible.reduce(
				(count, f) => count + (f as Fleet).tokens.reduce((c, t) => c + t.quantity, 0),
				0
			);
			if (value) {
				const { enemies, friends } = getEnemiesAndFriends(orbiting, player);
				draw.count = { value, color: colors.orbit.text[orbitKind(enemies, friends)] };
			}
		}

		return draw;
	});
}

function setNormal(
	draw: PlanetDraw,
	planet: Planet,
	player: CommandedPlayer,
	universe: Universe,
	settings: PlayerSettings,
	orbitingFleets: MapObjectLike[],
	commanded: boolean,
	colors: ScannerColors
) {
	// green for us, gray for unexplored, white for explored
	let fill = colors.planet.unexplored;
	if (planet.mapObject?.playerNum === player.num) {
		fill = colors.planet.owned;
	} else if (planet.mapObject?.playerNum) {
		fill = getDisplayColor(planet.mapObject.playerNum, player, universe, settings);
	} else if (planet.mapObject?.reportAge !== ReportAgeUnexplored) {
		fill = colors.planet.explored;
	}

	const radius = owned(planet) ? (commanded ? 6 : 3) : commanded ? 4 : 2;
	const ringRadius = radius * 2.5;
	const starbaseWidth = commanded ? 6 : 4;
	const starbaseSpec = planet.spec?.planetStarbaseSpec;

	let ring: NonNullable<PlanetDraw['normal']>['ring'];
	if (orbitingFleets.length > 0) {
		const { enemies, friends } = getEnemiesAndFriends(orbitingFleets, player);
		const kind = orbitKind(enemies, friends);
		ring = {
			radius: ringRadius,
			width: commanded ? 2 : 1.5,
			color: colors.orbit.ring[kind],
			dashed: kind === 'both'
		};
	}

	draw.normal = {
		radius,
		fill,
		strokeWidth: commanded ? 1 : 0.5,
		ring,
		starbaseWidth,
		starbaseXOffset: ringRadius * 0.75,
		starbaseYOffset: ringRadius + starbaseWidth,
		starbase: starbaseSpec?.hasStarbase
			? starbaseSpec.dockCapacity
				? colors.starbase
				: colors.starbaseFort
			: undefined,
		stargate: !!starbaseSpec?.hasStargate,
		massDriver: !!starbaseSpec?.hasMassDriver
	};

	const outer = ring ? ring.radius + ring.width / 2 : radius;
	draw.bottom = outer;
	draw.top = -outer;
	if (draw.normal.starbase || draw.normal.stargate) {
		draw.top = Math.min(draw.top, -draw.normal.starbaseYOffset);
	}
	if (draw.normal.massDriver) {
		draw.top = Math.min(draw.top, -draw.normal.starbaseYOffset - starbaseWidth / 2);
	}
}

/**
 * The planet value view, ported from the scanViewPlanetValue case of DrawScanner in Stars!
 * A planet's value is its habitability for us. If it's uninhabitable, show what it could be
 * terraformed to instead: yellow if that makes it habitable, red if not. CA races terraform
 * instantly, so their planets always show the terraformed value in green.
 */
function setPlanetValue(
	draw: PlanetDraw,
	planet: Planet,
	player: CommandedPlayer,
	universe: Universe,
	settings: PlayerSettings,
	colors: ScannerColors
) {
	const instantTerraform = player.race?.prt === Prt.CA;
	let value = planet.spec?.habitability ?? 0;
	let terraformable = false;
	if (value < 0 || instantTerraform) {
		value = planet.spec?.terraformedHabitability ?? value;
		terraformable = value >= 0 && !instantTerraform;
	}

	// the radius grows 1px per 11% of value, or per 5% of negative value, and the inner disc
	// is 2px smaller, but never smaller than 1px
	const outerRadius = Math.min(
		value >= 0 ? Math.trunc(value / 11) + 2 : Math.trunc(-value / 5) + 2,
		valueMaxRadius
	);
	let innerRadius = outerRadius - 2;
	if (innerRadius < 3) {
		innerRadius++;
	}
	innerRadius = Math.max(innerRadius, 1);

	draw.value = {
		outerRadius,
		innerRadius,
		outer:
			value < 0
				? colors.planet.uninhabitableRing
				: terraformable
					? colors.planet.terraformableRing
					: colors.planet.habitableRing,
		inner:
			value < 0
				? colors.planet.uninhabitable
				: terraformable
					? colors.planet.terraformable
					: colors.planet.habitable,
		flag: planet.mapObject?.playerNum
			? getDisplayColor(planet.mapObject.playerNum, player, universe, settings)
			: undefined
	};
	// Stars! ellipses cover 2r+1 pixels, so they extend half a pixel past the radius
	draw.bottom = outerRadius + 0.5;
	draw.top = draw.value.flag ? -valueFlagHeight : -(outerRadius + 0.5);
}

/**
 * The population view, ported from the scanViewPopulation case of DrawScanner in Stars!
 * The radius steps up through popRadiusThresholds, colored by who owns the planet.
 */
function setPopulation(
	draw: PlanetDraw,
	planet: Planet,
	player: CommandedPlayer,
	universe: Universe,
	settings: PlayerSettings,
	colors: ScannerColors
) {
	const pop = population(planet.cargo) / 100;
	let steps = 0;
	while (steps < popRadiusThresholds.length && popRadiusThresholds[steps] <= pop) {
		steps++;
	}

	const playerNum = planet.mapObject?.playerNum;
	let fill = colors.planet.popEnemy;
	let ring = colors.planet.popEnemyRing;
	if (playerNum === player.num) {
		fill = colors.planet.popOwned;
		ring = colors.planet.popOwnedRing;
	} else if (settings.showPlayerColors) {
		fill = ring = getDisplayColor(playerNum, player, universe, settings);
	} else if (player.isFriend(playerNum)) {
		fill = colors.planet.popFriend;
		ring = colors.planet.popFriendRing;
	}

	draw.pop = { radius: steps + 2, fill, ring };
}

export type FleetDraw = {
	fleet: Fleet;
	color: string;
	commanded: boolean;
	angle: number;
	count?: { value: number; color: string };
};

/**
 * Fleets in deep space. Like Stars!, every fleet at a location is drawn, but the location gets
 * one token count for all the fleets there instead of overlapping counts. The commanded fleet
 * is drawn last so it's never hidden under another fleet.
 */
export function buildFleets(
	universe: Universe,
	player: CommandedPlayer,
	settings: PlayerSettings,
	commandedFleet: CommandedFleet | undefined,
	colors: ScannerColors
): FleetDraw[] {
	const fleets = universe
		.getAllFleets()
		.filter((f) => !f.orbitingPlanetNum)
		.filter((f) => equal(commandedFleet, f) || filterFleet(player, f, settings));

	const draws: FleetDraw[] = fleets.map((fleet) => ({
		fleet,
		color: getDisplayColor(fleet.mapObject?.playerNum, player, universe, settings),
		commanded:
			commandedFleet?.mapObject.num === fleet.mapObject?.num &&
			commandedFleet?.mapObject.playerNum === fleet.mapObject?.playerNum,
		angle: headingAngle(fleet.heading)
	}));
	draws.sort((a, b) => Number(a.commanded) - Number(b.commanded));

	if (settings.showFleetTokenCounts) {
		// put one count for each location on the last fleet drawn there, so it's on top
		const byPosition = new Map<string, FleetDraw[]>();
		for (const draw of draws) {
			const key = positionKey(draw.fleet);
			byPosition.set(key, [...(byPosition.get(key) ?? []), draw]);
		}
		for (const here of byPosition.values()) {
			const { enemies, friends } = getEnemiesAndFriends(
				here.map((d) => d.fleet),
				player
			);
			here[here.length - 1].count = {
				value: here.reduce(
					(c, d) => c + d.fleet.tokens.reduce((t, token) => t + token.quantity, 0),
					0
				),
				color: colors.orbit.text[orbitKind(enemies, friends)]
			};
		}
	}
	return draws;
}

// the fleet triangle points up and to the right by default, so rotate it 225º more
export function headingAngle(heading: { x: number; y: number } | undefined): number {
	return Math.atan2(heading?.y ?? 0, heading?.x ?? 0) + (225 * Math.PI) / 180;
}

export function packetColor(
	playerNum: number | undefined,
	player: CommandedPlayer,
	universe: Universe,
	settings: PlayerSettings,
	colors: ScannerColors
): string {
	return playerNum === player.num
		? colors.myPacket
		: getDisplayColor(playerNum, player, universe, settings);
}

export type WormholeLink = { from: Position; to: Position };

export function buildWormholeLinks(universe: Universe): WormholeLink[] {
	const numsUsed = new Set<number>();
	const links: WormholeLink[] = [];
	for (const wormhole of universe.wormholes) {
		if (!wormhole.destinationNum) {
			continue;
		}
		const num = wormhole.mapObject?.num ?? 0;
		const used = numsUsed.has(num) || numsUsed.has(wormhole.destinationNum);
		numsUsed.add(num);
		numsUsed.add(wormhole.destinationNum);
		if (used) {
			continue;
		}

		const from = wormhole.mapObject?.position ?? emptyVector();
		const to = universe.getWormhole(wormhole.destinationNum)?.mapObject?.position ?? from;
		const dx = to.x - from.x;
		const dy = to.y - from.y;
		const length = Math.hypot(dx, dy) || 1;
		// pull each end in 3ly so the line doesn't cover the wormholes
		links.push({
			from: { x: from.x + (dx / length) * 3, y: from.y + (dy / length) * 3 },
			to: { x: to.x - (dx / length) * 3, y: to.y - (dy / length) * 3 }
		});
	}
	return links;
}
