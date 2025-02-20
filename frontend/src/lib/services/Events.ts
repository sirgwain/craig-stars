import type { CargoDest, CargoTransferRequest } from '$lib/types/CargoTransferRequest.svelte';
import type { Fleet, ShipToken, Waypoint, WaypointTransportTasks } from '$lib/types/cs';
import type { CommandedFleet } from '$lib/types/Fleet';
import type { CommandedPlanet } from '$lib/types/Planet';

export type OnOk<T> = (e: T) => void;
export type OnCancel = () => void;
export type OnClose = () => void;
export type OnShowDialog<T> = (e: T) => void;

export function getXFromPointerEvent(e: PointerEvent, elem: HTMLElement | undefined): number {
	if (!elem) {
		return 0;
	}
	return (e.clientX - elem.getBoundingClientRect().left) / elem.getBoundingClientRect()?.width;
}

export type CargoTransferDialogEvent = {
	src: CommandedFleet;
	dest?: CargoDest;
};

export type SplitFleetDialogEvent = {
	src: CommandedFleet;
	dest?: Fleet;
};

export type MergeFleetsDialogEvent = {
	fleet: CommandedFleet;
	otherFleetsHere: Fleet[];
};

export type ProductionQueueDialogEvent = {
	planet: CommandedPlanet;
};

export type TransportTasksDialogEvent = {
	fleet: CommandedFleet;
	waypoint: Waypoint;
	waypointIndex: number;
};

export type SplitFleetEvent = {
	src: CommandedFleet;
	dest: Fleet | undefined;
	srcTokens: ShipToken[];
	destTokens: ShipToken[];
	transferAmount: CargoTransferRequest;
};

export type SplitAllEvent = {
	fleet: CommandedFleet;
};

export type RenameFleetEvent = {
	fleet: CommandedFleet;
	name: string;
};

export type MergeFleetsEvent = {
	fleet: CommandedFleet;
	fleetNums: number[];
};

export type TransferCargoEvent = {
	src: CommandedFleet;
	dest?: CargoDest;
	transferAmount: CargoTransferRequest;
};

export type ChangeWaypointTransportTasksEvent = {
	transportTasks: WaypointTransportTasks;
} & ChangeWaypointEvent;

export type ClearProductionQueueEvent = {
	planet: CommandedPlanet;
};

export type ChangeMassDriverSpeedEvent = {
	planet: CommandedPlanet;
	warpSpeed: number;
};

export type BattlePlanChangedEvent = {
	fleet: CommandedFleet;
	battlePlanNum: number;
};

export type ChangeWaypointEvent = {
	fleet: CommandedFleet;
	waypoint: Waypoint;
	waypointIndex: number;
};

export type DeleteWaypointEvent = {
	fleet: CommandedFleet;
	waypoint: Waypoint;
};

export type SelectWaypointEvent = {
	fleet: CommandedFleet;
	waypoint: Waypoint;
};

export type NextPrevMapObjectProps = {
	onNextMapObject?: () => void;
	onPreviousMapObject?: () => void;
};

export type ClearProductionQueueProps = {
	onClearProductionQueue?: (e: ClearProductionQueueEvent) => Promise<void>;
};

export type ChangeMassDriverSpeedProps = {
	onChangeMassDriverSpeed?: (e: ChangeMassDriverSpeedEvent) => Promise<void>;
};

export type SplitAllProps = {
	onSplitAll: (e: SplitAllEvent) => Promise<void> | undefined;
};

export type BattlePlanChangedProps = {
	onBattlePlanChanged?: (e: BattlePlanChangedEvent) => Promise<void>;
};

export type RenameFleetProps = {
	onRenameFleet?: (e: RenameFleetEvent) => Promise<void>;
};

export type ChangeWaypointProps = {
	onChangeWaypoint?: (e: ChangeWaypointEvent) => Promise<void>;
};

export type DeleteWaypointProps = {
	// delete the currently selected waypoint
	// event data is optional and not currently implemented
	onDeleteWaypoint?: (e?: DeleteWaypointEvent) => Promise<void>;
};

export type SelectWaypointProps = {
	onSelectWaypoint?: (e: SelectWaypointEvent) => void;
};

// properties for a components that trigger a Dialog events

export type ShowCargoTransferDialogProps = {
	onShowCargoTransferDialog: OnShowDialog<CargoTransferDialogEvent> | undefined;
};

export type ShowSplitFleetDialogProps = {
	onShowSplitFleetDialog: OnShowDialog<SplitFleetDialogEvent> | undefined;
};

export type ShowMergeFleetsDialogProps = {
	onShowMergeFleetDialog: OnShowDialog<MergeFleetsDialogEvent> | undefined;
};

export type ShowProductionQueueDialogProps = {
	onShowProductionQueueDialog: OnShowDialog<ProductionQueueDialogEvent> | undefined;
};

export type ShowTransportTasksDialogEventProps = {
	onShowTransportTasksDialog: OnShowDialog<TransportTasksDialogEvent> | undefined;
};
