import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

const Dashboard = () => {
    const navigate = useNavigate();
    const [token, setToken] = useState(null);

    useEffect(() => {
        const fetchToken = async () => {
            try {
                const response = await fetch("http://localhost:8080/api/token", {
                    method: "GET",
                    credentials: "include", // Ensures cookies are sent
                });

                if (!response.ok) {
                    throw new Error("Failed to fetch token");
                }

                const data = await response.json();
                setToken(data.id_token); // Use token as needed
            } catch (error) {
                console.error("Error fetching token:", error);
                navigate("/login"); // Redirect to login if token fetch fails
            }
        };

        fetchToken();
    }, [navigate]);

    return (
        <div>
            <h1>Welcome to the Dashboard</h1>
        </div>
    );
};

export default Dashboard;
