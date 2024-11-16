import type { CargoTransferRequest } from '$lib/types/CargoTransferRequest';
import type {
	CargoTransferTarget,
	CommandedFleet,
	Fleet,
	ShipToken,
	Waypoint,
	WaypointTransportTasks
} from '$lib/types/Fleet';
import type { CommandedPlanet, Planet } from '$lib/types/Planet';
import type { Salvage } from '$lib/types/Salvage';

export type OnOk<T> = (e: T) => void;
export type OnCancel = () => void;
export type OnShowDialog<T> = (e: T) => void;

export type CargoTransferDialogEvent = {
	src: CommandedFleet;
	dest?: CargoTransferTarget;
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

export type MergeFleetsEvent = {
	fleet: CommandedFleet;
	fleetNums: number[];
};

export type TransferCargoEvent = {
	src: CommandedFleet;
	dest?: Fleet | Planet | Salvage;
	transferAmount: CargoTransferRequest;
};

export type TransportTasksUpdateEvent = {
	fleet: CommandedFleet;
	waypoint: Waypoint;
	transportTasks: WaypointTransportTasks;
};

export type SplitAllProps = {
	onSplitAll: (e: SplitAllEvent) => Promise<void> | undefined;
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
