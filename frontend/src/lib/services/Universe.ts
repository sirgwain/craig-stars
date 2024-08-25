import {
	battlesSortBy,
	getBattleRecordDetails,
	type BattleRecord,
	type BattleRecordDetails
} from '$lib/types/Battle';
import type { Cost } from '$lib/types/Cost';
import { fleetsSortBy, type Fleet, type Target, type Waypoint } from '$lib/types/Fleet';
import { MapObjectType, type MapObject } from '$lib/types/MapObject';
import type { MineField } from '$lib/types/MineField';
import type { MineralPacket } from '$lib/types/MineralPacket';
import type { MysteryTrader } from '$lib/types/MysteryTrader';
import { planetsSortBy, type Planet } from '$lib/types/Planet';
import type {
	Player,
	PlayerIntel,
	PlayerIntels,
	PlayerScore,
	PlayerUniverse
} from '$lib/types/Player';
import type { ProductionQueueItem } from '$lib/types/Production';
import type { Salvage } from '$lib/types/Salvage';
import type { ShipDesign } from '$lib/types/ShipDesign';
import type { TechStore } from '$lib/types/Tech';
import type { Vector } from '$lib/types/Vector';
import type { Wormhole } from '$lib/types/Wormhole';
import { groupBy, startCase } from 'lodash-es';

export interface DesignFinder {
	getDesign(playerNum: number, num: number): ShipDesign | undefined;
	getMyDesign(num: number | undefined): ShipDesign | undefined;
}

export interface CostFinder {
	getItemCost(
		item: ProductionQueueItem | undefined,
		designFinder: DesignFinder,
		techStore: TechStore,
		planet?: Planet,
		quantity?: number
	): Cost;
}

export interface PlayerFinder {
	getPlayerIntel(num: number): PlayerIntel | undefined;
	getPlayerName(playerNum: number | undefined): string;
	getPlayerPluralName(playerNum: number | undefined): string;
	getPlayerColor(playerNum: number | undefined): string;
}
const sortByNum = (a: MapObject, b: MapObject) => a.num - b.num;

function addtoDict(mo: MapObject, dict: Record<string, MapObject[]>) {
	const key = positionKey(mo);
	if (!dict[key]) {
		dict[key] = [];
	}
	dict[key].push(mo);
}

function positionKey(pos: MapObject | Vector): string {
	const mo = 'position' in pos && (pos as MapObject);
	const v = 'x' in pos && (pos as Vector);
	if (mo) {
		return `${mo.position.x},${mo.position.y}`;
	} else if (v) {
		return `${v.x},${v.y}`;
	}
	return '';
}

export class Universe implements PlayerUniverse, PlayerIntels, DesignFinder {
	playerNum = 0;
	planets: Planet[] = [];
	fleets: Fleet[] = [];
	salvages: Salvage[] = [];
	mineFields: MineField[] = [];
	mineralPackets: MineralPacket[] = [];
	starbases: Fleet[] = [];
	wormholes: Wormhole[] = [];
	mysteryTraders: MysteryTrader[] = [];
	designs: ShipDesign[] = [];
	players: PlayerIntel[] = [];
	scores: PlayerScore[][] = [];
	battles: BattleRecord[] = [];

	mapObjectsByPosition: Record<string, MapObject[]> = {};
	myMapObjectsByPosition: Record<string, MapObject[]> = {};

	public setData(data: PlayerUniverse & PlayerIntels): Universe {
		this.designs = data.designs ?? [];
		this.battles = data.battles ?? [];
		this.players = data.players ?? [];
		this.scores = data.scores ?? [];

		this.planets = data.planets ?? [];
		this.fleets = data.fleets ?? [];
		this.starbases = data.starbases ?? [];
		this.mineFields = data.mineFields ?? [];
		this.mineralPackets = data.mineralPackets ?? [];
		this.salvages = data.salvages ?? [];
		this.wormholes = data.wormholes ?? [];
		this.mysteryTraders = data.mysteryTraders ?? [];

		this.starbases.forEach((sb) => (this.planets[(sb.planetNum ?? 1) - 1].starbase = sb));
		this.resetMapObjectsByPosition();
		return this;
	}

	// set the player number and update player specific internal data
	public setPlayer(playerNum: number) {
		this.playerNum = playerNum;
		this.resetMyMapObjectsByPosition();
	}

	// reset all data in the universe
	public resetData(playerNum: number, data: PlayerUniverse & PlayerIntels): Universe {
		this.setData(data);
		this.setPlayer(playerNum);
		return this;
	}

	resetMapObjectsByPosition() {
		this.mapObjectsByPosition = {};
		this.planets.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.fleets.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.mineFields.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.mineralPackets.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.salvages.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.wormholes.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.mysteryTraders.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
	}

