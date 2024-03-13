
import React from 'react';
import { BrowserRouter as Router, Routes, Route}
    from 'react-router-dom';
import Main from "./Main";
import Health from "./Health";

function App() {
    return (
        <Router>
            <Routes>
                <Route path='/' element={<Main />} />
                <Route path='/health' element={<Health/>} />
            </Routes>
        </Router>
    );
}

export default App;