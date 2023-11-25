import React from "react";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router";

import { type TransferReadResponse } from "../types";

export default function TransferRead() {
	const { id } = useParams();

	const { isPending, isError, error, data } = useQuery({
		queryKey: ["transfer", id],
		queryFn: async () => {
			const response = await fetch(`/api/v1/transfer/${id}`);
			if (!response.ok) {
				throw new Error("Network response was not OK");
			}

			const data: TransferReadResponse = await response.json();
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

	const transfer = data.transfer;
	return (
		<div>
			<p>ID: {transfer.id}</p>
			<p>ItineraryID: {transfer.itineraryID}</p>
			<p>Status: {transfer.status}</p>
			<p>Progress: {transfer.progress}</p>
			<p>Error: {transfer.error}</p>
		</div>
	);
}
