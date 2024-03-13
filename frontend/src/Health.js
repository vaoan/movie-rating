import React, {useEffect, useState} from 'react';
import axios from "axios";

const Health = () => {
    const [healthy, setHealth] = useState([]);

    useEffect(() => {
        const fetchHealth = async () => {
            try {
                // When running in codespaces: Terminal -> Ports tab -> Copy forwarded address under port 8080 +/api/health 
                // Important you right click 8080 -> Port visibility -> Public
                const {data} = await axios.get('https://ideal-goggles-jjjjv55grqw254j-8080.app.github.dev/api/health');
                setHealth(data);
                return data;
            } catch (error) {
                return error;
            }
        };

        fetchHealth().then(result => console.log(result));
    }, []);

    return (
        <div>
            <h1>{JSON.stringify(healthy)}</h1>
        </div>
    )
};

export default Health;