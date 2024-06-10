// Function to calculate the average score of a movie

import React from "react";
import useMovies from "./useMovies";
import Movie from "./MovieCard";

const highestRaiting = (movies) => {
    let movie = {average_rating: 0};
    movies?.forEach(element => {
        if(Number(element?.average_rating) >= Number(movie?.average_rating)){
            movie = {...element}
        }
    });
    return movie;
}

const Trending = () => {

    const {data, isLoading, isError} = useMovies();
    if (isLoading) return <>Loading...</>
    if (isError) return <>Error!!</>
    const movie = highestRaiting(data);
    return <><h1>TRENDING!!</h1><Movie movie={movie} /></>
}

export default Trending;