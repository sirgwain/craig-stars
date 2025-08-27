import { battlesSortBy, getBattleRecordDetails, type BattleRecordDetails } from '$lib/types/Battle';
import type { CargoDest } from '$lib/types/CargoTransferRequest.svelte';
import type {
	BattleRecord,
	Fleet,
	Minefield,
	MineralPacket,
	PlayerScore,
	PlayerUniverse,
	ShipDesign,
	Waypoint
} from '$lib/types/cs-proto';
import {
	MapObjectType,
	VectorSchema,
	type MysteryTrader,
	type Planet,
	type PlayerIntel,
	type Salvage,
	type ScoreIntel,
	type Wormhole
} from '$lib/types/cs-proto';
import { enumToString } from '$lib/types/Enums';
import { fleetsSortBy } from '$lib/types/Fleet';
import {
	commandable,
	positionKey,
	type MapObjectLike,
	type MapObjectTargetLike,
	type Position
} from '$lib/types/MapObject';
import { planetsSortBy } from '$lib/types/Planet';
import type { CommandedPlayer } from '$lib/types/Player';
import { create, type UnknownField } from '@bufbuild/protobuf';
import { groupBy } from 'lodash-es';

export type PlayerDesigns = {
	designs: ShipDesign[];
};

export interface DesignFinder {
	getDesign(playerNum: number | undefined, num: number | undefined): ShipDesign | undefined;
	getMyDesign(num: number | undefined): ShipDesign | undefined;
}

export interface PlayerFinder {
	getPlayerIntel(num: number | undefined): PlayerIntel | undefined;
	getPlayerName(playerNum: number | undefined): string;
	getPlayerPluralName(playerNum: number | undefined): string;
	getPlayerColor(playerNum: number | undefined): string;
}
const sortByNum = (a: MapObjectLike, b: MapObjectLike) =>
	(a.mapObject?.num ?? 0) - (b.mapObject?.num ?? 0);

function addtoDict(mo: MapObjectLike, dict: Record<string, MapObjectLike[]>) {
	const key = positionKey(mo);
	if (!dict[key]) {
		dict[key] = [];
	}
	dict[key].push(mo);
}

export class Universe implements PlayerUniverse, DesignFinder {
	$typeName: 'craig_stars.v1.PlayerUniverse';
	$unknown?: UnknownField[] | undefined;
	playerNum = 0;
	planets: Planet[] = [];
	fleets: Fleet[] = [];
	starbases: Fleet[] = [];
	minefields: Minefield[] = [];
	mineralPackets: MineralPacket[] = [];
	designs: ShipDesign[] = [];
	mysteryTraders: MysteryTrader[] = [];
	wormholes: Wormhole[] = [];
	salvages: Salvage[] = [];

	battleRecords: BattleRecord[] = [];
	playerIntels: PlayerIntel[] = [];
	scoreIntels: ScoreIntel[] = [];
	planetIntels: Planet[] = [];
	fleetIntels: Fleet[] = [];
	shipDesignIntels: ShipDesign[] = [];
	minefieldIntels: Minefield[] = [];

	mapObjectsByPosition: Record<string, MapObjectLike[]> = {};
	myMapObjectsByPosition: Record<string, MapObjectLike[]> = {};
	allPlanets: Planet[] = [];

	constructor() {
		this.$typeName = 'craig_stars.v1.PlayerUniverse';
	}

	public get allFleets(): Fleet[] {
		return [...this.fleets, ...this.fleetIntels];
	}

	public get allMinefields(): Minefield[] {
		return [...this.minefields, ...this.minefieldIntels];
	}

	public get allDesigns(): ShipDesign[] {
		return [...this.designs, ...this.shipDesignIntels];
	}

	public get intels() {
		return {
			battleRecords: this.battleRecords,
			playerIntels: this.playerIntels,
			scoreIntels: this.scoreIntels,
			planetIntels: this.planetIntels,
			fleetIntels: this.fleetIntels,
			shipDesignIntels: this.shipDesignIntels,
			mineralPacketIntels: this.mineralPackets,
			minefieldIntels: this.minefieldIntels,
			wormholeIntels: this.wormholes,
			mysteryTraderIntels: this.mysteryTraders,
			salvageIntels: this.salvages
		};
	}

