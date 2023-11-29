import React from "react";
import { useQuery } from "@tanstack/react-query";

import ItineraryEmpty from "./ItineraryEmpty";
import ItineraryList from "./ItineraryList";
import { listItineraries } from "../fetch";

export default function ItineraryPage() {
	const {
		isPending,
		isError,
		error,
		data: itineraries,
	} = useQuery({
		queryKey: ["itinerary"],
		queryFn: async () => listItineraries(),
	});

	// TODO: build a generic loading component
	if (isPending) {
		return <div>Loading...</div>;
	}

	// TODO: build a generic error component
	if (isError) {
		return <div>Error: {error.message}</div>;
	}

	const hasItineraries = itineraries.length > 0;
	return hasItineraries ? <ItineraryList itineraries={itineraries} /> : <ItineraryEmpty />;
}
