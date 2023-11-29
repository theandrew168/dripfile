import React from "react";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router";

import { readTransfer } from "../fetch";

export default function TransferRead() {
	const { id } = useParams();
	if (!id) {
		return null;
	}

	const {
		isPending,
		isError,
		error,
		data: transfer,
	} = useQuery({
		queryKey: ["transfer", id],
		queryFn: async () => readTransfer(id),
	});

	// TODO: build a generic loading component
	if (isPending) {
		return <div>Loading...</div>;
	}

	// TODO: build a generic error component
	if (isError) {
		return <div>Error: {error.message}</div>;
	}

	return (
		<div>
			<p>ID: {transfer.id}</p>
			<p>ItineraryID: {transfer.itineraryID}</p>
			<p>Status: {transfer.status}</p>
			<p>Progress: {transfer.progress}</p>
			<p>Error: {transfer.error}</p>
			<p>CreatedAt: {transfer.createdAt.toString()}</p>
			<p>UpdatedAt: {transfer.updatedAt.toString()}</p>
		</div>
	);
}
