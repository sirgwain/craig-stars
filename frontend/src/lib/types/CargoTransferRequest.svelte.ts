import { negativeCargo } from './Cargo';
import { type Cargo } from './cs';

export class CargoTransferRequest {
	ironium = $state(0);
	boranium = $state(0);
	germanium = $state(0);
	colonists = $state(0);
	fuel = $state(0);

	constructor(cargo?: Cargo, fuel?: number) {
		this.ironium = cargo?.ironium ?? 0;
		this.boranium = cargo?.boranium ?? 0;
		this.germanium = cargo?.germanium ?? 0;
		this.colonists = cargo?.colonists ?? 0;
		this.fuel = fuel ?? 0;
	}

	public jsonData(): Cargo & { fuel: number } {
		return {
			ironium: this.ironium,
			boranium: this.boranium,
			germanium: this.germanium,
			colonists: this.colonists,
			fuel: this.fuel
		};
	}
}

// return a new CargoTransferRequest from this one, but negated
export function negative(req: CargoTransferRequest): CargoTransferRequest {
	return new CargoTransferRequest(negativeCargo(req), -req.fuel);
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
	return new CargoTransferRequest(
		{
			ironium: req.ironium + c.ironium,
			boranium: req.boranium + c.boranium,
			germanium: req.germanium + c.germanium,
			colonists: req.colonists + c.colonists
		},
		req.fuel + c.fuel
	);
}
