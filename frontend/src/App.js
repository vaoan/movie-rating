
import React from 'react';
import { BrowserRouter as Router, Routes, Route}
    from 'react-router-dom';
import Main from "./Main";
import Health from "./Health";
import reel from './Moviereel.gif'
import Movies from "./Movies"
import './Main.css';
import axios from 'axios';
import { QueryClient, QueryClientProvider } from 'react-query';
import Leaving from './Leaving';
import Trending from './Trending';
import Soon from './Soon';
import Movie from './Movie';

axios.defaults.baseURL = "https://vigilant-disco-wg4gg7vgp5c9q7-8080.app.github.dev/api"
const queryClient = new QueryClient();

function App() {
    return (
        <>
        <QueryClientProvider client={queryClient}>
        <Router>
            <Routes>
                <Route path='/' element={<Main />}>
                    <Route index element={<div><img src={reel} className={'Movie-gif'} alt="logo"/></div>} />
                    <Route path='health' element={<Health/>} />
                    <Route path='movies' element={<Movies />} />
                    <Route path="leaving" element={<Leaving />} />
                    <Route path="trending" element={<Trending />} />
                    <Route path='soon' element={<Soon />} />
                    <Route path='movie/:name' element={<Movie />} />
                </Route>
            </Routes>
        </Router>
        </QueryClientProvider>
        </>
    );
}

export default App;