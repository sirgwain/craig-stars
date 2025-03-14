import { battlesSortBy, getBattleRecordDetails, type BattleRecordDetails } from '$lib/types/Battle';
import type {
	Cost,
	FleetIntel,
	Intel,
	MapObjectTarget,
	MineField,
	MineFieldIntel,
	MineralPacket,
	MineralPacketIntel,
	MysteryTraderIntel,
	PlanetIntel,
	PlayerIntel,
	PlayerIntels,
	PlayerScore,
	ProductionQueueItem,
	SalvageIntel,
	ScoreIntel,
	ShipDesign,
	ShipDesignIntel,
	Vector,
	WormholeIntel
} from '$lib/types/cs';
import {
	MapObjectTypeFleet,
	MapObjectTypeMineField,
	MapObjectTypeMineralPacket,
	MapObjectTypeMysteryTrader,
	MapObjectTypePlanet,
	MapObjectTypeSalvage,
	MapObjectTypeWormhole,
	type BattleRecord,
	type Fleet,
	type MapObject,
	type Planet,
	type Waypoint
} from '$lib/types/cs';
import { fleetsSortBy } from '$lib/types/Fleet';
import { CommandedPlanet, planetsSortBy } from '$lib/types/Planet';
import type { CommandedPlayer } from '$lib/types/Player';
import type { CS } from '$lib/wasm';
import { groupBy, startCase } from 'lodash-es';

export type AnyPlanet = Planet | PlanetIntel;
export type AnyFleet = Fleet | FleetIntel;
export type AnyMineField = MineField | MineFieldIntel;
export type AnyMineralPacket = MineralPacket | MineralPacketIntel;
export type AnyShipDesign = ShipDesign | ShipDesignIntel;

/**
 * Check whether a given object is an {@linkcode Intel} or not.
 * @param obj the object to check
 * @returns `true` if obj implements {@linkcode Intel}.
 */
export function isIntel(obj: unknown): obj is Intel {
	if (!obj || typeof obj !== 'object') {
		return false
	}
	return 'reportAge' in obj && typeof obj.reportAge === "number"
}

export type PlayerUniverse = {
	planets: Planet[];
	fleets: Fleet[];
	starbases: Fleet[];
	mineFields: MineField[];
	mineralPackets: MineralPacket[];
	designs: ShipDesign[];
} & PlayerIntels;

export interface DesignFinder {
	getDesign(playerNum: number, num: number): AnyShipDesign | undefined;
	getMyDesign(num: number | undefined): ShipDesign | undefined;
}

export interface CostFinder {
	getItemCost(
		cs: CS,
		item: ProductionQueueItem | undefined,
		designFinder: DesignFinder,
		planet?: CommandedPlanet,
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
	// eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
	if (!dict[key]) { // This is weird due to golang jank with maps
		dict[key] = [];
	}
	dict[key].push(mo);
}

function positionKey(pos: MapObject | Vector): string {
	if ('position' in pos) {
		return `${pos.position.x},${pos.position.y}`;
	}
	return `${pos.x},${pos.y}`;
}

export class Universe implements PlayerUniverse, DesignFinder {
	playerNum = 0;
	planets: Planet[] = [];
	fleets: Fleet[] = [];
	starbases: Fleet[] = [];
	mineFields: MineField[] = [];
	mineralPackets: MineralPacket[] = [];
	designs: ShipDesign[] = [];

	battleRecords: BattleRecord[] = [];
	playerIntels: PlayerIntel[] = [];
	scoreIntels: ScoreIntel[] = [];
	planetIntels: PlanetIntel[] = [];
	fleetIntels: FleetIntel[] = [];
	shipDesignIntels: ShipDesignIntel[] = [];
	mineralPacketIntels: MineralPacketIntel[] = [];
	mineFieldIntels: MineFieldIntel[] = [];
	wormholeIntels: WormholeIntel[] = [];
	mysteryTraderIntels: MysteryTraderIntel[] = [];
	salvageIntels: SalvageIntel[] = [];

	mapObjectsByPosition: Record<string, MapObject[]> = {};
	myMapObjectsByPosition: Record<string, MapObject[]> = {};
	allPlanets: AnyPlanet[] = [];

	public get allFleets(): AnyFleet[] {
		return [...this.fleets, ...this.fleetIntels];
	}

