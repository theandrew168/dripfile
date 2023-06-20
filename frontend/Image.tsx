import React from "react";

type Props = {
	src: string;
};

export function Image({ src }: Props) {
	return <img src={src}></img>;
}
