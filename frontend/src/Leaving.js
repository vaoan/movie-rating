import Movie from "./MovieCard";
import useMovies from "./useMovies"

const lowestRaiting = (movies) => {
    let movie = {average_rating: 100};
    movies?.forEach(element => {
        if(Number(element?.average_rating) <= Number(movie?.average_rating)){
            movie = {...element}
        }
    });
    return movie;
}

const Leaving = () => {
    const {data, isLoading, isError} = useMovies();
    if (isLoading) return <>Loading...</>
    if (isError) return <>Error!!</>
    const movie = lowestRaiting(data);
    return <Movie movie={movie} />

}

export default Leaving;