	resetMyMapObjectsByPosition() {
		// build a map of objects owned by me
		this.myMapObjectsByPosition = {};
		const ownedByMe = (mo: MapObject) => mo.playerNum === this.playerNum;
		this.planets
			.filter(ownedByMe)
			.sort(sortByNum)
			.forEach((mo) => addtoDict(mo, this.myMapObjectsByPosition));
		this.fleets
			.filter(ownedByMe)
			.sort(sortByNum)
			.forEach((mo) => addtoDict(mo, this.myMapObjectsByPosition));
		this.mineFields
			.filter(ownedByMe)
			.sort(sortByNum)
			.forEach((mo) => addtoDict(mo, this.myMapObjectsByPosition));
		this.mineralPackets
			.filter(ownedByMe)
			.sort(sortByNum)
			.forEach((mo) => addtoDict(mo, this.myMapObjectsByPosition));
	}

	getPlayerIntel(num: number): PlayerIntel | undefined {
		if (num >= 1 && num <= this.players.length) {
			return this.players[num - 1];
		}
	}

	getPlayerScoreHistory(num: number): PlayerScore[] | undefined {
		if (num >= 1 && num <= this.scores.length && this.scores[num - 1]) {
			return this.scores[num - 1];
		}
	}

	getPlayerScore(num: number): PlayerScore | undefined {
		if (num >= 1 && num <= this.scores.length) {
			const history = this.scores[num - 1];
			if (history && history.length > 0) {
				return history[history.length - 1];
			}
		}
	}

	getPlayerName(playerNum: number | undefined): string {
		if (playerNum && playerNum > 0 && playerNum <= this.players.length) {
			const intel = this.players[playerNum - 1];
			return intel.raceName ?? intel.name;
		}
		return 'unknown';
	}

	getPlayerPluralName(playerNum: number | undefined): string {
		if (playerNum && playerNum > 0 && playerNum <= this.players.length) {
			const intel = this.players[playerNum - 1];
			return intel.racePluralName ?? intel.name;
		}
		return 'unknown';
	}

	getPlayerColor(playerNum: number | undefined): string {
		if (playerNum && playerNum > 0 && playerNum <= this.players.length) {
			const intel = this.players[playerNum - 1];
			return intel.color ?? '#FF0000';
		}
		return '#FF0000';
	}

	getMyDesigns(): ShipDesign[] {
		return this.designs.filter((d) => d.playerNum === this.playerNum);
	}

	getMyPlanets(sortKey: string, descending: boolean): Planet[] {
		const planets = this.planets.filter((d) => d.playerNum === this.playerNum);
		planets.sort(planetsSortBy(sortKey));
		if (descending) {
			planets.reverse();
		}
		return planets;
	}

	getPlanets(sortKey: string, descending: boolean): Planet[] {
		const planets = [...this.planets];
		planets.sort(planetsSortBy(sortKey));
		if (descending) {
			planets.reverse();
		}
		return planets;
	}

	getMyFleets(sortKey: string, descending: boolean): Fleet[] {
		const fleets = this.fleets.filter((d) => d.playerNum === this.playerNum);
		fleets.sort(fleetsSortBy(sortKey, this));
		if (descending) {
			fleets.reverse();
		}
		return fleets;
	}

	getFleets(sortKey: string, descending: boolean): Fleet[] {
		const fleets = [...this.fleets];
		fleets.sort(fleetsSortBy(sortKey, this));
		if (descending) {
			fleets.reverse();
		}
		return fleets;
	}

	getBattles(sortKey: string, descending: boolean, player: Player): BattleRecordDetails[] {
		const battles = this.battles.map((b) => getBattleRecordDetails(b, player, this));
		battles.sort(battlesSortBy(sortKey));
		if (descending) {
			battles.reverse();
		}
		return battles;
	}

	getDesign(playerNum: number, num: number): ShipDesign | undefined {
		return this.designs.find((d) => d.playerNum === playerNum && d.num === num);
	}

	getDesigns(playerNum: number): ShipDesign[] {
		return this.designs.filter((d) => d.playerNum === playerNum);
	}

	getMyDesign(num: number | undefined): ShipDesign | undefined {
		return this.designs.find((d) => d.playerNum === this.playerNum && d.num === num);
	}

	getBattle(num: number | undefined): BattleRecord | undefined {
		return this.battles.find((b) => b.num === num);
	}

	updateDesign(design: ShipDesign): Universe {
		const index = this.designs.findIndex((f) => f.num === design.num);
		if (index != -1) {
			this.designs = [...this.designs.slice(0, index), design, ...this.designs.slice(index + 1)];
		}

		return this;
	}

	addDesign(design: ShipDesign): Universe {
		this.designs = [...this.designs, design];
		return this;
	}

	getBattleLocation(battle: BattleRecord): string {
		if (battle.planetNum) {
			const planet = this.getPlanet(battle.planetNum);
			return planet?.name ?? 'Unknown';
		}
		return `Space (${battle.position.x}, ${battle.position.y})`;
	}

