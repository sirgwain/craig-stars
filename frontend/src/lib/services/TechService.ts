import techjson from '$lib/ssr/techs.json';
import {
	GetTechsResponseSchema,
	type GetTechsResponseJson,
	type TechDefense,
	type TechHull,
	type TechHullComponent,
	type TechPlanetary,
	type TechPlanetaryScanner,
	type TechTerraform
} from '$lib/types/cs-proto';
import type { TechLike, TechStore } from '$lib/types/Tech';
import { fromJson, type UnknownField } from '@bufbuild/protobuf';
import { kebabCase } from 'lodash-es';

export class TechService implements TechStore {
	$typeName: 'craig_stars.v1.GetTechsResponse';
	$unknown?: UnknownField[] | undefined;

	techs: TechLike[] = [];
	planetaryScanners: TechPlanetaryScanner[] = [];
	terraforms: TechTerraform[] = [];
	defenses: TechDefense[] = [];
	planetaries: TechPlanetary[] = [];
	hullComponents: TechHullComponent[] = [];
	hulls: TechHull[] = [];

	techsByName: Map<string, TechLike> = new Map();
	hullsByName: Map<string, TechHull> = new Map();
	hullComponentsByName: Map<string, TechHullComponent> = new Map();

	constructor(store?: TechStore) {
		this.$typeName = 'craig_stars.v1.GetTechsResponse';
		store = store ?? fromJson(GetTechsResponseSchema, techjson as unknown as GetTechsResponseJson);
		this.buildMaps(store);
	}

	private buildMaps(store: TechStore) {
		this.planetaryScanners = store.planetaryScanners ?? [];
		this.terraforms = store.terraforms ?? [];
		this.defenses = store.defenses ?? [];
		this.planetaries = store.planetaries ?? [];
		this.hullComponents = store.hullComponents ?? [];
		this.hulls = store.hulls ?? [];

		this.techs = [];
		this.techs = this.techs.concat(this.hullComponents);
		this.techs = this.techs.concat(this.planetaryScanners);
		this.techs = this.techs.concat(this.defenses);
		this.techs = this.techs.concat(this.planetaries);
		this.techs = this.techs.concat(this.hulls);
		this.techs = this.techs.concat(this.terraforms);

		this.techsByName = new Map(this.techs.map((t) => [kebabCase(t.tech?.name), t]));
		this.hullComponentsByName = new Map(
			this.hullComponents.map((t) => [kebabCase(t.tech?.name), t])
		);
		this.hullsByName = new Map(this.hulls.map((t) => [kebabCase(t.tech?.name), t]));
	}

	getTech(name: string): TechLike | undefined {
		return this.techsByName.get(kebabCase(name));
	}

	getHull(name: string): TechHull | undefined {
		return this.hullsByName.get(kebabCase(name));
	}

	getHullComponent(name: string): TechHullComponent | undefined {
		return this.hullComponentsByName.get(kebabCase(name));
	}
}
