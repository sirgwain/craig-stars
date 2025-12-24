import type { Fleet, MineralPacket, Planet, Salvage } from '$lib/types/cs-proto';
import { CargoSchema, type Cargo, type CargoJson } from '$lib/types/cs-proto';
import { create } from '@bufbuild/protobuf';
import { negativeCargo, totalCargo } from './Cargo';
import { clamp } from '$lib/services/Math';

// a destination that cargo can be transferred to/from
export type CargoDest = Fleet | Planet | MineralPacket | Salvage | undefined;

export type CargoTransferRequest = {
	ironium: number;
	boranium: number;
	germanium: number;
	colonists: number;
	fuel: number;
};

export function newCargoTransferRequest(
	cargo: CargoJson | undefined,
	fuel: number | undefined
): CargoTransferRequest {
	return {
		ironium: cargo?.ironium ?? 0,
		boranium: cargo?.boranium ?? 0,
		germanium: cargo?.germanium ?? 0,
		colonists: cargo?.colonists ?? 0,
		fuel: fuel ?? 0
	};
}

export function emptyCargoTransferRequest(): CargoTransferRequest {
	return {
		ironium: 0,
		boranium: 0,
		germanium: 0,
		colonists: 0,
		fuel: 0
	};
}

export function jsonData(r: CargoTransferRequest): CargoJson & { fuel: number } {
	return {
		ironium: r.ironium,
		boranium: r.boranium,
		germanium: r.germanium,
		colonists: r.colonists,
		fuel: r.fuel
	};
}

export function getCargo(r: CargoTransferRequest): Cargo {
	return create(CargoSchema, {
		ironium: r.ironium,
		boranium: r.boranium,
		germanium: r.germanium,
		colonists: r.colonists
	});
}

export function setCargo(r: CargoTransferRequest, c: CargoJson) {
	r.ironium = c.ironium ?? 0;
	r.boranium = c.boranium ?? 0;
	r.germanium = c.germanium ?? 0;
	r.colonists = c.colonists ?? 0;
}

// return a new CargoTransferRequest from this one, but negated
export function negative(req: CargoTransferRequest): CargoTransferRequest {
	return { ...negativeCargo(getCargo(req)), fuel: -req.fuel };
}

// return the absolute size of this transfer request
export function absoluteSize(req: CargoTransferRequest): number {
	return (
		Math.abs(req.ironium) +
		Math.abs(req.boranium) +
		Math.abs(req.germanium) +
		Math.abs(req.colonists) +
		Math.abs(req.fuel)
	);
}

// return the absolute size of this transfer request
export function absoluteCargoSize(req: CargoTransferRequest): number {
	return (
		Math.abs(req.ironium) +
		Math.abs(req.boranium) +
		Math.abs(req.germanium) +
		Math.abs(req.colonists)
	);
}

// add a CargoTransferRequest to another CargoTransferRequest, returning a new one
export function add(req: CargoTransferRequest, c: CargoTransferRequest): CargoTransferRequest {
	return {
		ironium: req.ironium + c.ironium,
		boranium: req.boranium + c.boranium,
		germanium: req.germanium + c.germanium,
		colonists: req.colonists + c.colonists,
		fuel: req.fuel + c.fuel
	};
}

export function suggestFuelTransfer(
	srcFuel: number,
	srcFuelCapacity: number,
	destFuel: number,
	destFuelCapacity: number
): number {
	const totalFuel = srcFuel + destFuel;
	const totalCap = srcFuelCapacity + destFuelCapacity;
	if (totalCap <= 0) return 0;

	const desiredDestFuel = Math.floor((totalFuel * destFuelCapacity) / totalCap);
	const currentDestFuel = destFuel;

	// delta = desired - current
	const suggested = currentDestFuel - desiredDestFuel;

	return clampFuelTransfer(suggested, srcFuel, srcFuelCapacity, destFuel, destFuelCapacity);
}

export function clampFuelTransfer(
	delta: number,
	srcFuel: number,
	srcCapacity: number,
	destFuel: number,
	destCapacity: number
) {
	// if src has 150/200mg of fuel, it can give -150 or take 50
	const minDeltaFromSrc = -srcFuel; // delta <= this (can't send more than you have)
	const maxDeltaFromSrc = srcCapacity - srcFuel; // delta >= this

	// if dest has 150/200mg of fuel, it can take -50, or give 150
	const minDeltaFromDest = -(destCapacity - destFuel);
	const maxDeltaFromDest = destFuel;

	const minDelta = Math.max(minDeltaFromSrc, minDeltaFromDest);
	const maxDelta = Math.min(maxDeltaFromSrc, maxDeltaFromDest);

	return clamp(delta, minDelta, maxDelta);
}

