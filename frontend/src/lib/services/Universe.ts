import type { BattleRecord } from '$lib/types/Battle';
import type { Fleet, Target, Waypoint } from '$lib/types/Fleet';
import { MapObjectType, type MapObject } from '$lib/types/MapObject';
import type { MineField } from '$lib/types/MineField';
import type { MineralPacket } from '$lib/types/MineralPacket';
import type { MysteryTrader } from '$lib/types/MysteryTrader';
import type { Planet } from '$lib/types/Planet';
import type { PlayerIntel, PlayerIntels, PlayerUniverse } from '$lib/types/Player';
import type { Salvage } from '$lib/types/Salvage';
import type { ShipDesign } from '$lib/types/ShipDesign';
import type { Vector } from '$lib/types/Vector';
import type { Wormhole } from '$lib/types/Wormhole';
import { groupBy, startCase } from 'lodash-es';

export interface DesignFinder {
	getDesign(playerNum: number, num: number): ShipDesign | undefined;
	getMyDesignByName(playerNum: number, name: string): ShipDesign | undefined;
}

export interface PlayerFinder {
	getPlayerIntel(num: number): PlayerIntel | undefined;
	getPlayerName(playerNum: number | undefined): string;
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
	battles: BattleRecord[] = [];

	mapObjectsByPosition: Record<string, MapObject[]> = {};
	myMapObjectsByPosition: Record<string, MapObject[]> = {};

	public setData(playerNum: number, data: PlayerUniverse & PlayerIntels) {
		this.mapObjectsByPosition = {};
		this.myMapObjectsByPosition = {};
		this.playerNum = playerNum;
		this.designs = data.designs ?? [];
		this.battles = data.battles ?? [];
		this.players = data.players ?? [];

		this.planets = data.planets ?? [];
		this.fleets = data.fleets ?? [];
		this.mineFields = data.mineFields ?? [];
		this.mineralPackets = data.mineralPackets ?? [];
		this.salvages = data.salvages ?? [];
		this.wormholes = data.wormholes ?? [];
		this.mysteryTraders = data.mysteryTraders ?? [];

		this.planets.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.fleets.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.mineFields.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.mineralPackets.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.salvages.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.wormholes.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.mysteryTraders.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));

		this.resetMyMapObjectsByPosition();
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
		this.salvages
			.filter(ownedByMe)
			.sort(sortByNum)
			.forEach((mo) => addtoDict(mo, this.myMapObjectsByPosition));
	}

	getPlayerIntel(num: number): PlayerIntel | undefined {
		if (num >= 1 && num <= this.players.length) {
			return this.players[num - 1];
		}
	}

	getPlayerName(playerNum: number | undefined): string {
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

	getMyDesigns(playerNum: number): ShipDesign[] {
		return this.designs.filter((d) => d.playerNum === playerNum);
	}

	getMyPlanets(playerNum: number): Planet[] {
		return this.planets.filter((d) => d.playerNum === playerNum);
	}

	getDesign(playerNum: number, num: number): ShipDesign | undefined {
		return this.designs.find((d) => d.playerNum === playerNum && d.num === num);
	}

	getMyDesignByName(playerNum: number, name: string): ShipDesign | undefined {
		return this.designs.find((d) => d.playerNum === playerNum && d.name == name);
	}

	getBattle(num: number): BattleRecord | undefined {
		return this.battles.find((b) => b.num === num);
	}

	updateDesign(design: ShipDesign) {
		const index = this.designs.findIndex((f) => f.num === design.num);
		if (index != -1) {
			this.designs = [...this.designs.slice(0, index), design, ...this.designs.slice(index + 1)];
		}
	}

	addDesign(design: ShipDesign) {
		this.designs = [...this.designs, design];
	}

	getBattleLocation(battle: BattleRecord, universe: Universe): string {
		if (battle.planetNum) {
			const planet = universe.getPlanet(battle.planetNum);
			return planet?.name ?? 'Unknown';
		}
		return `Space (${battle.position.x}, ${battle.position.y}`;
	}

	getOtherMapObjectsHereByType(position: Vector) {
		return groupBy(this.mapObjectsByPosition[positionKey(position)], (mo) => mo.type);
	}

	getMapObjectsByPosition(position: MapObject | Vector) {
		return this.mapObjectsByPosition[positionKey(position)];
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

	getPlanetStarbase(planetNum: number) {
		return this.starbases.find((sb) => sb.planetNum == planetNum);
	}

	getWormhole(num: number) {
		return this.wormholes.find((w) => w.num === num);
	}

	getMysteryTrader(num: number) {
		return this.mysteryTraders.find((mt) => mt.num === num);
	}

	addFleets(fleets: Fleet[]) {
		this.fleets = [...fleets, ...this.fleets];
		this.resetMyMapObjectsByPosition();
	}

	updateFleet(fleet: Fleet) {
		const index = this.fleets.findIndex((f) => f.num === fleet.num);
		if (index != -1) {
			this.fleets = [...this.fleets.slice(0, index), fleet, ...this.fleets.slice(index + 1)];
		}
		this.resetMyMapObjectsByPosition();
	}

	updatePlanet(planet: Planet) {
		const index = this.planets.findIndex((f) => f.num === planet.num);
		if (index != -1) {
			this.planets = [...this.planets.slice(0, index), planet, ...this.planets.slice(index + 1)];
		}
		this.resetMyMapObjectsByPosition();
	}

	removeFleets(fleetNums: number[]) {
		this.fleets = this.fleets.filter((f) => fleetNums.indexOf(f.num) == -1);
		this.resetMyMapObjectsByPosition();
	}

	getTargetName(wp: Waypoint): string {
		if (wp.targetName && wp.targetName !== '') {
			return wp.targetName;
		}
		const mo = this.getMapObject(wp);
		if (mo) {
			if (mo.name && mo.name !== '') {
				return mo.name;
			}

			return `${startCase(mo.type)} #${mo.num}`;
		}
		return `Space: (${wp.position.x}, ${wp.position.y})`;
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
					(mf) => mf.num === target.targetNum && mf.playerNum === target.targetPlayerNum
				);
			case MapObjectType.Salvage:
				return this.salvages.find(
					(mf) => mf.num === target.targetNum && mf.playerNum === target.targetPlayerNum
				);
			case MapObjectType.Wormhole:
				return target.targetNum ? this.getWormhole(target.targetNum) : undefined;
			case MapObjectType.MysteryTrader:
				return target.targetNum ? this.getMysteryTrader(target.targetNum) : undefined;
			case MapObjectType.PositionWaypoint:
				break;
		}
	}

	getHomeworld(playerNum: number) {
		return this.planets.find((p) => p.playerNum === playerNum && p.homeworld);
	}
}
