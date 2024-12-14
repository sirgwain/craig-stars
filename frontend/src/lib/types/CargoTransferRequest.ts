import { type Cargo, negativeCargo } from './Cargo';

export type CargoTransferRequest = {
	ironium: number;
	boranium: number;
	germanium: number;
	colonists: number;
	fuel: number;
};

// create a new CargoTransferRequest from a Cargo and Fuel
export function newCargoTransferRequest(
	cargo?: Cargo | CargoTransferRequest,
	fuel?: number
): CargoTransferRequest {
	const req = Object.assign(
		{
			ironium: 0,
			boranium: 0,
			germanium: 0,
			colonists: 0,
			fuel: 0
		},
		cargo
	);
	req.fuel = fuel ?? 0;

	return req;
}

// return a new CargoTransferRequest from this one, but negated
export function negative(req: CargoTransferRequest): CargoTransferRequest {
	return newCargoTransferRequest(negativeCargo(req), -req.fuel);
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
export function add(req: CargoTransferRequest, c: CargoTransferRequest) {
	return newCargoTransferRequest(
		{
			ironium: req.ironium + c.ironium,
			boranium: req.boranium + c.boranium,
			germanium: req.germanium + c.germanium,
			colonists: req.colonists + c.colonists
		},
		req.fuel + c.fuel
	);
}