	public setData(playerNum: number, playerUniverse: PlayerUniverse): Universe {
		this.playerNum = playerNum;
		this.planets =
			playerUniverse.planets.filter((mo) => mo.mapObject?.playerNum === playerNum) ?? [];
		this.fleets =
			playerUniverse.fleets.filter((f) => f.mapObject?.playerNum === playerNum && !f.starbase) ??
			[];
		this.starbases =
			playerUniverse.fleets.filter((f) => f.mapObject?.playerNum === playerNum && f.starbase) ?? [];
		this.minefields =
			playerUniverse.minefields.filter((mo) => mo.mapObject?.playerNum === playerNum) ?? [];
		this.designs = playerUniverse.designs.filter((d) => d.playerNum === playerNum) ?? [];
		this.mineralPackets = playerUniverse.mineralPackets;

		// set player intel (now embedded under player.intels)
		this.battleRecords = playerUniverse?.battleRecords;
		this.playerIntels = playerUniverse?.playerIntels;
		this.scoreIntels = playerUniverse?.scoreIntels;
		this.planetIntels = playerUniverse.planets;
		this.fleetIntels =
			playerUniverse.fleets.filter((f) => f.mapObject && f.mapObject?.playerNum !== playerNum && !f.starbase) ??
			[];
		this.minefieldIntels =
			playerUniverse.minefields.filter((mo) => mo.mapObject?.playerNum !== playerNum) ?? [];
		this.shipDesignIntels = playerUniverse.designs.filter((d) => d.playerNum !== playerNum) ?? [];
		this.wormholes = playerUniverse.wormholes;
		this.mysteryTraders = playerUniverse.mysteryTraders;
		this.salvages = playerUniverse.salvages;

		this.resetAllPlanets();
		this.resetMapObjectsByPosition();
		this.resetMyMapObjectsByPosition();
		return this;
	}

	public resetAllPlanets() {
		this.allPlanets = [...this.planetIntels];
		this.planets.forEach((planet) => (this.allPlanets[(planet.mapObject?.num ?? 1) - 1] = planet));
	}

