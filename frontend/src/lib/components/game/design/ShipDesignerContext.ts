import type { ShipDesignSlot } from '$lib/types/ShipDesign';
import type { HullSlot } from '$lib/types/Tech';
import { writable } from 'svelte/store';

export type ShipDesignerContext = {
	selectedSlotIndex: number | undefined;
	selectedSlot: HullSlot | undefined;
	selectedShipDesignSlot: ShipDesignSlot | undefined;
};

export const shipDesignerContext = writable<ShipDesignerContext>({
	selectedSlotIndex: undefined,
	selectedSlot: undefined,
	selectedShipDesignSlot: undefined
});
