import type { ShipDesignSlotJson as ShipDesignSlot } from '$lib/types/cs-proto';
import type { TechHullComponent, TechHullSlot } from '$lib/types/cs-proto';
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
