import React from 'react';
import './Main.css';
import reel from './Moviereel.gif'

const MovieDashboard = () => {


    return (
        <div>
            <section className={"Background"}>
                <header>
                    <div className={'Button-layout'}>
                        <button>Movies</button>
                        <button>Coming soon!</button>
                        <button>Trending</button>
                        <button>Leaving Soon</button>
                        <button>Health</button>
                    </div>
                    <h1 className={'Header'}>The Home Movie Depot</h1>

                </header>
                <img src={reel} className={'Movie-gif'} alt="logo"/>
            </section>
        </div>

    );

}

export default MovieDashboard;