	getOtherMapObjectsHereByType(position: Vector) {
		return groupBy(this.mapObjectsByPosition[positionKey(position)], (mo) => mo.type);
	}

	getMapObjectsByPosition(position: MapObject | Vector) {
		return this.mapObjectsByPosition[positionKey(position)];
	}

	getSalvageAtPosition(position: MapObject | Vector): Salvage | undefined {
		const mo = this.getMapObjectsByPosition(position)?.find(
			(mo) => mo.type === MapObjectType.Salvage
		);
		if (mo) {
			return mo as Salvage;
		}
	}

	getMyMapObjectsByPosition(position: MapObject | Vector) {
		return this.myMapObjectsByPosition[positionKey(position)];
	}

	getMyPlanetsByPosition(position: MapObject | Vector): Planet[] {
		return (
			(this.getMyMapObjectsByPosition(position)?.filter(
				(mo) => mo.type === MapObjectType.Planet
			) as Planet[]) ?? []
		);
	}

	getMyFleetsByPosition(position: MapObject | Vector): Fleet[] {
		return (
			(this.getMyMapObjectsByPosition(position)?.filter(
				(mo) => mo.type === MapObjectType.Fleet
			) as Fleet[]) ?? []
		);
	}

	getPlanet(num: number) {
		return this.planets.find((p) => p.num === num);
	}

	getFleet(playerNum: number, num: number) {
		return this.fleets.find((f) => f.playerNum === playerNum && f.num === num);
	}

	getPlanetStarbase(planetNum: number) {
		return this.starbases.find((sb) => sb.planetNum === planetNum);
	}

	getWormhole(num: number) {
		return this.wormholes.find((w) => w.num === num);
	}

	getMysteryTrader(num: number) {
		return this.mysteryTraders.find((mt) => mt.num === num);
	}

	addFleets(fleets: Fleet[]) {
		this.fleets = [...fleets, ...this.fleets];
		this.resetMapObjectsByPosition();
		this.resetMyMapObjectsByPosition();
	}

	updateFleet(fleet: Fleet) {
		const index = this.fleets.findIndex((f) => f.num === fleet.num);
		if (index != -1) {
			this.fleets = [...this.fleets.slice(0, index), fleet, ...this.fleets.slice(index + 1)];
		}
		this.resetMapObjectsByPosition();
		this.resetMyMapObjectsByPosition();
	}

	updatePlanet(planet: Planet) {
		this.planets[planet.num - 1] = planet;
		this.resetMapObjectsByPosition();
		this.resetMyMapObjectsByPosition();
	}

	updateSalvages(salvages: Salvage[]) {
		this.salvages = salvages;
		this.resetMapObjectsByPosition();
		this.resetMyMapObjectsByPosition();
	}

	removeFleets(fleetNums: number[]) {
		this.fleets = this.fleets.filter((f) => fleetNums.indexOf(f.num) == -1);
		this.resetMapObjectsByPosition();
		this.resetMyMapObjectsByPosition();
	}

	getTargetName(wp: Waypoint): string {
		// first see if we can load this waypoint target as a mapobject
		// if so, use its name
		const mo = this.getMapObject(wp);
		if (mo) {
			if (mo.name && mo.name !== '') {
				return mo.name;
			}

			return `${startCase(mo.type)} #${mo.num}`;
		} else if (wp.targetName && wp.targetName !== '') {
			// we can't load it from the universe, see if the server gave us a target name
			return wp.targetName;
		}

		// we don't have a target name and we can't find the map object, just point it to the space location
		return `Space: (${wp.position.x.toFixed()}, ${wp.position.y.toFixed()})`;
	}

	// get a mapobject by type, number, and optionally player num
	getMapObject(target: Target): MapObject | undefined {
		switch (target.targetType) {
			case MapObjectType.Planet:
				return target.targetNum ? this.getPlanet(target.targetNum) : undefined;
			case MapObjectType.Fleet:
				return this.fleets.find(
					(f) => f.num === target.targetNum && f.playerNum === target.targetPlayerNum
				);
			case MapObjectType.MineField:
				return this.mineFields.find(
					(mf) => mf.num === target.targetNum && mf.playerNum === target.targetPlayerNum
				);
			case MapObjectType.MineralPacket:
				return this.mineralPackets.find(
					(p) => p.num === target.targetNum && p.playerNum === target.targetPlayerNum
				);
			case MapObjectType.Salvage:
				return this.salvages.find(
					(s) => s.num === target.targetNum && s.playerNum === target.targetPlayerNum
				);
			case MapObjectType.Wormhole:
				return target.targetNum ? this.getWormhole(target.targetNum) : undefined;
			case MapObjectType.MysteryTrader:
				return target.targetNum ? this.getMysteryTrader(target.targetNum) : undefined;
			case MapObjectType.PositionWaypoint:
				break;
		}
	}

	getHomeworld() {
		return this.planets.find((p) => p.playerNum === this.playerNum && p.homeworld);
	}
}