	public get allMineFields(): AnyMineField[] {
		return [...this.mineFields, ...this.mineFieldIntels];
	}

	public get allMineralPackets(): AnyMineralPacket[] {
		return [...this.mineralPackets, ...this.mineralPacketIntels];
	}

	public get allDesigns(): AnyShipDesign[] {
		return [...this.designs, ...this.shipDesignIntels];
	}

	public setData(data: PlayerUniverse): Universe {
		this.planets = data.planets
		this.fleets = data.fleets
		this.starbases = data.starbases
		this.mineFields = data.mineFields
		this.mineralPackets = data.mineralPackets
		this.designs = data.designs

		this.battleRecords = data.battleRecords ?? [];
		this.playerIntels = data.playerIntels ?? [];
		this.scoreIntels = data.scoreIntels ?? [];
		this.planetIntels = data.planetIntels ?? [];
		this.fleetIntels = data.fleetIntels ?? [];
		this.shipDesignIntels = data.shipDesignIntels ?? [];
		this.mineralPacketIntels = data.mineralPacketIntels ?? [];
		this.mineFieldIntels = data.mineFieldIntels ?? [];
		this.wormholeIntels = data.wormholeIntels ?? [];
		this.mysteryTraderIntels = data.mysteryTraderIntels ?? [];
		this.salvageIntels = data.salvageIntels ?? [];

		this.allPlanets = [...this.planetIntels];
		this.planets.forEach((planet) => (this.allPlanets[planet.num - 1] = planet));

		this.resetMapObjectsByPosition();
		return this;
	}

	// set the player number and update player specific internal data
	public setPlayer(playerNum: number) {
		this.playerNum = playerNum;
		this.resetMyMapObjectsByPosition();
	}

	// reset all data in the universe
	public resetData(playerNum: number, data: PlayerUniverse): Universe {
		this.setData(data);
		this.setPlayer(playerNum);
		return this;
	}

	resetMapObjectsByPosition() {
		this.mapObjectsByPosition = {};

		// add all planets, both intel and owned
		this.allPlanets.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));

