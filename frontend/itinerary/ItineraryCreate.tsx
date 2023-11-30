import React, { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import type { Location, NewItinerary } from "../types";
import Alert from "../Alert";
import { createItinerary, listLocations } from "../fetch";

export default function ItineraryCreate() {
	const [locations, setLocations] = useState<Location[]>([]);

	const [fromLocationID, setFromLocationID] = useState("");
	const [toLocationID, setToLocationID] = useState("");
	const [pattern, setPattern] = useState("");

	const navigate = useNavigate();
	const queryClient = useQueryClient();

	const locationsQuery = useQuery({
		queryKey: ["location"],
		queryFn: async () => listLocations(),
	});
	useEffect(() => {
		if (locationsQuery.isPending) return;
		if (locationsQuery.isError) return;
		setLocations(locationsQuery.data);
	}, [locationsQuery]);

	const { mutate, isPending, isError, error } = useMutation({
		mutationFn: (itinerary: NewItinerary) => createItinerary(itinerary),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["itinerary"] });
			navigate("/itinerary");
		},
	});

	// https://tailwindui.com/components/application-ui/forms/form-layouts#component-dcf2bee8aa4fbef0d4623df5b9718da8
	return (
		<>
			{isError && <Alert type="failure" message={error.message} />}
			<div className="space-y-10 divide-y divide-gray-900/10">
				<div className="grid grid-cols-1 gap-x-8 gap-y-8 md:grid-cols-3">
					<div className="px-4 sm:px-0">
						<h2 className="text-base font-semibold leading-7 text-gray-900">Itinerary</h2>
						<p className="mt-1 text-sm leading-6 text-gray-600">A plan for transferring files between locations.</p>
					</div>

					<form
						onSubmit={(event) => {
							event.preventDefault();
							event.stopPropagation();
							mutate({ fromLocationID, toLocationID, pattern });
						}}
						className="bg-white shadow-sm ring-1 ring-gray-900/5 sm:rounded-xl md:col-span-2"
					>
						<div className="px-4 py-6 sm:p-8">
							<div className="grid max-w-2xl grid-cols-1 gap-x-6 gap-y-8 sm:grid-cols-6">
								<div className="sm:col-span-4">
									<label htmlFor="fromLocationID" className="block text-sm font-medium leading-6 text-gray-900">
										From
									</label>
									<div className="mt-2">
										<select
											id="fromLocationID"
											name="fromLocationID"
											value={fromLocationID}
											onChange={(event) => setFromLocationID(event.target.value)}
											className="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:max-w-xs sm:text-sm sm:leading-6"
										>
											{locations.map((location) => (
												<option key={location.id} value={location.id}>
													{location.id}
												</option>
											))}
										</select>
									</div>
								</div>
								<div className="sm:col-span-4">
									<label htmlFor="toLocationID" className="block text-sm font-medium leading-6 text-gray-900">
										To
									</label>
									<div className="mt-2">
										<select
											id="toLocationID"
											name="toLocationID"
											defaultValue={toLocationID}
											onChange={(event) => setToLocationID(event.target.value)}
											className="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:max-w-xs sm:text-sm sm:leading-6"
										>
											{locations.map((location) => (
												<option key={location.id} value={location.id}>
													{location.id}
												</option>
											))}
										</select>
									</div>
								</div>
								<div className="sm:col-span-4">
									<label htmlFor="pattern" className="block text-sm font-medium leading-6 text-gray-900">
										Pattern
									</label>
									<div className="mt-2">
										<div className="flex rounded-md shadow-sm ring-1 ring-inset ring-gray-300 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-600 sm:max-w-md">
											<input
												type="text"
												id="pattern"
												name="pattern"
												value={pattern}
												onChange={(event) => setPattern(event.target.value)}
												className="block flex-1 border-0 bg-transparent py-1.5 pl-2 text-gray-900 placeholder:text-gray-400 focus:ring-0 sm:text-sm sm:leading-6"
											/>
										</div>
									</div>
								</div>
							</div>
						</div>
						<div className="flex items-center justify-end gap-x-6 border-t border-gray-900/10 px-4 py-4 sm:px-8">
							<Link to="/itinerary">
								<button type="button" className="text-sm font-semibold leading-6 text-gray-900">
									Cancel
								</button>
							</Link>
							<button
								type="submit"
								disabled={isPending}
								className="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
							>
								{isPending ? "Creating..." : "Create"}
							</button>
						</div>
					</form>
				</div>
			</div>
		</>
	);
}
