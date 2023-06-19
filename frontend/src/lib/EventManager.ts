import type { EnumDictionary } from './types/EnumDictionary';
import type { CommandedFleet, Fleet } from './types/Fleet';
import type { MapObject } from './types/MapObject';
import type { CommandedPlanet, Planet } from './types/Planet';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
type EventFunction = (...args: any[]) => void;

type ProductionQueueDialogRequestedEvent = (planet: CommandedPlanet) => void;
type CargoTransferDialogRequestedEvent = (src: CommandedFleet, target?: Fleet | Planet) => void;
type CargoTransferredEvent = (mo: MapObject) => void;
type SplitFleetDialogRequestedEvent = (src: CommandedFleet) => void;
type MergeFleetDialogRequestedEvent = (src: CommandedFleet) => void;

enum EventType {
	ProductionQueueDialogRequested,
	CargoTransferDialogRequested,
	CargoTransferred,
	SplitFleetDialogRequested,
	MergeFleetDialogRequested
}

class Events {
	private events: EnumDictionary<EventType, EventFunction[]> = {} as EnumDictionary<
		EventType,
		EventFunction[]
	>;

	/**
	 * Subscribe to an event, return a callback to unsubscribe
	 * @param eventType
	 * @param event
	 * @returns
	 */
	private subscribe(eventType: EventType, event: EventFunction) {
		this.events[eventType] = this.events[eventType] ?? [];
		this.events[eventType].push(event);

		// return an unsubscribe method
		return () => {
			this.events[eventType] = this.events[eventType].filter((e) => e != event);
		};
	}

	// ProductionQueueDialog
	public subscribeProductionQueueDialogRequestedEvent(
		event: ProductionQueueDialogRequestedEvent
	): () => void {
		return this.subscribe(EventType.ProductionQueueDialogRequested, event);
	}

	publishProductionQueueDialogRequestedEvent(planet: CommandedPlanet) {
		this.events[EventType.ProductionQueueDialogRequested].forEach((e) => e.apply(e, [planet]));
	}

	// CargoTransferDialog
	public subscribeCargoTransferDialogRequestedEvent(
		event: CargoTransferDialogRequestedEvent
	): () => void {
		return this.subscribe(EventType.CargoTransferDialogRequested, event);
	}

	publishCargoTransferDialogRequestedEvent(src: Fleet, target?: Fleet | Planet) {
		this.events[EventType.CargoTransferDialogRequested].forEach((e) => e.apply(e, [src, target]));
	}

	// CargoTransferred
	public subscribeCargoTransferredEvent(event: CargoTransferredEvent): () => void {
		return this.subscribe(EventType.CargoTransferred, event);
	}

	publishCargoTransferredEvent(mo: MapObject) {
		this.events[EventType.CargoTransferred].forEach((e) => e.apply(e, [mo]));
	}

	publishSplitFleetDialogRequestedEvent(src: CommandedFleet) {
		this.events[EventType.SplitFleetDialogRequested].forEach((e) => e.apply(e, [src]));
	}

	public subscribeSplitFleetDialogRequestedEvent(event: SplitFleetDialogRequestedEvent): () => void {
		return this.subscribe(EventType.SplitFleetDialogRequested, event);
	}

	publishMergeFleetDialogRequestedEvent(src: CommandedFleet) {
		this.events[EventType.MergeFleetDialogRequested].forEach((e) => e.apply(e, [src]));
	}

	public subscribeMergeFleetDialogRequestedEvent(event: MergeFleetDialogRequestedEvent): () => void {
		return this.subscribe(EventType.MergeFleetDialogRequested, event);
	}

}

export const EventManager = new Events();
