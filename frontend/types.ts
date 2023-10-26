export type ErrorResponse = {
	error: string | Record<string, string>;
};

export function isErrorResponse(response: any): response is ErrorResponse {
	return "error" in response;
}

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

export type Itinerary = {
	id: string;
	pattern: string;
	fromLocationID: string;
	toLocationID: string;
};

export type ItineraryListResponse = {
	itineraries: Itinerary[];
};

export type ItineraryReadResponse = {
	itinerary: Itinerary;
};
