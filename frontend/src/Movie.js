import React from "react";
import { useParams } from "react-router-dom";

import MovieCard from "./MovieCard"
import useMovies from "./useMovies"

const Movie = () => {
    const { name } = useParams();
    const {data: movies, isLoading} = useMovies();

    if(isLoading) return <>Loading...</>;

    const movie = movies?.find((movie) => movie.title == name);


    return (<><MovieCard movie={movie} /></>);
}

export default Movie;