import type { ShipDesignSlot } from '$lib/types/ShipDesign';
import type { HullSlot, TechHullComponent } from '$lib/types/Tech';
import { writable } from 'svelte/store';

export type ShipDesignerContext = {
	selectedSlotIndex: number | undefined;
	selectedSlot: HullSlot | undefined;
	selectedShipDesignSlot: ShipDesignSlot | undefined;
	selectedHullComponent: TechHullComponent | undefined;
};

export const shipDesignerContext = writable<ShipDesignerContext>({
	selectedSlotIndex: undefined,
	selectedSlot: undefined,
	selectedShipDesignSlot: undefined,
	selectedHullComponent: undefined
});
