export type ErrorResponse = {
	error: string;
}

export function isErrorResponse(response: any): response is ErrorResponse {
	return 'error' in response;
}