		// add all owned mapobjects
		this.fleets.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.mineFields.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.mineralPackets.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));

		// add all intel
		this.fleetIntels.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.mineFieldIntels.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.mineralPacketIntels.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.salvageIntels.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.wormholeIntels.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.mysteryTraderIntels.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
	}

	resetMyMapObjectsByPosition() {
		// build a map of objects owned by me
		this.myMapObjectsByPosition = {};
		this.planets.sort(sortByNum).forEach((mo) => addtoDict(mo, this.myMapObjectsByPosition));
		this.fleets.sort(sortByNum).forEach((mo) => addtoDict(mo, this.myMapObjectsByPosition));
		this.mineFields.sort(sortByNum).forEach((mo) => addtoDict(mo, this.myMapObjectsByPosition));
		this.mineralPackets.sort(sortByNum).forEach((mo) => addtoDict(mo, this.myMapObjectsByPosition));
	}

	getPlayerIntel(num: number): PlayerIntel | undefined {
		if (num >= 1 && num <= this.playerIntels.length) {
			return this.playerIntels[num - 1];
		}
	}

	getPlayerScoreHistory(num: number): PlayerScore[] | undefined {
		if (num >= 1 && num <= this.scoreIntels.length && this.scoreIntels[num - 1]) {
			return this.scoreIntels[num - 1].scoreHistory;
		}
	}

	getPlayerScore(num: number): PlayerScore | undefined {
		if (num >= 1 && num <= this.scoreIntels.length) {
			const history = this.scoreIntels[num - 1].scoreHistory;
			if (history.length > 0) {
				return history[history.length - 1];
			}
		}
	}

	getPlayerName(playerNum: number | undefined): string {
		if (playerNum && playerNum > 0 && playerNum <= this.playerIntels.length) {
			const intel = this.playerIntels[playerNum - 1];
			return intel.raceName ?? intel.name;
		}
		return 'unknown';
	}

	getPlayerPluralName(playerNum: number | undefined): string {
		if (playerNum && playerNum > 0 && playerNum <= this.playerIntels.length) {
			const intel = this.playerIntels[playerNum - 1];
			return intel.racePluralName ?? intel.name;
		}
		return 'unknown';
	}

	getPlayerColor(playerNum: number | undefined): string {
		if (playerNum && playerNum > 0 && playerNum <= this.playerIntels.length) {
			const intel = this.playerIntels[playerNum - 1];
			return intel.color || '#FF0000';
		}
		return '#FF0000';
	}

	getMyDesigns(): ShipDesign[] {
		return this.designs.filter((d) => d.playerNum === this.playerNum);
	}

	getMyPlanets(sortKey: string, descending: boolean): Planet[] {
		const planets = [...this.planets];

		planets.sort(planetsSortBy(sortKey));
		if (descending) {
			planets.reverse();
		}
		return planets;
	}

	getPlanets(sortKey: string, descending: boolean): PlanetIntel[] {
		const planets = [...this.planetIntels];
		planets.sort(planetsSortBy(sortKey));
		if (descending) {
			planets.reverse();
		}
		return planets;
	}

	getMyFleets(sortKey: string, descending: boolean): Fleet[] {
		const fleets = [...this.fleets];
		fleets.sort(fleetsSortBy(sortKey, this));
		if (descending) {
			fleets.reverse();
		}
		return fleets;
	}

	// getAllFleets gets all of my fleets and other player fleets, optionally sorted
	getAllFleets(sortKey?: string, descending?: boolean): AnyFleet[] {
		const fleets = [...this.fleets, ...this.fleetIntels];
		if (sortKey) {
			fleets.sort(fleetsSortBy(sortKey, this));
			if (descending) {
				fleets.reverse();
			}
		}
		return fleets;
	}

	getBattles(sortKey: string, descending: boolean, player: CommandedPlayer): BattleRecordDetails[] {
		const battles = this.battleRecords.map((b) => getBattleRecordDetails(b, player, this));
		battles.sort(battlesSortBy(sortKey));
		if (descending) {
			battles.reverse();
		}
		return battles;
	}

	getDesign(playerNum: number, num: number): AnyShipDesign | undefined {
		if (playerNum === this.playerNum) {
			return this.designs.find((d) => d.num === num);
		}
		return this.shipDesignIntels.find((d) => d.playerNum === playerNum && d.num === num);
	}

	getDesigns(playerNum: number): AnyShipDesign[] {
		return this.allDesigns.filter((d) => d.playerNum === playerNum);
	}

	getMyDesign(num: number | undefined): ShipDesign | undefined {
		return this.designs.find((d) => d.playerNum === this.playerNum && d.num === num);
	}

	getBattle(num: number | undefined): BattleRecord | undefined {
		return this.battleRecords.find((b) => b.num === num);
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

	getSalvageAtPosition(position: MapObject | Vector): SalvageIntel | undefined {
		const mo = this.getMapObjectsByPosition(position).find(
			(mo) => mo.type === MapObjectTypeSalvage
		);
		if (mo) {
			return mo as SalvageIntel;
		}
	}

	getSalvage(num: number | undefined): SalvageIntel | undefined {
		return this.salvageIntels.find((s) => s.num === num);
	}

	getMyMapObjectsByPosition(position: MapObject | Vector) {
		return this.myMapObjectsByPosition[positionKey(position)];
	}

	getMyFleetsByPosition(position: MapObject | Vector): Fleet[] {
		return this.getMyMapObjectsByPosition(position).filter(
				(mo) => mo.type === MapObjectTypeFleet
			) as Fleet[];
	}

	// getPlanet returns either the player owned planet by a number
	getPlanet(num: number): AnyPlanet | undefined {
		return this.allPlanets[num - 1];
	}

	getFleet(playerNum: number | undefined, num: number | undefined): AnyFleet | undefined {
		return this.allFleets.find((f) => f.playerNum === playerNum && f.num === num);
	}

	getMyFleet(num: number | undefined): Fleet | undefined {
		return this.fleets.find((f) => f.num === num);
	}

	getMyPlanetStarbase(planetNum: number) {
		return this.starbases.find((sb) => sb.planetNum === planetNum);
	}

	getWormhole(num: number) {
		return this.wormholeIntels.find((w) => w.num === num);
	}

	getMysteryTrader(num: number) {
		return this.mysteryTraderIntels.find((mt) => mt.num === num);
	}

	getMineField(playerNum: number | undefined, num: number | undefined) {
		return this.mineFields.find((f) => f.playerNum === playerNum && f.num === num);
	}

	getMineralPacket(playerNum: number | undefined, num: number | undefined) {
		return this.mineralPacketIntels.find((f) => f.playerNum === playerNum && f.num === num);
	}

	addFleets(fleets: Fleet[]) {
		this.fleets = [...fleets, ...this.fleets];
		this.resetMapObjectsByPosition();
		this.resetMyMapObjectsByPosition();
	}

	updateFleet(fleet: Fleet) {
		const index = this.fleets.findIndex(
			(f) => f.num === fleet.num && f.playerNum === fleet.playerNum
		);
		if (index != -1) {
			this.fleets = [...this.fleets.slice(0, index), fleet, ...this.fleets.slice(index + 1)];
		}
		this.resetMapObjectsByPosition();
		this.resetMyMapObjectsByPosition();
	}

	updatePlanet(planet: Planet) {
		// try and update our own planet first
		const index = this.planets.findIndex(
			(p) => p.num === planet.num && p.playerNum === planet.playerNum
		);
		if (index != -1) {
			this.planets = [...this.planets.slice(0, index), planet, ...this.planets.slice(index + 1)];
		}
		// update intel as well
		this.planetIntels[planet.num - 1] = { ...planet, reportAge: 0 };

		this.allPlanets = [...this.planetIntels];
		this.planets.forEach((planet) => (this.allPlanets[planet.num - 1] = planet));

		this.resetMapObjectsByPosition();
		this.resetMyMapObjectsByPosition();
	}

	updateMineField(mineField: MineField) {
		const index = this.mineFields.findIndex(
			(mf) => mf.playerNum === mineField.playerNum && mf.num === mineField.num
		);
		if (index != -1) {
			this.mineFields = [
				...this.mineFields.slice(0, index),
				mineField,
				...this.mineFields.slice(index + 1)
			];
		}
		this.resetMapObjectsByPosition();
		this.resetMyMapObjectsByPosition();
	}

	updateMineralPacket(mineralPacket: AnyMineralPacket) {
		if (mineralPacket.playerNum === this.playerNum) {
			const index = this.mineralPackets.findIndex(
				(mf) => mf.playerNum === mineralPacket.playerNum && mf.num === mineralPacket.num
			);
			if (index != -1) {
				this.mineralPackets = [
					...this.mineralPackets.slice(0, index),
					mineralPacket as MineralPacket,
					...this.mineralPackets.slice(index + 1)
				];
			}
		} else {
			const index = this.mineralPacketIntels.findIndex(
				(mf) => mf.playerNum === mineralPacket.playerNum && mf.num === mineralPacket.num
			);
			if (index != -1) {
				this.mineralPacketIntels = [
					...this.mineralPacketIntels.slice(0, index),
					mineralPacket as MineralPacketIntel,
					...this.mineralPacketIntels.slice(index + 1)
				];
			}
		}
		this.resetMapObjectsByPosition();
		this.resetMyMapObjectsByPosition();
	}

	updateSalvages(salvages: SalvageIntel[]) {
		this.salvageIntels = salvages;
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
	getMapObject(target: MapObjectTarget): MapObject | undefined {
		switch (target.targetType) {
			case MapObjectTypePlanet:
				return target.targetNum ? this.getPlanet(target.targetNum) : undefined;
			case MapObjectTypeFleet:
				return this.allFleets.find(
					(f) => f.num === target.targetNum && f.playerNum === target.targetPlayerNum
				);
			case MapObjectTypeMineField:
				return this.allMineFields.find(
					(mf) => mf.num === target.targetNum && mf.playerNum === target.targetPlayerNum
				);
			case MapObjectTypeMineralPacket:
				return this.allMineralPackets.find(
					(p) => p.num === target.targetNum && p.playerNum === target.targetPlayerNum
				);
			case MapObjectTypeSalvage:
				return this.salvageIntels.find(
					(s) => s.num === target.targetNum && s.playerNum === target.targetPlayerNum
				);
			case MapObjectTypeWormhole:
				return target.targetNum ? this.getWormhole(target.targetNum) : undefined;
			case MapObjectTypeMysteryTrader:
				return target.targetNum ? this.getMysteryTrader(target.targetNum) : undefined;
		}
	}

	getHomeworld() {
		return this.planets.find((p) => p.homeworld);
	}
}
