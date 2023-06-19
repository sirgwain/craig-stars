import type { Fleet, Target } from '$lib/types/Fleet';
import { equal, MapObjectType, type MapObject } from '$lib/types/MapObject';
import type { MineField } from '$lib/types/MineField';
import type { MineralPacket } from '$lib/types/MineralPacket';
import type { MysteryTrader } from '$lib/types/MysteryTrader';
import type { Planet } from '$lib/types/Planet';
import type { PlayerIntel, PlayerIntels, PlayerMapObjects } from '$lib/types/Player';
import type { Salvage } from '$lib/types/Salvage';
import type { ShipDesignIntel } from '$lib/types/ShipDesign';
import type { Vector } from '$lib/types/Vector';
import type { Wormhole } from '$lib/types/Wormhole';
import { groupBy } from 'lodash-es';
import { commandMapObject, selectMapObject, zoomToMapObject } from './Context';

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

export class Universe implements PlayerMapObjects, PlayerIntels {
	playerNum = 0;
	planets: Planet[] = [];
	fleets: Fleet[] = [];
	salvages: Salvage[] = [];
	mineFields: MineField[] = [];
	mineralPackets: MineralPacket[] = [];
	starbases: Fleet[] = [];
	planetIntels: Planet[] = [];
	fleetIntels: Fleet[] = [];
	mineFieldIntels: MineField[] = [];
	mineralPacketIntels: MineralPacket[] = [];
	salvageIntels: Salvage[] = [];
	wormholeIntels: Wormhole[] = [];
	mysteryTraderIntels: MysteryTrader[] = [];
	shipDesignIntels: ShipDesignIntel[] = [];
	playerIntels: PlayerIntel[] = [];

	mapObjectsByPosition: Record<string, MapObject[]> = {};
	myMapObjectsByPosition: Record<string, MapObject[]> = {};

	constructor(data?: PlayerMapObjects & PlayerIntels) {
		Object.assign(this, data);
	}

	setIntels(intels: PlayerIntels) {
		this.mapObjectsByPosition = {};

		this.planetIntels = intels.planetIntels ?? [];
		this.fleetIntels = intels.fleetIntels ?? [];
		this.mineFieldIntels = intels.mineFieldIntels ?? [];
		this.mineralPacketIntels = intels.mineralPacketIntels ?? [];
		this.salvageIntels = intels.salvageIntels ?? [];
		this.wormholeIntels = intels.wormholeIntels ?? [];
		this.mysteryTraderIntels = intels.mysteryTraderIntels ?? [];
		this.planetIntels.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.fleetIntels.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.mineFieldIntels.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.mineralPacketIntels.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.salvageIntels.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.wormholeIntels.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.mysteryTraderIntels.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.fleets.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.mineFields.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.mineralPackets.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
	}

	setMapObjects(mos: PlayerMapObjects) {
		Object.assign(this, mos);

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
		return this.planets.find((p) => p.num === num) ?? this.planetIntels.find((p) => p.num === num);
	}

	getPlanetStarbase(planetNum: number) {
		return this.starbases.find((sb) => sb.planetNum == planetNum);
	}

	getWormhole(num: number) {
		return this.wormholeIntels.find((w) => w.num === num);
	}

	getMysteryTrader(num: number) {
		return this.mysteryTraderIntels.find((mt) => mt.num === num);
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

	getIntelMapObject(mo: MapObject): MapObject {
		switch (mo.type) {
			case MapObjectType.Planet:
				return this.planetIntels.find((planet) => planet.num == mo.num) ?? mo;
			case MapObjectType.Fleet:
				return this.fleetIntels.find((fleet) => fleet.num == mo.num) ?? mo;
			default:
				return mo;
		}
	}

	// get a mapobject by type, number, and optionally player num
	getMapObject(target: Target): MapObject | undefined {
		switch (target.targetType) {
			case MapObjectType.Planet:
				return target.targetNum ? this.getPlanet(target.targetNum) : undefined;
			case MapObjectType.Fleet:
				if (target.targetPlayerNum === this.playerNum) {
					return this.fleets.find((f) => f.num === target.targetNum);
				}
				return this.fleetIntels.find(
					(f) => f.num === target.targetNum && f.playerNum === target.targetPlayerNum
				);
			case MapObjectType.Wormhole:
				return target.targetNum ? this.getWormhole(target.targetNum) : undefined;
			case MapObjectType.MineField:
				if (target.targetPlayerNum === this.playerNum) {
					return this.mineFields.find((mf) => mf.num === target.targetNum);
				}
				return this.mineFieldIntels.find(
					(mf) => mf.num === target.targetNum && mf.playerNum === target.targetPlayerNum
				);
			case MapObjectType.MysteryTrader:
				return target.targetNum ? this.getMysteryTrader(target.targetNum) : undefined;
			case MapObjectType.Salvage:
				if (target.targetPlayerNum === this.playerNum) {
					return this.salvages.find((mf) => mf.num === target.targetNum);
				}
				return this.salvageIntels.find(
					(mf) => mf.num === target.targetNum && mf.playerNum === target.targetPlayerNum
				);
				break;
			case MapObjectType.MineralPacket:
				if (target.targetPlayerNum === this.playerNum) {
					return this.mineralPackets.find((mf) => mf.num === target.targetNum);
				}
				return this.mineralPacketIntels.find(
					(mf) => mf.num === target.targetNum && mf.playerNum === target.targetPlayerNum
				);
			case MapObjectType.PositionWaypoint:
				break;
		}
	}

	// command the player's homeworld (or the first planet they own, if their homeworld has been taken)
	commandHomeWorld() {
		const homeworld = this.planets.find((p) => p.homeworld);
		if (homeworld) {
			commandMapObject(homeworld);
			selectMapObject(homeworld);
			zoomToMapObject(homeworld);
		} else if (this.planets.length > 0) {
			commandMapObject(this.planets[0]);
			selectMapObject(this.planets[0]);
			zoomToMapObject(this.planets[0]);
		}
	}
}
