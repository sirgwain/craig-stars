import type { EnumDictionary } from './types/EnumDictionary';
import type { Fleet } from './types/Fleet';
import type { MapObject } from './types/MapObject';
import type { Planet } from './types/Planet';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
type EventFunction = (...args: any[]) => void;

type ProductionQueueDialogRequestedEvent = (planet?: Planet) => void;
type CargoTransferDialogRequestedEvent = (src: Fleet, target?: Fleet | Planet) => void;
type CargoTransferredEvent = (mo: MapObject) => void;

enum EventType {
	ProductionQueueDialogRequested,
	CargoTransferDialogRequested,
	CargoTransferred
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

	publishProductionQueueDialogRequestedEvent(planet?: Planet | undefined) {
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
}

export const EventManager = new Events();
