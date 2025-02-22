import type { ShipDesignSlot, TechHullComponent, TechHullSlot } from '$lib/types/cs';
import { writable } from 'svelte/store';

export type ShipDesignerContext = {
	selectedSlotIndex: number | undefined;
	selectedSlot: TechHullSlot | undefined;
	selectedShipDesignSlot: ShipDesignSlot | undefined;
	selectedHullComponent: TechHullComponent | undefined;
};

export const shipDesignerContext = writable<ShipDesignerContext>({
	selectedSlotIndex: undefined,
	selectedSlot: undefined,
	selectedShipDesignSlot: undefined,
	selectedHullComponent: undefined
});
