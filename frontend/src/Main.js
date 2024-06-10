import React from 'react';
import './Main.css';

import { NavLink, Outlet } from 'react-router-dom';

const MovieDashboard = () => {


    return (
        <div>
            <section className={"Background"}>
                <header>
                    <div className={'Button-layout'}>
                        <NavLink to="movies"><button>Movies</button></NavLink>
                        <NavLink to="soon"><button>Coming soon!</button></NavLink>
                        <NavLink to="trending"><button>Trending</button></NavLink>
                        <NavLink to="leaving"><button>Leaving Soon</button></NavLink>
                        <NavLink to={"health"}><button>Health</button></NavLink>
                    </div>
                    <h1 className={'Header'}>The Home Movie Depot</h1>
                </header>
                <Outlet />
                
            </section>
        </div>

    );

}

export default MovieDashboard;

