import type { MessageType } from './Message';

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

export class PlayerSettings {
	planetViewState = PlanetViewState.Normal;
	addWaypoint = false;
	setPacketDest = false;
	showPlanetNames = false;
	showFleetTokenCounts = false;
	showScanners = true;
	showAllyScanners = true;
	showMineFields = true;
	showIdleFleetsOnly = false;
	showMessagePane = false;
	scannerPercent = 100;
	mineralScale = 5000;
	messageTypeFilter = new Set<number>();
	messageTypeFilterArray: number[] = [];
	sortPlanetsKey = 'name';
	sortPlanetsDescending = false;
	showAllPlanets = false;
	sortFleetsKey = 'name';
	sortFleetsDescending = false;
	sortBattlesKey = 'num';
	sortBattlesDescending = false;

	constructor(public gameId = 0, public playerNum = 0) {}

	get key(): string {
		return PlayerSettings.key(this.gameId, this.playerNum);
	}

	static key(gameId: number, playerNum: number) {
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

	filterMessageType(type: MessageType) {
		this.messageTypeFilter.add(Number(type));
	}
	showMessageType(type: MessageType) {
		this.messageTypeFilter.delete(Number(type));
	}
	isMessageFiltered(type: MessageType) {
		return this.messageTypeFilter.has(Number(type));
	}
	isMessageVisible(type: MessageType) {
		return !this.messageTypeFilter.has(Number(type));
	}
}
