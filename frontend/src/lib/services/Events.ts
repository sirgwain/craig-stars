import type { Fleet, Waypoint, ShipToken, WaypointTransportTasks } from '#lib/types/cs-proto.js';
import type { CargoDest, CargoTransferRequest } from '#lib/types/CargoTransferRequest.js';
import type { CommandedFleet } from '#lib/types/Fleet.js';
import type { CommandedPlanet } from '#lib/types/Planet.js';

export type OnOk<T> = (e: T) => void;
export type OnCancel = () => void;
export type OnClose = () => void;
export type OnShowDialog<T> = (e: T) => void;

// Keyboard activation shares the existing pointer behavior and tooltip positioning.
export function onPointerKeyDown(e: KeyboardEvent) {
	if (e.key !== 'Enter' && e.key !== ' ') return;
	if (!(e.currentTarget instanceof HTMLElement)) return;
	e.preventDefault();
	e.stopPropagation();
	if (e.repeat) return;
	const rect = e.currentTarget.getBoundingClientRect();
	e.currentTarget.dispatchEvent(
		new PointerEvent('pointerdown', {
			bubbles: true,
			cancelable: true,
			clientX: rect.left + rect.width / 2,
			clientY: rect.top + rect.height / 2
		})
	);
}

export function onPointerKeyUp(e: KeyboardEvent) {
	if (e.key !== 'Enter' && e.key !== ' ') return;
	e.preventDefault();
	e.stopPropagation();
	window.dispatchEvent(new PointerEvent('pointerup'));
}

export function sliderValueForKey(key: string, value: number, min: number, max: number) {
	let next: number;
	switch (key) {
		case 'ArrowLeft':
		case 'ArrowDown':
			next = value - 1;
			break;
		case 'ArrowRight':
		case 'ArrowUp':
			next = value + 1;
			break;
		case 'PageDown':
			next = value - 10;
			break;
		case 'PageUp':
			next = value + 10;
			break;
		case 'Home':
			next = min;
			break;
		case 'End':
			next = max;
			break;
		default:
			return undefined;
	}
	return Math.min(max, Math.max(min, next));
}

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
	onShowCargoTransferDialog?: OnShowDialog<CargoTransferDialogEvent> | undefined;
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
