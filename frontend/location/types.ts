export type Location = {
	id: string;
	kind: string;
};

export type LocationListResponse = {
	locations: Location[];
};

export type LocationReadResponse = {
	location: Location;
};
