import type { Planet } from './types/Planet';

class Events {
	public productionQueueDialogRequestedEvent?: (planet?: Planet) => void;

	publishProductionQueueDialogRequestedEvent(planet?: Planet | undefined) {
		this.productionQueueDialogRequestedEvent?.apply(planet);
	}
}

export const EventManager = new Events();
