//{"title":"Life of Brian","ratings":[{"movie_ratings_id":1,"source":"Internet Movie Database","value":80},{"movie_ratings_id":1,"source":"Rotten Tomatoes","value":95},
import React from 'react';
import { NavLink } from 'react-router-dom';

const MovieCard = ({movie}) => {
    const {title, ratings} = movie;
    return <>
            <NavLink to={"/movie/" + title}><h3>{title}</h3></NavLink>
            <>{
                ratings?.map((rate) => <React.Fragment key={JSON.stringify(movie) + JSON.stringify(rate)}>
                    
                    <div>{rate?.source}</div>
                    <div>{rate?.value}</div>
                </React.Fragment>)
            }</>
        </>
}

export default MovieCard