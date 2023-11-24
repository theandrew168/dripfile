import React from "react";
import { useQuery } from "@tanstack/react-query";

import type { ItineraryListResponse } from "../types";
import ItineraryEmpty from "./ItineraryEmpty";
import ItineraryList from "./ItineraryList";

export default function ItineraryPage() {
	const { isPending, isError, error, data } = useQuery({
		queryKey: ["itineraries"],
		queryFn: async () => {
			const response = await fetch("/api/v1/itineraries");
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
	const hasLocations = itineraries.length > 0;
	return hasLocations ? <ItineraryList itineraries={itineraries} /> : <ItineraryEmpty />;
}