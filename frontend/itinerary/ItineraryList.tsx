import React from "react";
import { Link } from "react-router-dom";

import type { Itinerary } from "../types";

type Props = {
	itineraries: Itinerary[];
};

export default function ItineraryList({ itineraries }: Props) {
	// https://tailwindui.com/components/application-ui/lists/tables#component-4738eac883e67bf84a9f7db2446e838a
	return (
		<>
			<div className="sm:flex sm:items-center">
				<div className="sm:flex-auto">
					<h1 className="text-base font-semibold leading-6 text-gray-900">Itineraries</h1>
					<p className="mt-2 text-sm text-gray-700">A list of all your transfer itineraries.</p>
				</div>
				<div className="mt-4 sm:ml-16 sm:mt-0 sm:flex-none">
					<Link to="/itineraries/create">
						<button
							type="button"
							className="block rounded-md bg-indigo-600 px-3 py-2 text-center text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
						>
							Add itinerary
						</button>
					</Link>
				</div>
			</div>
			<div className="mt-8 flow-root">
				<div className="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
					<div className="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
						<div className="overflow-hidden shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
							<table className="min-w-full divide-y divide-gray-300">
								<thead className="bg-gray-50">
									<tr>
										<th scope="col" className="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-6">
											ID
										</th>
										<th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
											Pattern
										</th>
										<th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
											From
										</th>
										<th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
											To
										</th>
										<th scope="col" className="relative py-3.5 pl-3 pr-4 sm:pr-6">
											<span className="sr-only">Edit</span>
										</th>
									</tr>
								</thead>
								<tbody className="divide-y divide-gray-200 bg-white">
									{itineraries.map((itinerary) => (
										<tr key={itinerary.id}>
											<td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6">
												<Link to={`/itineraries/${itinerary.id}`}>{itinerary.id}</Link>
											</td>
											<td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{itinerary.pattern}</td>
											<td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
												<Link to={`/locations/${itinerary.fromLocationID}`}>{itinerary.fromLocationID}</Link>
											</td>
											<td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
												<Link to={`/locations/${itinerary.toLocationID}`}>{itinerary.toLocationID}</Link>
											</td>
											<td className="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-6">
												<a href="#" className="text-indigo-600 hover:text-indigo-900">
													Edit<span className="sr-only">, {itinerary.id}</span>
												</a>
											</td>
										</tr>
									))}
								</tbody>
							</table>
						</div>
					</div>
				</div>
			</div>
		</>
	);
}
