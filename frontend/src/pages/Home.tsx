import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { FaInfoCircle, FaLock } from "react-icons/fa";

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
                const response = await fetch("/api/v1/coffee/status", {
                    method: "GET",
                });
                if (!response.ok) {
                    throw new Error("Error retrieving status");
                }
                const json = await response.json();
                setStatus(json.value);
                setLoading(false);
            } catch (error) {
                console.error(error);
            }
        };
        getStatus();
    }, []);

    return (
        <>
            <div className="flex flex-col items-center text-center h-dvh">
                <h1 className="py-8 text-3xl md:text-6xl font-medium">Is Leo at Think Coffee?</h1>
                <div className="h-full flex flex-col items-center justify-center gap-12">
                    <div className="relative h-48 md:h-96 w-48 md:w-96">
                        <div className={"absolute inset-0 border-solid border-4 rounded-full" + (loading ? " border-t-transparent animate-spin " : " ") + (status)} />
                        <div className="absolute inset-0 flex items-center justify-center">
                            <img className="ml-3 md:ml-8 h-36 md:h-72 w-36 md:w-72" src="/coffee_cup.png" alt="coffee cup icon" />
                        </div>
                        <button className="hidden absolute bottom-0 right-0 md:flex items-center justify-center btn btn-soft btn-accent btn-circle" onClick={toggleInfo}>
                            <FaInfoCircle />
                        </button>
                    </div>
                    <div className={"card opacity-100 w-48 bg-accent text-accent-content card-sm shadow-md" + (showInfo ? "" : " md:opacity-0")}>
                        <div className="card-body">
                            <p>Green = "Yes"</p>
                            <p>Yellow = "On the Way"</p>
                            <p>Red = "No"</p>
                        </div>
                    </div>
                </div>
                <div className="fixed bottom-6 right-6">
                    <Link to="/admin">
                        <button className="flex items-center justify-center btn btn-lg btn-soft btn-accent btn-circle">
                            <FaLock />
                        </button>
                    </Link>
                </div>
            </div>
        </>
    );
};

export default Home;
