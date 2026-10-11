import type { PlayerMessageType, ShipClass } from './cs-proto';
import { MinefieldOwners, type FilterOptions, type MinefieldOwner } from './Filter';

export enum PlanetViewState {
	// I do enjoy the classics
	Normal,
	SurfaceMinerals,
	MineralConcentration,
	Percent,
	Population,

	/// Show a bunch of gray dots. How boring
	None
}

export class PlayerSettings implements FilterOptions {
	planetViewState = PlanetViewState.Normal;
	addWaypoint = false;
	fastestWaypoint = false;
	setPacketDest = false;
	setRouteDest = false;
	showNames = false;
	showFleetTokenCounts = false;
	showScanners = true;
	showAllyScanners = true;
	showMinefields = true;
	showMinefieldsOf: MinefieldOwner[] = [...MinefieldOwners];
	showMessagePane = false;
	showPlayerColors = false;
	scannerPercent = 100;
	mineralScale = 5000;
	messageTypeFilter = new Set<number>();
	messageTypeFilterArray: number[] = [];
	sortPlanetsKey = 'name';
	sortPlanetsDescending = false;
	showAllPlanets = false;
	sortFleetsKey = 'num';
	sortFleetsDescending = false;
	sortBattlesKey = 'num';
	sortBattlesDescending = false;
	searchQuery = '';

	// filters
	showIdleFleetsOnly = false;
	filterMyDesigns = false;
	filterAllyDesigns = false;
	filterEnemyDesigns = false;
	filterDesigns: number[] = [];
	filterAllyShipClasses: ShipClass[] = [];
	filterEnemyShipClasses: ShipClass[] = [];

	constructor(
		public gameId = '0',
		public playerNum = 0
	) {}

	get key(): string {
		return PlayerSettings.key(this.gameId, this.playerNum);
	}

	static key(gameId: string | bigint, playerNum: number) {
		return `playerSettings-${gameId}-${playerNum}`;
	}

	beforeSave() {
		this.messageTypeFilterArray = Array.from(this.messageTypeFilter);
	}

	// some state we don't want to persist on load
	afterLoad() {
		this.addWaypoint = false;
		this.messageTypeFilter = new Set<number>();
		this.messageTypeFilterArray.forEach((t) => this.messageTypeFilter.add(t));
	}

	filterMessageType(type: PlayerMessageType) {
		this.messageTypeFilter.add(Number(type));
	}
	showMessageType(type: PlayerMessageType) {
		this.messageTypeFilter.delete(Number(type));
	}
	isMessageFiltered(type: PlayerMessageType) {
		return this.messageTypeFilter.has(Number(type));
	}
	// show these players' minefields, or hide all minefields if there are none to show
	setShownMinefieldOwners(owners: MinefieldOwner[]) {
		this.showMinefieldsOf = owners;
		this.showMinefields = owners.length > 0;
	}
	// show or hide minefields. Showing them brings back the ones that were shown before
	toggleMinefields() {
		this.showMinefields = !this.showMinefields;
		if (this.showMinefields && this.showMinefieldsOf.length === 0) {
			this.showMinefieldsOf = [...MinefieldOwners];
		}
	}
	isMessageVisible(type: PlayerMessageType) {
		return !this.messageTypeFilter.has(Number(type));
	}
}
