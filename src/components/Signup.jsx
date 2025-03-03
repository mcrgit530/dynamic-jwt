import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { jwtDecode } from 'jwt-decode';

const Signup = () => {
    const [username, setUsername] = useState('');
        const [password, setPassword] = useState('');
        const navigate = useNavigate();
    
        const handleSubmit = async (e) => {
            e.preventDefault();
            const response = await fetch('http://localhost:8080/signup', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username, password }),
            });
    
            if (response.ok) {
                const token = await response.text(); // Getting the token from the response
                localStorage.setItem('token', token);  // Storing the token in localStorage
    
                // Decoding the token to get the username
                const decoded = jwtDecode(token);
                const decodedUsername = decoded.username;
    
                // Navigate to the /tokens route with the username
                navigate(`/tokens/${decodedUsername}`);
            } else {
                alert('Signup failed');
            }
        };

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100">
            <h2 className="text-2xl font-bold mb-4">Signup</h2>
            <form onSubmit={handleSubmit} className="bg-white p-6 rounded shadow-md w-full max-w-sm">
                <div className="mb-4">
                    <label className="block text-gray-700">Username:</label>
                    <input
                        type="text"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        required
                        className="w-full px-3 py-2 border rounded"
                    />
                </div>
                <div className="mb-4">
                    <label className="block text-gray-700">Password:</label>
                    <input
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                        className="w-full px-3 py-2 border rounded"
                    />
                </div>
                <button type="submit" className="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600">
                    Signup
                </button>
            </form>
        </div>
    );
};

export default Signup;