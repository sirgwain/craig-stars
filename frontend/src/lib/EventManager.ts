import type { Fleet } from './types/Fleet';
import type { Planet } from './types/Planet';

class Events {
	public productionQueueDialogRequestedEvent?: (planet?: Planet) => void;
	public cargoTransferDialogRequestedEvent?: (src: Fleet, target?: Fleet | Planet) => void;

	publishProductionQueueDialogRequestedEvent(planet?: Planet | undefined) {
		this.productionQueueDialogRequestedEvent?.apply(planet);
	}

	publishCargoTransferDialogRequestedEvent(src: Fleet, target?: Fleet | Planet) {
		this.cargoTransferDialogRequestedEvent?.apply(this, [src, target]);
	}
}

export const EventManager = new Events();