	resetMapObjectsByPosition() {
		this.mapObjectsByPosition = {};

		// add all planets, both intel and owned
		this.allPlanets.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));

		// add all owned mapobjects
		this.fleets.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.minefields.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));

		// add all intel
		this.fleetIntels.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.minefieldIntels.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.salvages.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.wormholes.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.mysteryTraders.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
		this.mineralPackets.forEach((mo) => addtoDict(mo, this.mapObjectsByPosition));
	}

	resetMyMapObjectsByPosition() {
		// build a map of objects owned by me
		this.myMapObjectsByPosition = {};
		this.planets.sort(sortByNum).forEach((mo) => addtoDict(mo, this.myMapObjectsByPosition));
		this.fleets.sort(sortByNum).forEach((mo) => addtoDict(mo, this.myMapObjectsByPosition));
		this.minefields.sort(sortByNum).forEach((mo) => addtoDict(mo, this.myMapObjectsByPosition));
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
			if (history && history.length > 0) {
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
			return intel.color ?? '#FF0000';
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

	getPlanets(sortKey: string, descending: boolean): Planet[] {
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
	getAllFleets(sortKey?: string, descending?: boolean): Fleet[] {
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

	getDesign(playerNum: number | undefined, num: number | undefined): ShipDesign | undefined {
		if (!playerNum || !num) return;
		if (playerNum === this.playerNum) {
			return this.designs.find((d) => d.num === num);
		}
		return this.shipDesignIntels.find((d) => d.playerNum === playerNum && d.num === num);
	}

	getDesigns(playerNum: number): ShipDesign[] {
		return this.allDesigns.filter((d) => d.playerNum === playerNum);
	}

	getMyDesign(num: number | undefined): ShipDesign | undefined {
		return this.designs.find((d) => d.playerNum === this.playerNum && d.num === num);
	}

	getBattle(num: number | undefined): BattleRecord | undefined {
		return this.battleRecords.find((b) => b.num === num);
	}

	validateDesign(design: ShipDesign): { valid: boolean; reason?: string } {
		// if we have a design with this name already, it is invalid
		const designsWithName = this.getMyDesigns().filter(
			(d) => d.gameDbObject?.id !== design.gameDbObject?.id && d.name === design.name
		);
		if (designsWithName.length > 0) {
			return { valid: false, reason: `Another design named ${design.name} exists` };
		}
		return { valid: true };
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
			return planet?.mapObject?.name ?? 'Unknown';
		}
		return `Space (${battle.position?.x ?? 0}, ${battle.position?.y ?? 0})`;
	}

	getOtherMapObjectsHereByType(position: Position) {
		return groupBy(this.mapObjectsByPosition[positionKey(position)], (mo) => mo.mapObject?.type);
	}

	getMapObjectsByPosition(position: Position) {
		return this.mapObjectsByPosition[positionKey(position ?? create(VectorSchema))];
	}

	getCargoDestsByPosition(position: Position): CargoDest[] {
		return this.mapObjectsByPosition[positionKey(position)]
			?.filter(
				(mo) =>
					[
						MapObjectType.FLEET,
						MapObjectType.MINERAL_PACKET,
						MapObjectType.SALVAGE,
						MapObjectType.PLANET
					].indexOf(mo.mapObject?.type ?? MapObjectType.UNSPECIFIED) != -1
			)
			.map((mo) => mo as CargoDest);
	}

	getSalvageAtPosition(position: Position): Salvage | undefined {
		const mo = this.getMapObjectsByPosition(position)?.find(
			(mo) => mo.mapObject?.type === MapObjectType.SALVAGE
		);
		if (mo) {
			return mo as Salvage;
		}
	}

	getSalvage(num: number | undefined): Salvage | undefined {
		return this.salvages.find((s) => s.mapObject?.num === num);
	}

	getMyMapObjectsByPosition(position: Position) {
		return this.myMapObjectsByPosition[positionKey(position)];
	}

	getCommandableMapObjectsByPosition(position: Position) {
		return (
			this.myMapObjectsByPosition[positionKey(position)]?.filter((mo) =>
				commandable(this.playerNum, mo)
			) ?? []
		);
	}

	getMyFleetsByPosition(position: Position): Fleet[] {
		return (
			(this.getMyMapObjectsByPosition(position)?.filter(
				(mo) => mo.mapObject?.type === MapObjectType.FLEET
			) as Fleet[]) ?? []
		);
	}

	getFleetsByPosition(position: Position): Fleet[] {
		return (
			(this.getMapObjectsByPosition(position)?.filter(
				(mo) => mo.mapObject?.type === MapObjectType.FLEET
			) as Fleet[]) ?? []
		);
	}

	// getPlanet returns either the player owned planet by a number
	getPlanet(num: number | undefined): Planet | undefined {
		if (!num || num > this.allPlanets.length) return;
		return this.allPlanets[num - 1];
	}

	getFleet(playerNum: number | undefined, num: number | undefined): Fleet | undefined {
		return this.allFleets.find(
			(f) => f.mapObject?.playerNum === playerNum && f.mapObject?.num === num
		);
	}

	getMyFleet(num: number | undefined): Fleet | undefined {
		return this.fleets.find((f) => f.mapObject?.num === num);
	}

	getMyPlanetStarbase(planetNum: number) {
		return this.starbases.find((sb) => sb.planetNum === planetNum);
	}

	getWormhole(num: number) {
		return this.wormholes.find((w) => w.mapObject?.num === num);
	}

	getMysteryTrader(num: number) {
		return this.mysteryTraders.find((mt) => mt.mapObject?.num === num);
	}

	getMinefield(playerNum: number | undefined, num: number | undefined) {
		return this.minefields.find(
			(f) => f.mapObject?.playerNum === playerNum && f.mapObject?.num === num
		);
	}

	getMineralPacket(playerNum: number | undefined, num: number | undefined) {
		return this.mineralPackets.find(
			(f) => f.mapObject?.playerNum === playerNum && f.mapObject?.num === num
		);
	}

	addFleets(fleets: Fleet[]) {
		this.fleets = [...fleets, ...this.fleets];
		this.resetMapObjectsByPosition();
		this.resetMyMapObjectsByPosition();
	}

	updateFleet(fleet: Fleet | undefined) {
		if (!fleet) {
			return;
		}
		const index = this.fleets.findIndex(
			(f) =>
				f.mapObject?.num === fleet.mapObject?.num &&
				f.mapObject?.playerNum === fleet.mapObject?.playerNum
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
			(p) =>
				p.mapObject?.num === planet.mapObject?.num &&
				p.mapObject?.playerNum === planet.mapObject?.playerNum
		);
		if (index != -1) {
			this.planets = [...this.planets.slice(0, index), planet, ...this.planets.slice(index + 1)];
		}
		// update intel as well
		this.planetIntels[(planet.mapObject?.num ?? 1) - 1] = planet;
		this.allPlanets = [...this.planetIntels];
		this.planets.forEach((planet) => (this.allPlanets[(planet.mapObject?.num ?? 1) - 1] = planet));

		this.resetMapObjectsByPosition();
		this.resetMyMapObjectsByPosition();
	}

	updateMinefield(minefield: Minefield) {
		const index = this.minefields.findIndex(
			(mf) =>
				mf.mapObject?.playerNum === minefield.mapObject?.playerNum &&
				mf.mapObject?.num === minefield.mapObject?.num
		);
		if (index != -1) {
			this.minefields = [
				...this.minefields.slice(0, index),
				minefield,
				...this.minefields.slice(index + 1)
			];
		}
		this.resetMapObjectsByPosition();
		this.resetMyMapObjectsByPosition();
	}

	updateMineralPacket(mineralPacket: MineralPacket) {
		const index = this.mineralPackets.findIndex(
			(mf) =>
				mf.mapObject?.playerNum === mineralPacket.mapObject?.playerNum &&
				mf.mapObject?.num === mineralPacket.mapObject?.num
		);
		if (index != -1) {
			this.mineralPackets = [
				...this.mineralPackets.slice(0, index),
				mineralPacket as MineralPacket,
				...this.mineralPackets.slice(index + 1)
			];
		}
		this.resetMapObjectsByPosition();
		this.resetMyMapObjectsByPosition();
	}

	updateSalvage(salvage: Salvage) {
		const index = this.salvages.findIndex(
			(mf) =>
				mf.mapObject?.playerNum === salvage.mapObject?.playerNum &&
				mf.mapObject?.num === salvage.mapObject?.num
		);
		if (index != -1) {
			this.salvages = [
				...this.salvages.slice(0, index),
				salvage,
				...this.salvages.slice(index + 1)
			];
		}
		this.resetMapObjectsByPosition();
	}

	updateSalvages(salvages: Salvage[]) {
		this.salvages = salvages;
		this.resetMapObjectsByPosition();
		this.resetMyMapObjectsByPosition();
	}

	removeFleets(fleetNums: number[]) {
		this.fleets = this.fleets.filter((f) => fleetNums.indexOf(f.mapObject?.num ?? 0) == -1);
		this.resetMapObjectsByPosition();
		this.resetMyMapObjectsByPosition();
	}

	getTargetName(wp: Waypoint): string {
		// first see if we can load this waypoint target as a mapobject
		// if so, use its name
		const mo = this.getMapObject(wp.mapObjectTarget);
		if (mo) {
			if (mo.mapObject?.name && mo.mapObject.name !== '') {
				return mo.mapObject.name;
			}

			return `${enumToString(MapObjectType, mo.mapObject?.type ?? MapObjectType.UNSPECIFIED)} #${mo.mapObject?.num ?? ''}`;
		} else if (wp.mapObjectTarget?.targetName && wp.mapObjectTarget?.targetName !== '') {
			// we can't load it from the universe, see if the server gave us a target name
			return wp.mapObjectTarget?.targetName;
		}

		const pos = wp.position ?? create(VectorSchema);
		// we don't have a target name and we can't find the map object, just point it to the space location
		return `Space: (${pos.x.toFixed()}, ${pos.y.toFixed()})`;
	}

	// get a mapobject by type, number, and optionally player num
	getMapObject(target: MapObjectTargetLike | undefined): MapObjectLike | undefined {
		if (!target) {
			return;
		}
		switch (target.targetType) {
			case MapObjectType.PLANET:
				return target.targetNum ? this.getPlanet(target.targetNum) : undefined;
			case MapObjectType.FLEET:
				return this.allFleets.find(
					(f) =>
						f.mapObject?.num === target.targetNum &&
						f.mapObject?.playerNum === target.targetPlayerNum
				);
			case MapObjectType.MINEFIELD:
				return this.allMinefields.find(
					(mf) =>
						mf.mapObject?.num === target.targetNum &&
						mf.mapObject?.playerNum === target.targetPlayerNum
				);
			case MapObjectType.MINERAL_PACKET:
				return this.mineralPackets.find(
					(p) =>
						p.mapObject?.num === target.targetNum &&
						p.mapObject?.playerNum === target.targetPlayerNum
				);
			case MapObjectType.SALVAGE:
				return this.salvages.find(
					(s) =>
						s.mapObject?.num === target.targetNum &&
						s.mapObject?.playerNum === target.targetPlayerNum
				);
			case MapObjectType.WORMHOLE:
				return target.targetNum ? this.getWormhole(target.targetNum) : undefined;
			case MapObjectType.MYSTERY_TRADER:
				return target.targetNum ? this.getMysteryTrader(target.targetNum) : undefined;
		}
	}

	getHomeworld() {
		return this.planets.find((p) => p.homeworld);
	}
}
