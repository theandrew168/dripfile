import React from "react";

type Props = {
	src: string;
};

export default function Image({ src }: Props) {
	return <img src={src}></img>;
}
