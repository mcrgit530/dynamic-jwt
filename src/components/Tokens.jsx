import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';

const Tokens = () => {
    const { username } = useParams();
    const [tokens, setTokens] = useState([]);

    useEffect(() => {
        const fetchTokens = async () => {
            const response = await fetch(`http://localhost:8080/fetchTokens/${username}`);
            if (response.ok) {
                const data = await response.text();
                setTokens(data.split("\n").filter(Boolean)); // Split tokens by new lines and filter out empty ones
            } else {
                alert('Failed to fetch tokens');
            }
        };

        fetchTokens();
    }, [username]);

    return (
        <div className="p-4">
            <h2 className="text-2xl font-bold mb-4">Tokens of {username}</h2>
            <ul className="list-disc pl-5">
                {tokens.length > 0 ? (
                    tokens.map((token, index) => (
                        <li key={index} className="mb-2 p-2 bg-gray-100 rounded shadow">
                            {token}
                        </li>
                    ))
                ) : (
                    <li>No tokens found</li>
                )}
            </ul>
        </div>
    );
};

export default Tokens;
