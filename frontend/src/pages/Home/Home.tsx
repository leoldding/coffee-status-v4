import React, { useEffect, useState } from "react";
import "./Home.css";
import { FaInfoCircle } from "react-icons/fa";

const Home: React.FC = () => {
    const [loading, setLoading] = useState<boolean>(true);
    const [showInfo, setShowInfo] = useState<boolean>(false);
    const [status, setStatus] = useState<string>("");
    const toggleInfo = () => {
        setShowInfo(!showInfo);
    };

    useEffect(() => {
        const getStatus = async () => {
            try {
                const response = await fetch("/api/status", {
                    method: "GET",
                });
                if (!response.ok) {
                    throw new Error("Error retrieving status");
                }
                const json = await response.json();
                setStatus(json.Item.value.S);
                setLoading(false);
            } catch (error) {
                console.error(error);
            }
        };
        getStatus();
    }, []);

    return (
        <>
            <div className="home-container">
                <h1>Is Leo at Think Coffee?</h1>
                <div className="format-container">
                    <div />
                    <div className="status-display">
                        <div className={"status-border " + (loading ? "status-border-loading" : "status-border-not-loading " + status)} />
                        <img className="status-image" src="/coffee_cup.png" alt="coffee cup icon" />
                        <button className="info-button" onClick={toggleInfo}>
                            <FaInfoCircle />
                        </button>
                    </div>
                    <div className={"info-card" + (showInfo ? "" : " hide")}>
                        <p>Green = "Yes"</p>
                        <p>Yellow = "On the Way"</p>
                        <p>Red = "No"</p>
                    </div>
                </div>
            </div>
        </>
    );
};

export default Home;
