export type EnumDictionary<T extends string | symbol | number, U> = {
	[K in T]: U;
};

export type PartialEnumDictionary<T extends string | symbol | number, U> = {
	[K in T]?: U;
};
