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
	showMineFields = true;
	showIdleFleetsOnly = false;
	scannerPercent = 100;
	mineralScale = 5000;
	messageTypeFilter = new Set<string>();
}
