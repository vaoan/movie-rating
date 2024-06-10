import React from "react";
import { NavLink } from "react-router-dom";
import Movie from "./MovieCard";
import useMovies from "./useMovies";

const Movies = () => {

    const {data:movies, isLoading, isError} = useMovies();

    if(isLoading) return <div>Loading...</div>;
    if(isError) return <>Error...</>

    return <>
        <div>{movies?.map((movie) => <React.Fragment key={JSON.stringify(movie)}>
            <Movie movie={movie} />
        </React.Fragment>)}</div>
    </>
}

export default Movies;