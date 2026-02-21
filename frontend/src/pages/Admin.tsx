import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { FaHome } from "react-icons/fa";

const Admin: React.FC = () => {
    const [authenticated, setAuthenticated] = useState<boolean>(false);

    useEffect(() => {
        const checkAuth = async () => {
            try {
                const resp = await fetch("/api/v1/auth/check", {
                    method: "GET",
                    credentials: "include",
                });

                if (!resp.ok) {
                    throw new Error("not authenticated");
                }

                setAuthenticated(true);
            } catch (error) {
                setAuthenticated(false);
                console.error(error);
                window.location.href = "https://auth.leoding.com?redirect=https://coffee.leoding.com/admin"
            }
        };
        checkAuth();
    }, []);

    const putStatus = async (status: string) => {
        try {
            const response = await fetch("/api/v1/coffee/status", {
                method: "POST",
                credentials: "include",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ "value": status }),
            });
            if (!response.ok) {
                throw new Error("Error updating status");
            }
        } catch (error) {
            console.error(error);
        }
    };

    return (
        authenticated &&
        <>
            <div className="h-dvh flex flex-col items-center justify-center gap-8">
                <button className="btn btn-xl w-48 btn-success" onClick={(() => putStatus("yes"))}>Yes</button>
                <button className="btn btn-xl w-48 btn-warning" onClick={(() => putStatus("otw"))}>On the way</button>
                <button className="btn btn-xl w-48 btn-error" onClick={(() => putStatus("no"))}>No</button>
            </div>
            <div className="fixed bottom-6 right-6">
                <Link to="/">
                    <button className="flex items-center justify-center btn btn-lg btn-soft btn-accent btn-circle">
                        <FaHome />
                    </button>
                </Link>
            </div>

        </>
    )
}

export default Admin;