export function suggestCargoTransfer(
	src: CargoJson,
	srcCapacity: number,
	dest: CargoJson,
	destCapacity: number
): CargoJson {
	const totalCap = srcCapacity + destCapacity;
	if (totalCap <= 0) return { ironium: 0, boranium: 0, germanium: 0, colonists: 0 };

	const lane = (key: keyof CargoJson) => {
		const s = src[key] ?? 0;
		const d = dest[key] ?? 0;
		const total = s + d;

		// desired dest lane amount proportional to dest capacity
		const desiredDest = Math.floor((total * destCapacity) / totalCap);

		// delta is "change to source" (same as fuel): currentDest - desiredDest
		return d - desiredDest;
	};

	const suggested: CargoJson = {
		ironium: lane('ironium'),
		boranium: lane('boranium'),
		germanium: lane('germanium'),
		colonists: lane('colonists')
	};

	return clampCargoTransfer(suggested, src, srcCapacity, dest, destCapacity);
}

/**
 * Clamp a 4-lane cargo transfer with a shared capacity.
 *
 * delta is change-to-source (same sign semantics as your fuel):
 *   negative => src gives to dest
 *   positive => dest gives to src
 */
export function clampCargoTransfer(
	delta: CargoJson,
	src: CargoJson,
	srcCapacity: number,
	dest: CargoJson,
	destCapacity: number
): CargoJson {
	const clampComponentDelta = (
		delta: number | undefined,
		srcAmt: number | undefined,
		destAmt: number | undefined
	) =>
		// src' = src + delta must be >=0  => delta >= -srcAmt
		// dest' = dest - delta must be >=0 => delta <= destAmt
		clamp(delta ?? 0, -(srcAmt ?? 0), destAmt ?? 0);

	// 1) per-component clamp so no component goes negative on either side
	const out: CargoJson = {
		ironium: clampComponentDelta(delta.ironium, src.ironium, dest.ironium),
		boranium: clampComponentDelta(delta.boranium, src.boranium, dest.boranium),
		germanium: clampComponentDelta(delta.germanium, src.germanium, dest.germanium),
		colonists: clampComponentDelta(delta.colonists, src.colonists, dest.colonists)
	};

	// 2) shared-capacity clamp on TOTAL
	const srcTotal = totalCargo(src);
	const destTotal = totalCargo(dest);
	const sum = totalCargo(out); // total change to src

	// total constraints:
	// srcTotal' = srcTotal + sum <= srcCapacity  => sum <= srcCapacity - srcTotal
	// destTotal' = destTotal - sum <= destCapacity => -sum <= destCapacity - destTotal => sum >= destTotal - destCapacity
	const minSum = destTotal - destCapacity; // usually <= 0
	const maxSum = srcCapacity - srcTotal; // usually >= 0

	// already feasible
	if (sum >= minSum && sum <= maxSum) return out;

	// We only ever need to pull deltas back toward 0.
	// If sum is too positive: reduce positive components (dest->src) toward 0
	// If sum is too negative: reduce negative components (src->dest) toward 0
	const want = clamp(sum, minSum, maxSum);
	let need = Math.abs(sum - want);
	const reducePositive = sum > want;

	// Reduce biggest movers first (less surprising than fixed order)
	const keys: (keyof CargoJson)[] = ['ironium', 'boranium', 'germanium', 'colonists'];
	keys.sort((a, b) => Math.abs(out[b] ?? 0) - Math.abs(out[a] ?? 0));

	for (const k of keys) {
		if (need <= 0) break;
		const v = out[k] ?? 0;

		if (reducePositive) {
			if (v <= 0) continue;
			const take = Math.min(v, need);
			out[k] = v - take; // toward 0
			need -= take;
		} else {
			if (v >= 0) continue;
			const take = Math.min(-v, need);
			out[k] = v + take; // toward 0 (less negative)
			need -= take;
		}
	}

	// (At this point, given sane capacities (totals <= caps), need should reach 0.)
	return out;
}
