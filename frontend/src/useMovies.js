import axios from "axios";
import { useQuery } from "react-query";

export const KEY = "movies"

const calculateAverageScore = (movie) => {
    const totalScore = movie.ratings.reduce((acc, rate) => acc + Number(rate.value), 0);
    return totalScore / movie.ratings.length;
};

const fetchMovies = async () => {
    const { data } = await axios.get(`/movies`);
    return data;
  };

  const useMovies = () => useQuery(KEY, fetchMovies, {
    select: (data) => {
        return data.map((d) => ({
            ...d,
            avScore: calculateAverageScore(d)
        }))
    }
  })

  export default useMovies;