export type ErrorResponse = {
	error: string | Record<string, string>;
};

export function isErrorResponse(response: any): response is ErrorResponse {
	return "error" in response;
}

export type Location = {
	id: string;
	kind: string;
	createdAt: Date;
	updatedAt: Date;
};

export type LocationListResponse = {
	locations: Location[];
};

export type LocationReadResponse = {
	location: Location;
};

export type Itinerary = {
	id: string;
	fromLocationID: string;
	toLocationID: string;
	pattern: string;
	createdAt: Date;
	updatedAt: Date;
};

export type ItineraryListResponse = {
	itineraries: Itinerary[];
};

export type ItineraryReadResponse = {
	itinerary: Itinerary;
};

export type Transfer = {
	id: string;
	itineraryID: string;
	status: string;
	progress: number;
	error: string;
	createdAt: Date;
	updatedAt: Date;
};

export type TransferListResponse = {
	transfers: Transfer[];
};

export type TransferReadResponse = {
	transfer: Transfer;
};
