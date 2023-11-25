import React from "react";
import { useQuery } from "@tanstack/react-query";

import type { ItineraryListResponse } from "../types";
import ItineraryEmpty from "./ItineraryEmpty";
import ItineraryList from "./ItineraryList";

export default function ItineraryPage() {
	const { isPending, isError, error, data } = useQuery({
		queryKey: ["itinerary"],
		queryFn: async () => {
			const response = await fetch("/api/v1/itinerary");
			if (!response.ok) {
				throw new Error("Network response was not OK");
			}

			const data: ItineraryListResponse = await response.json();
			return data;
		},
	});

	// TODO: build a generic loading component
	if (isPending) {
		return <div>Loading...</div>;
	}

	// TODO: build a generic error component
	if (isError) {
		return <div>Error: {error.message}</div>;
	}

	const itineraries = data.itineraries;
	const hasItineraries = itineraries.length > 0;
	return hasItineraries ? <ItineraryList itineraries={itineraries} /> : <ItineraryEmpty />;
}
