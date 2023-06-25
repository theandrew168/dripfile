import React from "react";
import { isRouteErrorResponse, useRouteError } from "react-router-dom";

export default function ErrorPage() {
	let statusCode = 500;
	let statusText = "Internal server error";
	let message = "Sorry, something went wrong.";

	const error = useRouteError();
	if (isRouteErrorResponse(error)) {
		statusCode = error.status;
		statusText = error.statusText;
		if (error.error?.message) {
			message = error.error.message + ".";
		}
	}

	return (
		<main className="grid min-h-full place-items-center bg-white px-6 py-24 sm:py-32 lg:px-8">
			<div className="text-center">
				<p className="text-base font-semibold text-indigo-600">{statusCode}</p>
				<h1 className="mt-4 text-3xl font-bold tracking-tight text-gray-900 sm:text-5xl">{statusText}</h1>
				<p className="mt-6 text-base leading-7 text-gray-600">{message}</p>
				<div className="mt-10 flex items-center justify-center gap-x-6">
					<a
						href="/"
						className="rounded-md bg-indigo-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
					>
						Go back home
					</a>
					<a href="#" className="text-sm font-semibold text-gray-900">
						View documentation <span aria-hidden="true">&rarr;</span>
					</a>
				</div>
			</div>
		</main>
	);
}
