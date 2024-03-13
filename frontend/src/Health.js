import React, {useEffect, useState} from 'react';
import axios from "axios";

const Health = () => {
    const [healthy, setHealth] = useState([]);

    useEffect(() => {
        const fetchHealth = async () => {
            try {
                const {data} = await axios.get('http://localhost:8080/api/health');
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