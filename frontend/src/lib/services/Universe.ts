import type { Fleet } from '$lib/types/Fleet';
import { MapObjectType, type MapObject } from '$lib/types/MapObject';
import type { MineField } from '$lib/types/MineField';
import type { MineralPacket } from '$lib/types/MineralPacket';
import type { Planet } from '$lib/types/Planet';
import type { PlayerIntel, PlayerIntels, PlayerMapObjects } from '$lib/types/Player';
import type { ShipDesignIntel } from '$lib/types/ShipDesign';
import type { Vector } from '$lib/types/Vector';
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
	mineFields: MineField[] = [];
	mineralPackets: MineralPacket[] = [];
	starbases: Fleet[] = [];
	planetIntels: Planet[] = [];
	fleetIntels: Fleet[] = [];
	mineFieldIntels: MineField[] = [];
	mineralPacketIntels: MineralPacket[] = [];
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
		this.planetIntels.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.fleetIntels.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
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
	getMapObject(type: MapObjectType, num: number, playerNum?: number): MapObject | undefined {
		let mo: MapObject;
		switch (type) {
			case MapObjectType.Planet:
				mo = this.planetIntels[num - 1];
				if (mo.playerNum === this.playerNum) {
					return this.getPlanet(mo.num);
				}
				return mo;
			case MapObjectType.Fleet:
				if (playerNum === this.playerNum) {
					return this.fleets.find((f) => f.num == num);
				}
				return this.fleetIntels.find((f) => f.num == num && f.playerNum == playerNum);
			case MapObjectType.Wormhole:
				break;
			case MapObjectType.MineField:
				if (playerNum === this.playerNum) {
					return this.mineFields.find((mf) => mf.num == num);
				}
				return this.mineFieldIntels.find((mf) => mf.num == num && mf.playerNum == playerNum);
			case MapObjectType.MysteryTrader:
				break;
			case MapObjectType.Salvage:
				break;
			case MapObjectType.MineralPacket:
				if (playerNum === this.playerNum) {
					return this.mineralPackets.find((mf) => mf.num == num);
				}
				return this.mineralPacketIntels.find((mf) => mf.num == num && mf.playerNum == playerNum);
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
