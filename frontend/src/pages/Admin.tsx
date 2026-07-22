import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { FaHome } from "react-icons/fa";
import { MdCheck, MdClose } from "react-icons/md";

import { api } from "../api/client";

const Admin: React.FC = () => {
    const [authenticated, setAuthenticated] = useState(false);
    const [checkingAuth, setCheckingAuth] = useState(true);
    const [updateStatus, setUpdateStatus] = useState(0);

    useEffect(() => {
        const checkAuth = async () => {
            try {
                const response = await api("/api/v1/auth/check");

                if (!response.ok) {
                    throw new Error("unauthorized");
                }

                setAuthenticated(true);
            } catch (err) {
                console.error(err);
                // api() already redirected if necessary
            } finally {
                setCheckingAuth(false);
            }
        };

        checkAuth();
    }, []);

    const putStatus = async (status: string) => {
        setUpdateStatus(1);

        try {
            const response = await api("/api/v1/coffee/status", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    value: status,
                }),
            });

            if (!response.ok) {
                throw new Error("failed to update");
            }

            setUpdateStatus(2);
        } catch (err) {
            console.error(err);
            setUpdateStatus(3);
        }
    };

    if (checkingAuth) {
        return (
            <div className="h-dvh flex items-center justify-center">
                <span className="loading loading-spinner loading-lg" />
            </div>
        );
    }

    if (!authenticated) {
        return null;
    }

    return (
        <>
            <div className="h-dvh flex flex-col items-center justify-center gap-8">
                <button
                    className="btn btn-xl w-48 btn-success"
                    onClick={() => putStatus("yes")}
                >
                    Yes
                </button>

                <button
                    className="btn btn-xl w-48 btn-warning"
                    onClick={() => putStatus("otw")}
                >
                    On the way
                </button>

                <button
                    className="btn btn-xl w-48 btn-error"
                    onClick={() => putStatus("no")}
                >
                    No
                </button>

                {updateStatus === 1 && (
                    <span className="loading loading-spinner" />
                )}

                {updateStatus === 2 && (
                    <span
                        className="btn btn-success btn-circle btn-xs"
                        onClick={() => setUpdateStatus(0)}
                    >
                        <MdCheck />
                    </span>
                )}

                {updateStatus === 3 && (
                    <span
                        className="btn btn-error btn-circle btn-xs"
                        onClick={() => setUpdateStatus(0)}
                    >
                        <MdClose />
                    </span>
                )}
            </div>

            <div className="fixed bottom-6 right-6">
                <Link to="/">
                    <button className="btn btn-lg btn-soft btn-accent btn-circle">
                        <FaHome />
                    </button>
                </Link>
            </div>
        </>
    );
};

export default Admin;
